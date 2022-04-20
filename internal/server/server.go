// Package server provides a convenient way to create and start a new server.
// that serves both gRPC and HTTP over a single port.
// with default useful APIs for authentication, health-checking, metrics, tracing, etc.
package server

import (
	"context"
	"github.com/gorilla/mux"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/realHoangHai/awesome/internal/health"
	"github.com/realHoangHai/awesome/internal/middleware/auth"
	"github.com/realHoangHai/awesome/pkg/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
	// register default codecs
	_ "github.com/realHoangHai/awesome/pkg/encoding/json"
	_ "google.golang.org/grpc/encoding/proto"
)

type (
	// Server holds the configuration options for the server instance.
	Server struct {
		lis         net.Listener
		httpSrv     *http.Server
		address     string
		tlsCertFile string
		tlsKeyFile  string

		// HTTP
		readTimeout          time.Duration
		writeTimeout         time.Duration
		shutdownTimeout      time.Duration
		routes               []HandlerOptions
		apiPrefix            string
		httpInterceptors     []func(http.Handler) http.Handler
		routesPrioritization bool
		notFoundHandler      http.Handler

		serverOptions   []grpc.ServerOption
		serveMuxOptions []runtime.ServeMuxOption

		// Interceptors
		streamInterceptors []grpc.StreamServerInterceptor
		unaryInterceptors  []grpc.UnaryServerInterceptor

		log           log.Logger
		enableMetrics bool

		auth auth.Authenticator

		// health checks
		healthCheckPath string
		healthSrv       health.Server
	}

	// Option is a configuration option.
	Option func(*Server)

	// Service implements a registration interface for services to attach themselves to the grpc.Server.
	// If the services support gRPC gateway, they must also implement the EndpointService interface.
	Service interface {
		Register(srv *grpc.Server)
	}

	// EndpointService implement an endpoint registration interface for service to attach their endpoints to gRPC gateway.
	EndpointService interface {
		RegisterWithEndpoint(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption)
	}
)

// New return new server with the given options.
// If address is not set, default address ":8088" will be used.
func New(opts ...Option) *Server {
	server := &Server{}
	for _, opt := range opts {
		opt(server)
	}
	if server.log == nil {
		server.log = log.Root()
	}
	if server.address == "" {
		server.address = defaultAddr
	}
	if server.healthSrv == nil {
		server.healthSrv = health.NewServer(map[string]health.Checker{})
	}
	return server
}

// Run call RunContext with background context.
func (s *Server) Run(services ...Service) error {
	return s.RunWithContext(context.Background(), services...)
}

// RunWithContext opens a tcp listener used by a grpc.Server and a HTTP server,
// and registers each Service with the grpc.Server. If the Service implements EndpointService
// its endpoints will be registered to the HTTP Server running on the same port.
// The server starts with default metrics and health endpoints.
// If the context is canceled or times out, the gRPC server will attempt a graceful shutdown.
func (s *Server) RunWithContext(ctx context.Context, services ...Service) error {
	if s.lis == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	if s.auth != nil {
		s.streamInterceptors = append(s.streamInterceptors, auth.StreamInterceptor(s.auth))
		s.unaryInterceptors = append(s.unaryInterceptors, auth.UnaryInterceptor(s.auth))
	}
	if s.enableMetrics {
		s.streamInterceptors = append(s.streamInterceptors, grpc_prometheus.StreamServerInterceptor)
		s.unaryInterceptors = append(s.unaryInterceptors, grpc_prometheus.UnaryServerInterceptor)
	}
	isSecured := s.tlsCertFile != "" && s.tlsKeyFile != ""

	// server options
	if len(s.streamInterceptors) > 0 {
		s.serverOptions = append(s.serverOptions, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(s.streamInterceptors...)))
	}
	if len(s.unaryInterceptors) > 0 {
		s.serverOptions = append(s.serverOptions, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(s.unaryInterceptors...)))
	}
	if isSecured {
		creds, err := credentials.NewServerTLSFromFile(s.tlsCertFile, s.tlsKeyFile)
		if err != nil {
			return err
		}
		s.serverOptions = append(s.serverOptions, grpc.Creds(creds))
	}
	grpcServer := grpc.NewServer(s.serverOptions...)
	muxOpts := s.serveMuxOptions
	if len(muxOpts) == 0 {
		muxOpts = []runtime.ServeMuxOption{DefaultHeaderMatcher()}
	}
	gw := runtime.NewServeMux(muxOpts...)
	router := mux.NewRouter()

	dialOpts := make([]grpc.DialOption, 0)
	if isSecured {
		creds, err := credentials.NewClientTLSFromFile(s.tlsCertFile, "")
		if err != nil {
			return err
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	}
	if !isSecured {
		s.log.Context(ctx).Warn("server: insecured mode is enabled.")
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}
	// expose health service via gRPC.
	services = append(services, s.healthSrv)

	for _, svc := range services {
		svc.Register(grpcServer)
		if epSrv, ok := svc.(EndpointService); ok {
			epSrv.RegisterWithEndpoint(ctx, gw, s.address, dialOpts)
		}
	}
	// Make sure Prometheus metrics are initialized.
	if s.enableMetrics {
		grpc_prometheus.Register(grpcServer)
	}
	// Add internal handlers.
	s.routes = append([]HandlerOptions{
		{
			p: s.getHealthCheckPath(),
			h: s.healthSrv,
			m: []string{http.MethodGet},
		},
	}, s.routes...)
	// Serve gRPC and GW only and only if there is at least one service registered.
	if len(services) > 0 {
		s.routes = append(s.routes, HandlerOptions{p: s.getAPIPrefix(), h: gw, prefix: true})
	}
	// register all http handlers to the router.
	s.registerHTTPHandlers(ctx, router)

	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	handler := grpcHandlerFunc(isSecured, grpcServer, router)
	for i := len(s.httpInterceptors) - 1; i >= 0; i-- {
		handler = s.httpInterceptors[i](handler)
	}
	s.httpSrv = &http.Server{
		Addr:         s.address,
		Handler:      handler,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}
	go func() {
		if isSecured {
			errChan <- s.httpSrv.ServeTLS(s.lis, s.tlsCertFile, s.tlsKeyFile)
			return
		}
		errChan <- s.httpSrv.Serve(s.lis)
	}()

	// init health check service.
	if err := s.healthSrv.Init(health.StatusServing); err != nil {
		s.log.Context(ctx).Error("server: start health check server, err: %v", err)
		s.Shutdown(ctx)
		return err
	}
	defer func() {
		_ = s.healthSrv.Close()
	}()
	s.log.Context(ctx).Infof("server: listening at: %s", s.address)
	select {
	case <-ctx.Done():
		s.Shutdown(ctx)
		return ctx.Err()
	case err := <-errChan:
		return err
	case sig := <-sigChan:
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			s.log.Context(ctx).Info("server: gracefully shutdown...")
			s.Shutdown(ctx)
		}
	}
	return nil
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise.
func grpcHandlerFunc(isSecured bool, grpcServer *grpc.Server, mux http.Handler) http.Handler {
	if isSecured {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isGRPCRequest(r) {
				grpcServer.ServeHTTP(w, r)
				return
			}
			mux.ServeHTTP(w, r)
		})
	}
	// work-around in case TLS is disabled. See: https://github.com/grpc/grpc-go/issues/555
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isGRPCRequest(r) {
			grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func (s *Server) getHealthCheckPath() string {
	if s.healthCheckPath == "" {
		return "/internal/health"
	}
	return s.healthCheckPath
}

// With allows user to add more options to the server after created.
func (s *Server) With(opts ...Option) *Server {
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Address return address that the server is listening.
func (s *Server) Address() string {
	return s.address
}

func (s *Server) getAPIPrefix() string {
	if s.apiPrefix == "" {
		return "/"
	}
	return s.apiPrefix
}

// Shutdown shutdown the server gracefully.
func (s *Server) Shutdown(ctx context.Context) {
	if s.healthSrv != nil {
		if err := s.healthSrv.Close(); err != nil {
			s.log.Errorf("server: shutdown health check service error: %v", err)
		}
	}
	if s.httpSrv != nil {
		ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
		defer cancel()
		if err := s.httpSrv.Shutdown(ctx); err != nil {
			s.log.Errorf("server: shutdown error: %v", err)
		}
	}
}

func (s *Server) getLogger() log.Logger {
	if s.log == nil {
		return log.Root()
	}
	return s.log
}

func (s *Server) registerHTTPHandlers(ctx context.Context, router *mux.Router) {
	// Longer patterns take precedence over shorter ones.
	if s.routesPrioritization {
		sort.Sort(sort.Reverse(handlerOptionsSlice(s.routes)))
	}
	if s.notFoundHandler != nil {
		router.NotFoundHandler = s.notFoundHandler
	}
	for _, r := range s.routes {
		var route *mux.Route
		h := r.h
		info := make([]interface{}, 0)
		for _, interceptor := range r.interceptors {
			h = interceptor(h)
		}
		if r.prefix {
			route = router.PathPrefix(r.p).Handler(h)
			info = append(info, "path_prefix", r.p)
		} else {
			route = router.Path(r.p).Handler(h)
			info = append(info, "path", r.p)
		}
		if r.m != nil {
			route.Methods(r.m...)
			info = append(info, "methods", r.m)
		}
		if r.q != nil {
			route.Queries(r.q...)
			info = append(info, "queries", r.q)
		}
		if r.hdr != nil {
			route.Headers(r.hdr...)
			info = append(info, "headers", r.hdr)
		}
		s.log.Context(ctx).Fields(info...).Infof("server: registered HTTP handler")
	}
}

func isGRPCRequest(r *http.Request) bool {
	return r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc")
}

package server

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/realHoangHai/awesome/config"
	"github.com/realHoangHai/awesome/internal/auth"
	"github.com/realHoangHai/awesome/internal/health"
	"github.com/realHoangHai/awesome/pkg/jwt"
	"github.com/realHoangHai/awesome/pkg/log"
	"github.com/realHoangHai/awesome/pkg/utils/header"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
	"net/http/pprof"
	"net/textproto"
	"os"
	"time"
)

const (
	defaultAddr = ":8088"
)

// FromEnv is an option to create a new server from environment variables configuration.
// See config.Config for the available options.
func FromEnv(cfg *config.Config) Option {
	return func(opts *Server) {
		FromConfig(cfg)(opts)
	}
}

// FromConfig is an option to create a new server from an existing config.
func FromConfig(cfg *config.Config) Option {
	return func(server *Server) {
		opts := []Option{
			UnaryInterceptors(CorrelationIDUnaryInterceptor()),
			StreamInterceptors(CorrelationIDStreamInterceptor()),
			Address(cfg.Core.Address),
			TLS(cfg.Core.TLSKeyFile, cfg.Core.TLSCertFile),
			Timeout(cfg.Core.ReadTimeout, cfg.Core.WriteTimeout),
			JWT(cfg.Core.JWTSecret),
			APIPrefix(cfg.Core.APIPrefix),
			CORS(cfg.Core.CORSAllowedCredential, cfg.Core.CORSAllowedHeaders, cfg.Core.CORSAllowedMethods, cfg.Core.CORSAllowedOrigins),
			ShutdownTimeout(cfg.Core.ShutdownTimeout),
			RoutesPrioritization(cfg.Core.RoutesPrioritization),
			ShutdownHook(cfg.Core.ShutdownHook),
		}
		if cfg.Core.Metrics {
			opts = append(opts, Metrics(cfg.Core.MetricsPath))
		}
		if cfg.Core.WebDir != "" {
			opts = append(opts, Web(cfg.Core.WebPrefix, cfg.Core.WebDir, cfg.Core.WebIndex))
		}
		// context logger
		if cfg.Core.ContextLogger {
			logger := log.Root()
			if cfg.Core.Name != "" {
				logger = logger.Fields("name", cfg.Core.Name)
			}
			opts = append(opts, Logger(logger))
		}
		// create health check by default
		opts = append(opts, HealthCheck(cfg.Core.HealthCheckPath,
			health.NewServer(map[string]health.Checker{},
				health.Logger(server.getLogger()))))
		// recovery
		if cfg.Core.Recovery {
			opts = append(opts, Recovery(nil))
		}
		if cfg.Core.PProf {
			opts = append(opts, PProf(cfg.Core.PProfPrefix))
		}
		// apply all
		for _, opt := range opts {
			opt(server)
		}
	}
}

// Address is an option to set address.
// Default address is :8000
func Address(addr string) Option {
	return func(opts *Server) {
		opts.address = addr
		if opts.lis != nil {
			opts.getLogger().Debugf("server: address is set to %s, Listener will be overridden", addr)
			opts.lis = nil
		}
	}
}

// Listener is an option allows server to be served on an existing listener.
func Listener(lis net.Listener) Option {
	return func(opts *Server) {
		if opts.address != "" {
			opts.getLogger().Debugf("server: listener is set to %s, address will be overridden", lis.Addr().String())
			opts.address = lis.Addr().String()
		}
		opts.lis = lis
	}
}

// CorrelationIDStreamInterceptor returns a grpc.StreamServerInterceptor that provides
// a context with correlation_id for tracing. It will try to looks for value of X-Correlation-ID or X-Request-ID
// in the metadata of the incoming request. If no value is provided, a new UUID will be generated.
func CorrelationIDStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		id, ok := header.CorrelationIDFromContext(ss.Context())
		if ok {
			return handler(srv, ss)
		}
		wrapped := grpc_middleware.WrapServerStream(ss)
		md := metadata.Pairs(header.XCorrelationID, id)
		if imd, ok := metadata.FromIncomingContext(ss.Context()); ok {
			md = metadata.Join(md, imd)
		}
		wrapped.WrappedContext = metadata.NewIncomingContext(ss.Context(), md)

		return handler(srv, wrapped)
	}
}

// StreamInterceptors is an option allows user to add additional stream interceptors to the server.
func StreamInterceptors(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(opts *Server) {
		opts.streamInterceptors = append(opts.streamInterceptors, interceptors...)
	}
}

// CorrelationIDUnaryInterceptor returns a grpc.UnaryServerInterceptor that provides
// a context with correlation_id for tracing. It will try to looks for value of X-Correlation-ID or X-Request-ID
// in the metadata of the incoming request. If no value is provided, a new UUID will be generated.
func CorrelationIDUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		id, ok := header.CorrelationIDFromContext(ctx)
		if ok {
			return handler(ctx, req)
		}
		md := metadata.Pairs(header.XCorrelationID, id)
		if imd, ok := metadata.FromIncomingContext(ctx); ok {
			md = metadata.Join(md, imd)
		}
		newCtx := metadata.NewIncomingContext(ctx, md)
		return handler(newCtx, req)
	}
}

// UnaryInterceptors is an option allows user to add additional unary interceptors to the server.
func UnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(opts *Server) {
		opts.unaryInterceptors = append(opts.unaryInterceptors, interceptors...)
	}
}

// JWT is an option allows user to use jwt authenticator for authentication.
func JWT(secret string) Option {
	return func(opts *Server) {
		if secret == "" {
			return
		}
		opts.auth = jwt.Authenticator([]byte(secret))
	}
}

// Auth is an option allows user to use an authenticator for authentication.
// Find more about authenticators in auth package.
func Auth(f auth.Authenticator) Option {
	return func(opts *Server) {
		opts.auth = f
	}
}

// Logger is an option allows user to add a custom logger into the server.
func Logger(logger log.Logger) Option {
	return func(opts *Server) {
		opts.log = logger
		opts.serveMuxOptions = append(opts.serveMuxOptions, DefaultHeaderMatcher())
		opts.streamInterceptors = append(opts.streamInterceptors, log.StreamInterceptor(logger))
		opts.unaryInterceptors = append(opts.unaryInterceptors, log.UnaryInterceptor(logger))
	}
}

// TLS is an option allows user to add TLS for transport security to the server.
// Note that host name in ADDRESS must be configured accordingly. Otherwise, you
// might encounter TLS handshake error.
// TIP: for local testing, take a look at https://github.com/FiloSottile/mkcert.
func TLS(key, cert string) Option {
	return func(opts *Server) {
		if key == "" || cert == "" {
			return
		}
		opts.tlsKeyFile = key
		opts.tlsCertFile = cert
		// server/dial options will be handled in server.go
	}
}

// Timeout is an option to override default read/write timeout.
func Timeout(read, write time.Duration) Option {
	if read == 0 {
		read = 30 * time.Second
	}
	if write == 0 {
		write = 30 * time.Second
	}
	return func(opts *Server) {
		opts.readTimeout = read
		opts.writeTimeout = write
		opts.serverOptions = append(opts.serverOptions, grpc.ConnectionTimeout(read))
	}
}

// ServeMuxOptions is an option allows user to add additional ServeMuxOption.
func ServeMuxOptions(muxOpts ...runtime.ServeMuxOption) Option {
	return func(opts *Server) {
		opts.serveMuxOptions = append(opts.serveMuxOptions, muxOpts...)
	}
}

// Options is an option allows user to add additional grpc.ServerOption.
func Options(serverOpts ...grpc.ServerOption) Option {
	return func(opts *Server) {
		opts.serverOptions = append(opts.serverOptions, serverOpts...)
	}
}

// HealthCheck is an option allows user to provide a custom health check server.
func HealthCheck(path string, srv health.Server) Option {
	return func(opts *Server) {
		opts.healthCheckPath = path
		opts.healthSrv = srv
	}
}

// HealthChecks is an option allows user to provide custom health checkers
// using default health check server.
func HealthChecks(checkers map[string]health.Checker) Option {
	return func(opts *Server) {
		opts.healthCheckPath = opts.getHealthCheckPath()
		opts.healthSrv = health.NewServer(checkers, health.Logger(opts.getLogger()))
	}
}

// AddressFromEnv is an option allows user to set address using environment configuration.
// It looks for PORT and then ADDRESS variables.
// This option is mostly used for cloud environment like Heroku where the port is randomly set.
func AddressFromEnv(cfg *config.Config) Option {
	return func(srv *Server) {
		srv.address = GetAddressFromEnv(cfg)
	}
}

// GetAddressFromEnv returns address from configured environment variables: PORT or ADDRESS.
// This function prioritizes PORT over ADDRESS.
// If non of the variables is configured, return default address.
func GetAddressFromEnv(cfg *config.Config) string {
	port := os.Getenv("PORT")
	if port != "" {
		return fmt.Sprintf(":%s", port)
	}
	if cfg.Core.Address != "" {
		return cfg.Core.Address
	}
	return defaultAddr
}

// Handler is an option allows user to add additional HTTP handlers.
// Longer patterns take precedence over shorter ones by default,
// use RoutesPrioritization option to disable this rule.
// See github.com/gorilla/mux for defining path with variables/patterns.
//
// For more options, use HTTPHandlerX.
func Handler(path string, h http.Handler, methods ...string) Option {
	return HandlerWithOptions(path, h, NewHandlerOptions().Methods(methods...))
}

// HandlerFunc is an option similar to HTTPHandler, but for http.HandlerFunc.
func HandlerFunc(path string, h func(http.ResponseWriter, *http.Request), methods ...string) Option {
	return HandlerWithOptions(path, http.HandlerFunc(h), NewHandlerOptions().Methods(methods...))
}

// PrefixHandler is an option to quickly define a prefix HTTP handler.
// For more options, please use HTTPHandlerX.
func PrefixHandler(path string, h http.Handler, methods ...string) Option {
	return HandlerWithOptions(path, h, NewHandlerOptions().Prefix().Methods(methods...))
}

// NotFoundHandler is an option to provide a custom not found HTTP Handler.
func NotFoundHandler(h http.Handler) Option {
	return func(opts *Server) {
		opts.notFoundHandler = h
	}
}

// HandlerWithOptions is an option to define full options such as method, query, header matchers
// and interceptors for a HTTP handler.
// Longer patterns take precedence over shorter ones by default,
// use RoutesPrioritization option to disable this rule.
// See github.com/gorilla/mux for defining path with variables/patterns.
func HandlerWithOptions(path string, h http.Handler, hopt *HandlerOptions) Option {
	return func(opts *Server) {
		hopt.p = path
		hopt.h = h
		opts.routes = append(opts.routes, *hopt)
	}
}

// RoutesPrioritization enable/disable the routes prioritization.
func RoutesPrioritization(enable bool) Option {
	return func(opts *Server) {
		opts.routesPrioritization = enable
	}
}

// HTTPInterceptors is an option allows user to set additional interceptors to the root HTTP handler.
// If interceptors are applied to gRPC, it is required that the interceptors don't hijack the response writer,
// otherwise panic "Hijack not supported" will be thrown.
func HTTPInterceptors(interceptors ...HTTPInterceptor) Option {
	return func(opts *Server) {
		opts.httpInterceptors = append(opts.httpInterceptors, interceptors...)
	}
}

// APIPrefix is an option allows user to route only the specified path prefix to gRPC Gateway.
// This option is used mostly when you serve both gRPC APIs along with other internal HTTP APIs.
// The default prefix is /, which will route all paths to gRPC Gateway.
func APIPrefix(prefix string) Option {
	return func(opts *Server) {
		opts.apiPrefix = prefix
	}
}

// Web is an option to allows user to serve Web/Single Page Application
// along with API Gateway and gRPC. API Gateway must be served in a
// different path prefix with the web path prefix.
func Web(pathPrefix, dir, index string) Option {
	if pathPrefix == "" {
		pathPrefix = "/"
	}
	return PrefixHandler(pathPrefix, spaHandler{
		index: index,
		dir:   dir,
	})
}

// Recovery is an option allows user to add an ability to recover a handler/API from a panic.
// This applies for both unary and stream handlers/APIs.
// If the provided error handler is nil, a default error handler will be used.
func Recovery(handler func(context.Context, interface{}) error) Option {
	return func(opts *Server) {
		if handler == nil {
			if opts.log == nil {
				opts.log = log.Root()
			}
			handler = recoveryHandler(opts.log)
		}
		recoverOpt := recovery.WithRecoveryHandlerContext(handler)
		opts.unaryInterceptors = append(opts.unaryInterceptors, recovery.UnaryServerInterceptor(recoverOpt))
		opts.streamInterceptors = append(opts.streamInterceptors, recovery.StreamServerInterceptor(recoverOpt))
	}
}

// CORS is an option allows users to enable CORS on the server.
func CORS(allowCredential bool, headers, methods, origins []string) Option {
	return func(opts *Server) {
		options := []handlers.CORSOption{}
		if allowCredential {
			options = append(options, handlers.AllowCredentials())
		}
		if headers != nil {
			options = append(options, handlers.AllowedHeaders(headers))
		}
		if methods != nil {
			options = append(options, handlers.AllowedMethods(methods))
		}
		if origins != nil {
			options = append(options, handlers.AllowedOrigins(origins))
		}
		if len(options) > 0 {
			opts.httpInterceptors = append(opts.httpInterceptors, handlers.CORS(options...))
		}
	}
}

// PProf is an option allows user to enable Go profiler.
func PProf(pathPrefix string) Option {
	return func(opts *Server) {
		opts.routes = append(opts.routes, HandlerOptions{
			p: pathPrefix + "/debug/pprof/",
			h: http.HandlerFunc(pprof.Index),
		})
		opts.routes = append(opts.routes, HandlerOptions{
			p: pathPrefix + "/debug/pprof/cmdline",
			h: http.HandlerFunc(pprof.Cmdline),
		})
		opts.routes = append(opts.routes, HandlerOptions{
			p: pathPrefix + "/debug/pprof/profile",
			h: http.HandlerFunc(pprof.Profile),
		})
		opts.routes = append(opts.routes, HandlerOptions{
			p: pathPrefix + "/debug/pprof/symbol",
			h: http.HandlerFunc(pprof.Symbol),
		})
		opts.routes = append(opts.routes, HandlerOptions{
			p: pathPrefix + "/debug/pprof/trace",
			h: http.HandlerFunc(pprof.Trace),
		})
	}
}

// Metrics is an option to register standard Prometheus metrics for HTTP.
// Default path is /internal/metrics.
func Metrics(path string) Option {
	return func(opts *Server) {
		p := path
		if p == "" {
			p = "/internal/metrics"
			opts.getLogger().Warnf("metrics path is switched automatically to %s", p)
		}
		opts.enableMetrics = true
		opts.routes = append(opts.routes, HandlerOptions{
			p: p,
			h: promhttp.Handler(),
			m: []string{http.MethodGet},
		})
	}
}

// ShutdownTimeout is an option to override default shutdown timeout of server.
// Set to -1 for no timeout.
func ShutdownTimeout(t time.Duration) Option {
	return func(opts *Server) {
		opts.shutdownTimeout = t
	}
}

// DefaultHeaderMatcher is an ServerMuxOption that forward
// header keys X-Request-Id, X-Correlation-ID, Api-Key to gRPC Context.
func DefaultHeaderMatcher() runtime.ServeMuxOption {
	return HeaderMatcher([]string{"X-Request-Id", "X-Correlation-ID", "Api-Key"})
}

// HeaderMatcher is an ServeMuxOption for matcher header
// for passing a set of non IANA headers to gRPC context
// without a need to prefix them with Grpc-Metadata.
func HeaderMatcher(keys []string) runtime.ServeMuxOption {
	return runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		canonicalKey := textproto.CanonicalMIMEHeaderKey(key)
		for _, k := range keys {
			if k == canonicalKey || textproto.CanonicalMIMEHeaderKey(k) == canonicalKey {
				return k, true
			}
		}
		return runtime.DefaultHeaderMatcher(key)
	})
}

// Tracing is an option to enable tracing on unary requests.
func Tracing(tracer opentracing.Tracer) Option {
	return UnaryInterceptors(otgrpc.OpenTracingServerInterceptor(tracer))
}

// StreamTracing is an option to enable tracing on stream requests.
func StreamTracing(tracer opentracing.Tracer) Option {
	return StreamInterceptors(otgrpc.OpenTracingStreamServerInterceptor(tracer))
}

// ShutdownHook expose an API for shutdown the server remotely.
//
// WARNING: this is an experiment API and
// it should be enabled only in development mode for live reload.
func ShutdownHook(path string) Option {
	return func(opts *Server) {
		// do nothing if path is empty.
		if path == "" {
			return
		}
		hopt := HandlerOptions{}
		hopt.p = path
		hopt.h = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			opts.Shutdown(context.Background())
		})
		opts.routes = append(opts.routes, hopt)
	}
}

// recoveryHandler print the context log to the configured writer and return
// a general error to the caller.
func recoveryHandler(l log.Logger) func(context.Context, interface{}) error {
	return func(ctx context.Context, p interface{}) error {
		l.Context(ctx).Errorf("server: panic recovered, err: %v", p)
		return status.Errorf(codes.Internal, codes.Internal.String())
	}
}

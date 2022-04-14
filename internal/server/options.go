package server

import (
	"context"
	"github.com/gorilla/handlers"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/realHoangHai/awesome/config"
	"github.com/realHoangHai/awesome/pkg/log"
	"github.com/realHoangHai/awesome/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
	"net/textproto"
	"time"
)

const (
	defaultAddr = ":8088"
)

// FromEnv is an option to create a new server from environment variables configuration.
// See Config for the available options.
func FromEnv(cfg *config.Config) Option {
	return func(opts *Server) {
		FromConfig(cfg)(opts)
	}
}

func FromConfig(cfg *config.Config) Option {
	return func(server *Server) {
		opts := []Option{
			UnaryInterceptors(CorrelationIDUnaryInterceptor()),
			StreamInterceptors(CorrelationIDStreamInterceptor()),
			Address(cfg.Core.Address),
			TLS(cfg.Core.TLSKeyFile, cfg.Core.TLSCertFile),
			Timeout(cfg.Core.ReadTimeout, cfg.Core.WriteTimeout),
			APIPrefix(cfg.Core.APIPrefix),
			CORS(cfg.Core.CORSAllowedCredential, cfg.Core.CORSAllowedHeaders, cfg.Core.CORSAllowedMethods, cfg.Core.CORSAllowedOrigins),
			ShutdownTimeout(cfg.Core.ShutdownTimeout),
		}
		// context logger
		if cfg.Core.ContextLogger {
			logger := log.Root()
			if cfg.Core.Name != "" {
				logger = logger.Fields("name", cfg.Core.Name)
			}
			opts = append(opts, Logger(logger))
		}
		// recovery
		if cfg.Core.Recovery {
			opts = append(opts, Recovery(nil))
		}
		// apply options
		for _, opt := range opts {
			opt(server)
		}
	}
}

// Address is an option to set address.
// Default address is :8088
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

// CorrelationIDUnaryInterceptor returns a grpc.UnaryServerInterceptor that provides
// a context with correlation_id for tracing. It will try to looks for value of X-Correlation-ID or X-Request-ID
// in the metadata of the incoming request. If no value is provided, a new UUID will be generated.
func CorrelationIDUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		id, ok := utils.CorrelationIDFromContext(ctx)
		if ok {
			return handler(ctx, req)
		}
		md := metadata.Pairs(utils.XCorrelationID, id)
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

// CorrelationIDStreamInterceptor returns a grpc.StreamServerInterceptor that provides
// a context with correlation_id for tracing. It will try to looks for value of X-Correlation-ID or X-Request-ID
// in the metadata of the incoming request. If no value is provided, a new UUID will be generated.
func CorrelationIDStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		id, ok := utils.CorrelationIDFromContext(ss.Context())
		if ok {
			return handler(srv, ss)
		}
		wrapped := grpc_middleware.WrapServerStream(ss)
		md := metadata.Pairs(utils.XCorrelationID, id)
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
		var options []handlers.CORSOption
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

// recoveryHandler print the context log to the configured writer and return
// a general error to the caller.
func recoveryHandler(l log.Logger) func(context.Context, interface{}) error {
	return func(ctx context.Context, p interface{}) error {
		l.Context(ctx).Errorf("server: panic recovered, err: %v", p)
		return status.Errorf(codes.Internal, codes.Internal.String())
	}
}

package auth

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"net/http"
)

// StreamInterceptor returns a new grpc.StreamClientInterceptor that performs an authentication
// check for each request by usin Authenticator.Authenticate(ctx context.Context)
func StreamInterceptor(auth Authenticator) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var newCtx context.Context
		// if server overides Authenticator, use it instead
		if srvAuth, ok := srv.(Authenticator); ok {
			auth = srvAuth
		}
		if wl, ok := auth.(WhiteListAuthenticator); ok && wl.IsWhiteListed(info.FullMethod) {
			newCtx = ss.Context()
		} else {
			newCtx, err = auth.Authenticate(ss.Context())
		}
		if err != nil {
			return
		}
		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

// UnaryInterceptor returns a grpc.UnaryServerInterceptor that performs an authentication
// check for each request by using Authenticator.Authenticate(ctx context.Context).
func UnaryInterceptor(auth Authenticator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var newCtx context.Context
		// if server override Authenticator, use it instead.
		if srvAuth, ok := info.Server.(Authenticator); ok {
			auth = srvAuth
		}
		if wl, ok := auth.(WhiteListAuthenticator); ok && wl.IsWhiteListed(info.FullMethod) {
			newCtx, err = ctx, nil
		} else {
			newCtx, err = auth.Authenticate(ctx)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// HTTPInterceptor return a HTTP interceptor that perform an authentication check
// for each request using the given authenticator.
func HTTPInterceptor(auth Authenticator) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if wl, ok := auth.(WhiteListAuthenticator); ok && wl.IsWhiteListed(r.URL.Path) {
				h.ServeHTTP(w, r)
				return
			}
			token := tokenString(r.Context())
			if token == "" {
				token = r.Header.Get(AuthorizationMD)
			}
			if token == "" {
				t, err := r.Cookie(AuthorizationMD)
				if err == nil {
					token = t.Value
				}
			}
			md := metadata.MD{}
			if v, ok := metadata.FromIncomingContext(r.Context()); ok {
				md = v.Copy()
			}
			md.Set(AuthorizationMD, token)
			newCtx, err := auth.Authenticate(metadata.NewIncomingContext(r.Context(), md))
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"code":%d,"message":"%s"}`, codes.Unauthenticated, err), http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

func tokenString(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	slice, ok := md[AuthorizationMD]
	if !ok || len(slice) == 0 {
		return ""
	}
	if len(slice) > 1 {
		return ""
	}
	return slice[0]
}

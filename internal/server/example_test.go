package server_test

//import (
//	"context"
//	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
//	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/tags"
//	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/tracing"
//	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
//	"github.com/realHoangHai/awesome/internal/server"
//	"github.com/realHoangHai/awesome/pkg/log"
//	"net/http"
//)
//
//func ExampleListenAndServe() {
//	if err := server.ListenAndServe( /*services ...Service*/ ); err != nil {
//		panic(err)
//	}
//}
//
//func ExampleListenAndServeContext() {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	if err := server.ListenAndServeContext(ctx /*, services ...Service*/); err != nil {
//		panic(err)
//	}
//}
//
//func ExampleNew_fromEnvironmentVariables() {
//	srv := server.New(server.FromEnv())
//	if err := srv.ListenAndServe( /*services ...Service*/ ); err != nil {
//		log.Panic(err)
//	}
//}
//
//func ExampleNew_withOptions() {
//	srv := server.New(
//		server.Address(":8088"),
//		server.JWT("iloveu"),
//		server.Logger(log.Fields("service", "awesome")),
//	)
//	if err := srv.ListenAndServe( /*services ...Service*/ ); err != nil {
//		panic(err)
//	}
//}
//
//func ExampleNew_withInternalHTTPAPI() {
//	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("doc"))
//	})
//	srv := server.New(
//		server.FromEnv(),
//		server.Handler("/doc", h),
//	)
//	if err := srv.ListenAndServe( /*services ...Service*/ ); err != nil {
//		panic(err)
//	}
//}
//
//func ExampleNew_withExternalInterceptors() {
//	srv := server.New(
//		server.FromEnv(),
//		server.StreamInterceptors(
//			tags.StreamServerInterceptor(),
//			tracing.StreamServerInterceptor(),
//			grpc_prometheus.StreamServerInterceptor,
//			recovery.StreamServerInterceptor(),
//		),
//		server.UnaryInterceptors(
//			tags.UnaryServerInterceptor(),
//			tracing.UnaryServerInterceptor(),
//			grpc_prometheus.UnaryServerInterceptor,
//			recovery.UnaryServerInterceptor(),
//		),
//	)
//	if err := srv.ListenAndServe( /*services ...Service*/ ); err != nil {
//		panic(err)
//	}
//}

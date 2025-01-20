package interceptor

import (
	"context"
	"log"
	"net/http"

	"github.com/felixge/httpsnoop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func WithLoggerInterceptor(handler http.Handler) http.Handler {
	// the create a handler
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// pass the handler to httpsnoop to get http status and latency
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		// printing extracted data
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}

func UnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// extracting metadata from context
	md, _ := metadata.FromIncomingContext(ctx)
	log.Printf("unary server interceptor: %v\n", md)
	log.Printf("received request on method: %s", info.FullMethod) // 被调用的RPC的完整方法名
	// calling the handler
	resp, err := handler(ctx, req)
	log.Printf("sending response from method: %s", info.FullMethod)
	return resp, err
}

func StreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// extracting metadata from context
	md, _ := metadata.FromIncomingContext(ss.Context())
	log.Printf("stream server interceptor: %v\n", md)
	log.Printf("received request on method: %s", info.FullMethod)
	// calling the handler
	err := handler(srv, ss)
	log.Printf("sending response from method: %s", info.FullMethod)
	return err
}

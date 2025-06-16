package grpcmw

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		log.Printf("gRPC request started: method: %s, request: %s", info.FullMethod, req)

		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			log.Printf(
				"gRPC request failed: method: %s, request: %s, duration: %s error: %v",
				info.FullMethod,
				req,
				duration,
				err,
			)
		} else {
			log.Printf("gRPC request succeeded: method: %s, request: %s, duration: %s", info.FullMethod, req, duration)
		}

		return resp, err
	}
}

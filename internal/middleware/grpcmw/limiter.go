package grpcmw

import (
	"context"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RateLimiter(rps int) grpc.UnaryServerInterceptor {
	limiter := rate.NewLimiter(rate.Limit(rps), rps)

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := limiter.Wait(ctx); err != nil {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(ctx, req)
	}
}

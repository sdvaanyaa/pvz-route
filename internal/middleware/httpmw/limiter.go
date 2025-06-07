package httpmw

import (
	"golang.org/x/time/rate"
	"net/http"
)

func RateLimiter(rps int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rps), rps)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if err := limiter.Wait(ctx); err != nil {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

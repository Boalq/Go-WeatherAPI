package functions

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func RateLimiterMiddleWare(next http.HandlerFunc) http.HandlerFunc{

	limiter := rate.NewLimiter(10, 20)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "Too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}

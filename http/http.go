package http

import (
	"github.com/althk/ratelimiter"
	"net"
	"net/http"
)

// WithLimiter wraps the given handler with a rate limiter.
// The rate limiter algorithm and the store type can be configured from the available options.
func WithLimiter(next func(w http.ResponseWriter, r *http.Request), opts *ratelimiter.LimiterOptions) http.Handler {
	limiter := ratelimiter.Build(opts)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !limiter.Allow(r.Context(), ip) {
			w.Header().Set("Content-Type", "plain/text")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte("Quota Exceeded\n"))
			return
		}
		next(w, r)
	})
}

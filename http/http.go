package http

import (
	"github.com/althk/ratelimiter"
	"github.com/althk/ratelimiter/store"
	"github.com/althk/ratelimiter/types"
	"github.com/redis/go-redis/v9"
	"net"
	"net/http"
	"time"
)

type LimiterStore byte
type LimiterAlgo byte

const (
	Local LimiterStore = iota
	Redis
)

const (
	TokenBucket LimiterAlgo = iota
	SlidingWindow
)

// buildLimiter builds a limiter with the given algorithm and storage driver.
// TODO: make all options in this method configurable
func buildLimiter(algo LimiterAlgo, storeType LimiterStore) types.Limiter {
	var s types.Store
	var l types.Limiter
	switch storeType {
	case Local:
		s = store.NewInMemory(time.Minute)
	case Redis:
		s = store.NewRedis(redis.NewClient(&redis.Options{Addr: "localhost:6379"}), time.Minute)
	}
	switch algo {
	case TokenBucket:
		l = ratelimiter.NewTBL(1, 1, s)
	case SlidingWindow:
		l = ratelimiter.NewSWL(time.Second, 1, s)
	}
	return l
}

// WithLimiter wraps the given handler with a rate limiter.
// The rate limiter algorithm and the store type can be configured from the available options.
func WithLimiter(next func(w http.ResponseWriter, r *http.Request), algo LimiterAlgo, storeType LimiterStore) http.Handler {
	limiter := buildLimiter(algo, storeType)

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

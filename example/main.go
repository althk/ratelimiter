package main

import (
	"github.com/althk/ratelimiter"
	rlhttp "github.com/althk/ratelimiter/http"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

// ping simply responds with a "pong" plain text response for all requests
func ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong\n"))
}

func main() {
	// wrap the ping handler with a redis backed sliding window limiter
	// configure the limiter to 1 request per sec per IP
	opts := &ratelimiter.LimiterOptions{
		// rate limiting algorithm
		Algo: ratelimiter.SlidingWindow,

		// the storage backend for the rate limiter
		StoreType: ratelimiter.Redis,

		// The max burst rate, or limit per window/bucket
		ReqLimit: 1,

		// TokenBucketRate is needed when Algo is set to ratelimiter.TokenBucket
		// TokenBucketRate:       0,

		// SlidingWindowDuration is needed when Algo is set to ratelimiter.SlidingWindow
		SlidingWindowDuration: time.Second,

		// RedisOpts is needed when StoreType is ratelimiter.Redis
		RedisOpts: &redis.Options{Addr: "localhost:6379"},
	}
	limitedPing := rlhttp.WithLimiter(ping, opts)

	http.HandleFunc("/unlimitedping", ping)

	// configure the route /limitedping to be handled by the rate limiter handler
	http.Handle("/limitedping", limitedPing)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

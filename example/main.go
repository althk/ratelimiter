package main

import (
	rlhttp "github.com/althk/ratelimiter/http"
	"log"
	"net/http"
)

// ping simply responds with a "pong" plain text response for all requests
func ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong\n"))
}

func main() {
	// wrap the ping handler with a local limiter
	// this limiter is configured to 1 request per sec per IP by default
	limitedPing := rlhttp.WithLimiter(ping, rlhttp.SlidingWindow, rlhttp.Redis)

	http.HandleFunc("/unlimitedping", ping)

	// configure the route /limitedping to be handled by the rate limiter handler
	http.Handle("/limitedping", limitedPing)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

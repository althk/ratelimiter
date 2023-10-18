package types

import (
	"context"
	"time"
)

// Limiter defines the interface for a rate limiter.
// The API is meant for simple use cases. For advanced use cases,
// prefer https://pkg.go.dev/golang.org/x/time/rate
type Limiter interface {

	// Allow returns true if the current request can be processed.
	// Returns false if the client has exceeded their allotted quota for the current time frame.
	Allow(ctx context.Context, id string) bool

	// AllowN returns true if n requests can be performed by the client at the given time.
	AllowN(ctx context.Context, id string, t time.Time, n int64) bool
}

type Client struct {
	data     []byte
	LastSeen time.Time
}

type Store interface {
	Get(ctx context.Context, id string) ([]byte, bool)
	Set(ctx context.Context, id string, data []byte)
}

package store

import (
	"context"
	"github.com/althk/ratelimiter/types"
	"github.com/redis/go-redis/v9"
	"time"
)

// Redis provides the ability to save to and retrieve key/value pairs
// from a Redis server.
// It implements types.Store interface to be used by rate limiters.
type Redis struct {
	// Redis client connection
	cli *redis.Client

	// the duration after which a key should expire
	ttl time.Duration
}

var _ types.Store = new(Redis)

func (r *Redis) Get(ctx context.Context, id string) ([]byte, bool) {
	v, err := r.cli.Get(ctx, id).Bytes()
	if err == redis.Nil {
		return nil, false
	}
	return v, true
}

func (r *Redis) Set(ctx context.Context, id string, data []byte) {
	r.cli.Set(ctx, id, data, r.ttl)
}

// NewRedis returns a Redis store client to be used by a rate limiter.
func NewRedis(cli *redis.Client, ttl time.Duration) *Redis {
	return &Redis{
		cli: cli,
		ttl: ttl,
	}
}

package ratelimiter

import (
	"github.com/althk/ratelimiter/store"
	"github.com/althk/ratelimiter/types"
	"github.com/redis/go-redis/v9"
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

// LimiterOptions allows selection and configuration of the rate limiter
// and other related dependencies.
type LimiterOptions struct {
	Algo                  LimiterAlgo
	StoreType             LimiterStore
	ReqLimit              int64
	TokenBucketRate       float32
	SlidingWindowDuration time.Duration
	RedisOpts             *redis.Options
}

// Build builds a limiter with the given algorithm and storage driver.
func Build(opts *LimiterOptions) types.Limiter {
	var s types.Store
	var l types.Limiter
	switch opts.StoreType {
	case Local:
		s = store.NewInMemory(time.Minute)
	case Redis:
		s = store.NewRedis(redis.NewClient(opts.RedisOpts), time.Minute)
	}
	switch opts.Algo {
	case TokenBucket:
		l = NewTBL(opts.TokenBucketRate, opts.ReqLimit, s)
	case SlidingWindow:
		l = NewSWL(opts.SlidingWindowDuration, opts.ReqLimit, s)
	}
	return l
}

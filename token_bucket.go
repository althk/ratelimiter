package ratelimiter

import (
	"context"
	"github.com/althk/ratelimiter/types"
	"github.com/goccy/go-json"
	"math"
	"time"
)

// tokenBucket represents one token bucket for a client, providing the core token bucket
// rate limiting functionality.
type tokenBucket struct {
	// RefillRate is the rate at which the bucket is refilled in tokens per second.
	RefillRate float32

	// Cap is the max capacity of the bucket. It is also the max burst rate.
	Cap int64

	// Available is the number of tokens currently Available in the bucket.
	Available int64

	// LastRefillTime holds the time the bucket was last refilled at.
	LastRefillTime time.Time
}

// refill calculates the amount of new tokens to be added to the
// bucket since last refill and adds them to the bucket.
// This method must be called from within a synchronised method.
func (tb *tokenBucket) refill(t time.Time) {
	diff := t.Sub(tb.LastRefillTime)
	newTokens := diff.Seconds() * float64(tb.RefillRate)
	tb.Available = int64(math.Min(float64(tb.Available)+newTokens, float64(tb.Cap)))
	tb.LastRefillTime = t
}

func (tb *tokenBucket) allowN(t time.Time, n int64) bool {
	tb.refill(t)
	if tb.Available >= n {
		tb.Available = tb.Available - n
		return true
	}
	return false
}

// TBL implements the token bucket rate limiting algorithm.
// It provides a thread-safe implementation of the Limiter interface.
type TBL struct {
	// refillRate is the rate at which the bucket is refilled in tokens per second.
	refillRate float32

	// cap is the max capacity of the bucket. It is also the max burst rate.
	cap int64

	s types.Store

	initTime time.Time
}

// assert interface implementation
var _ types.Limiter = new(TBL)

func (tbl *TBL) Allow(ctx context.Context, id string) bool {
	return tbl.AllowN(ctx, id, time.Now(), 1)
}

func (tbl *TBL) AllowN(ctx context.Context, id string, t time.Time, n int64) bool {
	b, exists := tbl.s.Get(ctx, id)
	var tb = &tokenBucket{}
	if !exists {
		tb = newTokenBucket(tbl.refillRate, tbl.cap, tbl.initTime)
	} else {
		err := json.Unmarshal(b, tb)
		if err != nil {
			return false
		}
	}

	res := tb.allowN(t, n)

	data, err := json.Marshal(tb)
	if err != nil {
		return false
	}
	tbl.s.Set(ctx, id, data)
	return res
}

// NewTBL returns a new TBL instance which implements the token bucket algorithm
// with the given rate and capacity.
// The new bucket starts with zero Available tokens.
func NewTBL(rate float32, cap int64, s types.Store) *TBL {
	return &TBL{
		refillRate: rate,
		cap:        cap,
		s:          s,
		initTime:   time.Now(),
	}
}

func newTokenBucket(refillRate float32, cap int64, initTime time.Time) *tokenBucket {
	return &tokenBucket{
		RefillRate:     refillRate,
		Cap:            cap,
		LastRefillTime: initTime,
		Available:      0,
	}
}

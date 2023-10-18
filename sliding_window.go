package ratelimiter

import (
	"context"
	"encoding/json"
	"github.com/althk/ratelimiter/types"
	"sync"
	"time"
)

// window represents one time period and holds the number of requests in that time period.
type window struct {
	Start time.Time
	Count int64
}

type slidingWindow struct {
	Curr window
	Prev window
}

func (sw *slidingWindow) slide(d time.Duration, t time.Time) {
	// calculate the newest window Start
	newStart := t.Truncate(d)

	// calculate how many windows our existing current window is from this new window
	diff := newStart.Sub(sw.Curr.Start).Seconds() / d.Seconds()
	// diff < 1 means our existing current window is still
	// the newest window
	if diff < 1 {
		return
	}

	newPrevCount := int64(0)
	// diff < 2 means our current window has become the previous window,
	// so carry fwd the Count
	if diff == 1.000 {
		newPrevCount = sw.Curr.Count
	}

	sw.Curr = window{Start: newStart}
	sw.Prev = window{Start: newStart.Add(-d), Count: newPrevCount}
}

func (sw *slidingWindow) allowN(d time.Duration, t time.Time, n, limit int64) bool {
	sw.slide(d, t)

	elapsed := t.Sub(sw.Curr.Start)
	wt := elapsed.Seconds() / d.Seconds()
	c := int64(wt*float64(sw.Prev.Count)) + sw.Curr.Count

	if c+n > limit {
		return false
	}
	sw.Curr.Count = sw.Curr.Count + n
	return true
}

// SWL implements the sliding window rate limiting algorithm.
// It implements the Limiter interface.
type SWL struct {
	limit          int64
	windowDuration time.Duration
	mu             sync.Mutex
	s              types.Store
}

var _ types.Limiter = new(SWL)

func (swl *SWL) Allow(ctx context.Context, id string) bool {
	return swl.AllowN(ctx, id, time.Now(), 1)
}

func (swl *SWL) AllowN(ctx context.Context, id string, t time.Time, n int64) bool {
	b, exists := swl.s.Get(ctx, id)
	var w = &slidingWindow{}
	if !exists {
		w = swl.newSlidingWindow(t)
	} else {
		err := json.Unmarshal(b, w)
		if err != nil {
			return false
		}
	}

	res := w.allowN(swl.windowDuration, t, n, swl.limit)

	data, err := json.Marshal(w)
	if err != nil {
		return false
	}
	swl.s.Set(ctx, id, data)
	return res
}

func (swl *SWL) newSlidingWindow(t time.Time) *slidingWindow {
	start := t.Truncate(swl.windowDuration)
	return &slidingWindow{Curr: window{Start: start},
		Prev: window{Start: start.Add(-swl.windowDuration)},
	}
}

// NewSWL creates a new instance of SWL, a sliding window rate limiter
// with the given window duration and the limit.
// The window Start is always rounded down to the nearest multiple of duration.
func NewSWL(d time.Duration, l int64, s types.Store) *SWL {
	return &SWL{
		limit:          l,
		windowDuration: d,
		s:              s,
	}
}

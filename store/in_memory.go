package store

import (
	"context"
	"github.com/althk/ratelimiter/types"
	"sync"
	"time"
)

type entry struct {
	data     []byte
	lastSeen time.Time
}
type InMemory struct {
	cMap   map[string]*entry // TODO: use a distributed map
	mu     sync.RWMutex
	maxAge time.Duration
}

var _ types.Store = new(InMemory)

func (m *InMemory) Get(_ context.Context, id string) ([]byte, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.cMap[id]
	if !ok {
		return nil, false
	}
	c.lastSeen = time.Now()
	return c.data, ok
}

func (m *InMemory) Set(_ context.Context, id string, data []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cMap[id] = &entry{data: data, lastSeen: time.Now()}
}

// removeOld removes entries which are older than the given maxAge.
// TODO: use a proper structure to clean up entries instead of looping over the entire keyset
func (m *InMemory) removeOld(maxAge time.Duration) {
	t := time.NewTicker(15 * time.Second)
	for {
		<-t.C
		// using read lock to not starve actual operations
		m.mu.RLock()
		for k, c := range m.cMap {
			if time.Since(c.lastSeen) > maxAge {
				delete(m.cMap, k)
			}
		}
		m.mu.RUnlock()
	}
}

// NewInMemory creates a new in-memory rate limiter that purges entries older than maxAge.
// The purge happens every 15 secs, so the true max age for an entry could be at most 1m15s
func NewInMemory(maxAge time.Duration) *InMemory {
	rl := &InMemory{cMap: make(map[string]*entry), maxAge: maxAge}
	// go rl.removeOld(maxAge)
	return rl
}

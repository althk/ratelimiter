package store

import (
	"context"
	"github.com/althk/ratelimiter/types"
)

// Test provides a primitive implementation of Store
// to be used only for tests/local dev.
type Test struct {
	cMap map[string][]byte
}

var _ types.Store = new(Test)

func (t *Test) Get(ctx context.Context, id string) ([]byte, bool) {
	v, ok := t.cMap[id]
	return v, ok
}

func (t *Test) Set(ctx context.Context, id string, data []byte) {
	t.cMap[id] = data
}

func NewTesting() *Test {
	return &Test{cMap: make(map[string][]byte)}
}

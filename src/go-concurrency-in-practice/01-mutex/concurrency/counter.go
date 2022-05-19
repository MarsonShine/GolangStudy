package concurrency

import "sync"

type Counter struct {
	mux   sync.Mutex
	count int64
}

func (c *Counter) Increament() {
	c.mux.Lock()
	c.count++
	c.mux.Unlock()
}
func (c *Counter) Count() int64 {
	return c.count
}

func NewCounter(n int64) *Counter {
	return &Counter{
		count: n,
	}
}

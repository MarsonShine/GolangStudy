package main

import (
	"sync"
	"time"
)

type Counter struct {
	mu    sync.RWMutex
	count int64
}

func (c *Counter) Count() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}
func (c *Counter) Increament() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

func main() {
	counter()
}

func counter() {
	var counter Counter
	for i := 0; i < 10; i++ {
		go func() {
			for {
				counter.Count()
				time.Sleep(time.Millisecond)
			}
		}()
	}

	for {
		counter.Increament() // 些操作
		time.Sleep(time.Second)
	}

}

/*
RWMutex 读写所，适用于读多写少的场景。
读操作允许并行。一旦发生些操作，所有操作在写锁释放之前一直阻塞；
在尝试持有写锁时，如果发现存在没有释放的读锁，那么在释放这些读锁之前，会一直等待
*/

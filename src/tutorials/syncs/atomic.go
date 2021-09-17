package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 利用原子操作做到无锁编程

func atomicCounter() {
	var ops uint64

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 1000; j++ {
				// 原字操作这样就不用加锁来达到最高性能
				atomic.AddUint64(&ops, 1)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("ops:", ops)
}

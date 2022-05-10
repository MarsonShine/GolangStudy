package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Printf("%d ", i)
		}()
	}

	// v2
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("%d ", i)
		}(i)
	}

	sign := make(chan int, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("%d ", i)
			sign <- i
		}(i)
	}

	for i := 0; i < 10; i++ {
		fmt.Printf("序列 %d 已完成", <-sign)
	}
	time.Sleep(time.Second * 1)

	var count uint32

	trigger := func(i uint32, fn func()) {
		for {
			if n := atomic.LoadUint32(&count); n == i {
				fn()
				atomic.AddUint32(&count, 1)
				break
			}
			time.Sleep(time.Nanosecond)
		}
	}
	// 按顺序执行
	for i := uint32(0); i < 10; i++ {
		go func(i uint32) {
			fn := func() {
				fmt.Printf("%d ", i)
			}
			trigger(i, fn)
		}(i)
	}

	trigger(10, func() {})
}

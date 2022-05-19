package main

import (
	"fmt"
	"go-cip/mutex/concurrency"
	"sync"
)

func main() {
	// 计数器
	counter() // 非并发安全，每次运行结果不一样
	counter_safe()
	counter_concurrency()
}

func counter() {
	var count = 0
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				count++
			}
		}()
	}
	wg.Wait()
	fmt.Println(count)
}

func counter_safe() {
	var mux sync.Mutex
	var count = 0
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				mux.Lock()
				count++
				mux.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Println(count)
}

func counter_concurrency() {
	var wg sync.WaitGroup
	counter := concurrency.NewCounter(0)
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				counter.Increament()
			}
		}()
	}
	wg.Wait()
	fmt.Println(counter.Count())
}

// doc:
/*
可以用 go run -race main.go 来检测是否有并发情况.
启用了-race命令，会影响程序的性能。可以通过生成中间代码查看开启race检测之后，生成的编译代码：go tool compile -race -S main.go
*/

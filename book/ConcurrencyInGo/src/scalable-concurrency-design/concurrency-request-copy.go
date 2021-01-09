package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 请求并发复制处理
// 对于多个并行的任务，如果其中一个最先响应任务，并立即返回结果并有效的取消其它全部未结束的操作
// 通过 10 个处理程序复制模拟请求
func concurrencyRequestCopyExample() {
	doWork := func(done <-chan interface{}, id int, wg *sync.WaitGroup, result chan<- int) {
		start := time.Now()
		defer wg.Done()

		// 模拟随机加载消耗时间
		simulcatedLoadTime := time.Duration(1*rand.Intn(5)) * time.Second
		select {
		case <-done:
		case <-time.After(simulcatedLoadTime):
		}

		select {
		case <-done:
		case result <- id:
		}

		took := time.Since(start)
		// 显示处理程序将花费多少时间
		if took < simulcatedLoadTime {
			took = simulcatedLoadTime
		}
		fmt.Printf("%v took %v\n", id, took)
	}

	done := make(chan interface{})
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go doWork(done, i, &wg, result)
	}

	firstReturned := <-result
	close(done)
	wg.Wait()

	fmt.Printf("Received an answer from #%v\n", firstReturned)
}

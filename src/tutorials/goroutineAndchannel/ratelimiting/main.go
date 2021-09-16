package main

import (
	"fmt"
	"time"
)

// 速率限制是控制服务资源利用和质量的重要机制。基于协程、通道和打点器（timer.Ticker）实现的

func main() {
	requests := make(chan int, 5)
	for i := 0; i < 5; i++ {
		requests <- i
	}
	close(requests)

	limiter := time.Tick(200 * time.Millisecond)

	for req := range requests {
		// 每200ms执行，能控制每次请求控制在200ms内
		<-limiter
		fmt.Println("request", req, time.Now())
	}

	// 通过构建3的缓冲通道来实现
	// 允许最多同事3个并发请求
	burstyLimiter := make(chan time.Time, 3)

	// 在限制的速率下允许并行三个
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	// 每200ms添加新的值插入缓冲通道中
	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	// 模拟5个请求并发执行
	burstyRequests := make(chan int, 5)
	for i := 0; i < 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)
	for req := range burstyRequests {
		// 因为burstyLimiter已经存在3个等待的发送者，所以允许并发处理3个
		// 而后就是控制200ms
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}
}

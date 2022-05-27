package main

import (
	"log"
	"time"

	"golang.org/x/sync/singleflight"
)

// 扩展并发原语
// SingleFlight 和 CyclicBarrier
// SingleFlight 处理多个goroutine同时调用同一个函数时，只让一个goroutine去调用这个函数。
// 等到这个结果返回时，在把结果返回给这几个同时调用的goroutine，这样可以减少并发数

// CyclicBarrier 循环栅栏
// 常用于重复进行一组goroutine同时执行的场景中
// CyclicBarrier允许一组goroutine批次等待，到达一个共同的执行点。同时因为可以被重用，所以叫循环栅栏
// 具体机制是：大家都在栅栏前等待，等全部都到齐了，就抬起栅栏放行
// 多适合用在固定数量的 goroutine 等待同一个执行点

func main() {
	singleFlight()
}

func singleFlight() {
	var singleSetCache singleflight.Group
	getAndSetCache := func(requestId int, cacheKey string) (string, error) {
		log.Printf("request %v start to get and set cache...", requestId)
		value, _, _ := singleSetCache.Do(cacheKey, func() (interface{}, error) {
			log.Printf("request %v is setting cache...", requestId)
			time.Sleep(3 * time.Second) // 模拟耗时
			log.Printf("request %v set cache success!", requestId)
			return "VALUE", nil
		})
		return value.(string), nil
	}
	cacheKey := "cacheKey"
	for i := 0; i < 10; i++ { // 模拟并发请求
		go func(requestId int) {
			value, _ := getAndSetCache(requestId, cacheKey)
			log.Printf("request %v get value: %v", requestId, value)
		}(i)
	}
	time.Sleep(20 * time.Second)
}

func cyclicBarrier() {

}

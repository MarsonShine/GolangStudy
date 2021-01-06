package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// fan-out: 用于描述启动多个 goroutines 以处理来自管道的输入的过程
// fan-in: 用于描述将多个结果组合到一个channel中的过程
func fanoutinExample() {
	// 生成最大值为 50000000 的一串随机数
	rand := func() interface{} { return rand.Intn(50000000) }
	done := make(chan interface{})
	defer close(done)
	start := time.Now()
	// 将转换为整数流传递
	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	fmt.Printf("Search took: %v", time.Since(start))
}

func fanoutinExample2() {
	done := make(chan interface{})
	defer close(done)
	start := time.Now()
	rand := func() interface{} { return rand.Intn(50000000) }
	randIntStream := toInt(done, repeatFn(done, rand))
	// CPU 核心数
	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Primes:")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	fmt.Printf("Search took: %v", time.Since(start))
}

var toInt = func(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case intStream <- v.(int):
			}
		}
	}()
	return intStream
}
var primeFinder = func(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
	primeStream := make(chan interface{})
	go func() {
		defer close(primeStream)
		for integer := range intStream {
			integer -= 1
			prime := true
			for divisor := integer - 1; divisor > 1; divisor-- {
				if integer%divisor == 0 {
					prime = false
					break
				}
			}
			if prime {
				select {
				case <-done:
					return
				case primeStream <- integer:
				}
			}
		}
	}()
	return primeStream
}

// fanoutinExample 函数执行效率非常慢
// 这时我们可以采用fan-out，启动多个 goroutine 执行计算
// 再用fan-in，将多个 goroutine 的执行结果汇聚一起到一个 channel 中
// 这里多个 goroutine 我们采用 CPU 的核芯数
var fanIn = func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})
	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}
	// 等待所有数据汇总
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

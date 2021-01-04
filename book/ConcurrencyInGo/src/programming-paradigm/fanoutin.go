package main

import (
	"fmt"
	"math/rand"
	"time"
)

// fan-out: 用于描述启动多个 goroutines 以处理来自管道的输入的过程
// fan-in: 用于描述将多个结果组合到一个channel中的过程
func fanoutinExample() {
	rand := func() interface{} { return rand.Intn(50000000) }
	done := make(chan interface{})
	defer close(done)
	start := time.Now()
	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	fmt.Printf("Search took: %v", time.Since(start))
}

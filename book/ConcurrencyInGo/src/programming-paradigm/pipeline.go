package main

import (
	"fmt"
	"math/rand"
)

// 管道
// 执行将数据输入，对其进行转换并将数据发回
func pipelineExample() {
	// 阶段1
	multipy := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}
	// 阶段2
	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}
	ints := []int{1, 2, 3, 4, 5}
	// 合并
	for _, v := range add(multipy(ints, 2), 1) {
		fmt.Println(v)
	}
	// 如何定义一个（管道）流水线函数
	// 一个阶段消耗并返回相同类型，
	// 一个阶段必须通过语言来体现，以便它可以被传递（传递函数），与函数式编程思想一致
	// 可以将在加上一个阶段
	for _, v := range multipy(add(multipy(ints, 2), 1), 2) {
		fmt.Println(v)
	}
	// 使用管道的好处之一是能够同时处理数据的各个阶段
}

// 用channel改造pipelineExample
// 虽然代码变多了，但是并发性变强了
func pipelineExample2() {
	// 传递 done，是为了防止 goroutine 泄露，能够安全的退出
	generator := func(done <-chan interface{}, intergers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range intergers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multipy := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4, 5)
	pipeline := multipy(done, add(done, multipy(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}

// 这个函数会重复你传给它的值，直到你告诉它停止
var repeat = func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

// 这个函数会从其传入的valueStream中取出第一个元素然后退出
var take = func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

func pipelineExample3() {
	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}
}

func pipelineExample4() {
	// 创建一个重复调用的生成函数
	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}
	// 生成10个随机数
	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} {
		return rand.Int()
	}

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}
}

func pipelineExample5() {
	toString := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case stringStream <- v.(string):
				}
			}
		}()
		return stringStream
	}

	// useage
	done := make(chan interface{})
	defer close(done)
	var message string
	for token := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
		message += token
	}
	fmt.Printf("message: %s...", message)
}

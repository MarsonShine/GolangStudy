package main

import (
	"fmt"
	"math"
	"sync"
)

/*
Pipeline 这种技术在可以很容易的把代码按单一职责的原则拆分成多个高内聚低耦合的小模块，然后可以很方便地拼装起来去完成比较复杂的功能。
*/

// channel 转发函数
func echo(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// 平方函数
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// 过滤奇数函数
func odd(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if n%2 != 0 {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

// 求和函数
func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var sum = 0
		for n := range in {
			sum += n
		}
		close(out)
	}()
	return out
}

func main() {
	var nums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for n := range sum(sq(odd(echo(nums)))) {
		fmt.Println(n)
	}
	// 如果不像嵌套这么多函数，可以通过下面的方式实现
	// 但是我个人觉得影响了阅读和代码理解
	for n := range pipeline(nums, echo, odd, sq, sum) {
		fmt.Println(n)
	}

	// 分段并发求和，高效
	nums = makeRange(1, 10000)
	in := echo(nums)

	const nProcess = 5
	var chans [nProcess]<-chan int
	for i := range chans {
		chans[i] = sum(prime(in))
	}
	// 合并
	for n := range sum(merge(chans[:])) {
		fmt.Println(n)
	}
}
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + 1
	}
	return a
}

type EchoFunc func([]int) <-chan int
type PipeFunc func(<-chan int) <-chan int

func pipeline(nums []int, echo EchoFunc, pipeFuncs ...PipeFunc) <-chan int {
	ch := echo(nums)
	for i := range pipeFuncs {
		ch = pipeFuncs[i](ch)
	}
	return ch
}

func is_prime(value int) bool {
	for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func prime(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if is_prime(n) {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

func merge(cs []<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan int) {
			for n := range c {
				out <- n
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

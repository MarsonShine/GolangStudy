package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// WaitGroup 最佳实践：先统一Add，再并发Done，最后Wait
func main() {
	// coordinateWithWaitGroup()
	coordinateWithWaitGroup2()
	coordinateWithWaitContext()
}

func coordinateWithWaitGroup() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	num := int32(0)
	fmt.Printf("The number: %d [with sync.WaitGroup]\n", num)
	max := int32(10)
	go addNum(&num, 3, max, wg.Done)
	go addNum(&num, 3, max, wg.Done)
	wg.Wait()
}

func addNum(n *int32, id, max int32, deferFunc func()) {
	defer func() {
		deferFunc()
	}()
	for i := 0; ; i++ {
		currNum := atomic.LoadInt32(n)
		if currNum >= max {
			break
		}
		newNum := currNum + 2
		time.Sleep(time.Millisecond * 200)
		if atomic.CompareAndSwapInt32(n, currNum, newNum) {
			fmt.Printf("The number: %d [%d-%d]\n", newNum, id, i)
		} else {
			fmt.Printf("The CAS operation failed. [%d-%d]\n", id, i)
		}
	}
}

func coordinateWithWaitGroup2() {
	total := 12
	stride := 3
	var num int32
	fmt.Printf("The number: %d [with sync.WaitGroup]\n", num)
	var wg sync.WaitGroup
	for i := 1; i <= total; i = i + stride {
		wg.Add(stride)
		for j := 0; j < stride; j++ {
			go addNum2(&num, i+j, wg.Done)
		}
		wg.Wait()
	}
	fmt.Println("End.")
}

func addNum2(n *int32, id int, deferFunc func()) {
	defer func() {
		deferFunc()
	}()
	for i := 0; ; i++ {
		currNum := atomic.LoadInt32(n)
		newNum := currNum + 1
		time.Sleep(time.Millisecond * 200)
		if atomic.CompareAndSwapInt32(n, currNum, newNum) {
			fmt.Printf("The number: %d [%d-%d]\n", newNum, id, i)
			break
		} else {
			//fmt.Printf("The CAS operation failed. [%d-%d]\n", id, i)
		}
	}
}

// 用context.Context实现与sync.WaitGroup一样的效果
func coordinateWithWaitContext() {
	total := 12
	var num int32
	fmt.Printf("The number: %d [with context.Context]\n", num)
	cxt, cancelFunc := context.WithCancel(context.Background())
	for i := 1; i <= total; i++ {
		go addNum2(&num, i, func() {
			if atomic.LoadInt32(&num) == int32(total) {
				cancelFunc() // 调用context.WithCancel，该goroutine执行完就会调用cxt.Done()
			}
		})
	}
	<-cxt.Done() // 通道结束
	// 从而实现了等待所有的 addNum2 函数都执行完毕
	fmt.Println("End.")
}

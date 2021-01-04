package main

import (
	"fmt"
	"time"
)

// 将一个或多个done通道合并到一个done通道中
// 这个通道任何组件关闭时，这个通道就会关闭
// 这种模式使用递归和 goroutine 创建一个复合done通道
var or func(channels ...<-chan interface{}) <-chan interface{}

func orchannelExample() {
	// 递归
	or = func(channels ...<-chan interface{}) <-chan interface{} { // 1: 接收可变数量的通道并返回一个channel
		switch len(channels) {
		case 0: // 2: 设置终止条件
			return nil
		case 1: // 3: 设置终止条件
			return channels[0]

		}
		orDone := make(chan interface{})
		go func() { // 4 通过开启协程可以不受阻塞的等待channel上的结果
			defer close(orDone)
			switch len(channels) {
			case 2: // 5
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default: // 6
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
					fmt.Println("channels[2] 完成")
				case <-or(append(channels[3:], orDone)...): // 6 传递 orDone通过，这样当该树状结构顶层的 goroutine 退出时，结构底层的 goroutines 也会退出
				}
			}
		}()
		return orDone
	}

	// useage
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(sig(2*time.Hour), sig(5*time.Minute), sig(1*time.Second), sig(1*time.Hour), sig(1*time.Minute)) // 将多个channel组合并返回单个channel
	fmt.Printf("done after %v", time.Since(start))
}

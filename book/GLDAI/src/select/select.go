// select 类似于C语言的 select 函数，监听多个文件描述符的可读可写状态
// go 的 select 关键字能够让 goroutine 同时等待多个 Channel 的可读可写状态
// 在多个文件或 Channel 发生状态改变之前，select 会一直阻塞当前线程或 Gorountine
// select 用法与 switch 相似，与 case 块一起使用
// 无论满足哪一个 case 都会立即执行相应 case 下的语句代码块。
// 如果多个 case 都满足，则会随机选择一个 case 语句块执行

// 在 go 中 select 与 channel 一起运行可以进行非阻塞的收发操作
// select 在遇到多个 channel 同时响应会随机选择一个 case 执行
// 如何用 select 实现非阻塞的收发？见 #1
// 当遇到同时满足多个 case 的情况？见 #2
package main

import (
	"fmt"
	"time"
)

func main() {
	// #1
	// 用一个默认的 case 来实现这种情况：当遇到指定的 case 会立即执行 case 语句块，如果没有满足任意的 case，则会执行 default 语句块，这样就不会阻塞了，继续监听下一个事件
	ch := make(chan int)
	// select {
	// case i := <-ch:
	// 	println(i)
	// default:
	// 	println("default")
	// }
	// // 非阻塞的收发也可以使用另一种形式
	// x, ok := <-ch
	// if ok {
	// 	println(x)
	// } else {
	// 	println("default")
	// }
	// #2
	go func() {
		for range time.Tick(1 * time.Second) {
			ch <- 0
		}
	}()

	for {
		select {
		case <-ch:
			println("case1")
		case <-ch:
			println("case2")
		}
	}
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

package main

import "fmt"

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	// go func() {
	// 	for x := 0; ; x++ {
	// 		naturals <- x
	// 	}
	// }()

	// Squarer
	// go func() {
	// 	for {
	// 		x := <-naturals // 等待 Counter 协程处理程序的 naturals <- x 执行完毕
	// 		squares <- x * x
	// 	}
	// }()
	// 两个协程都是一个死循环，如果没有数据发送了，可以通知接收者停止无限读取发送的值
	// Printer (in main goroutine)
	// for {
	// 	fmt.Println(<-squares)
	// }
	go func() {
		for x := 0; x < 100; x++ {
			// 发送
			naturals <- x
		}
		close(naturals)
	}()
	go func() {
		for x := range naturals { // 消费
			squares <- x * x
		}
		close(squares)
	}()

	// Printer (in main goroutine)
	for x := range squares { // 消费
		fmt.Println(x)
	}

}

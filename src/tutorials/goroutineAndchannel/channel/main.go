package main

import (
	"fmt"
	"time"
)

// 一个 channel 就是一个通信机制，是一个 goroutine 与另一个 goroutine 通信的通道。
// 每个 channel 都有一个类型，是发送数据的类型，如发送 int 的channel 就是 chan int
// 发送和接收都是用 <- 操作符
// 发送：ch <- x；x发送至 channel
// 接收: <-ch / x <-ch；channel 接收数据赋值给弃元或者接收并赋值给 x

func main() {
	ch := make(chan int) // make 创建的channel 是一个引用，底层引用了数据结构，所以在传递 channel 都是传递引用。
	// channel 有发送和接收两个操作，都是通信行为
	ch = make(chan int, 3) // 初始化 channel 的容量，容量大于0 即外表有缓存
	// 无缓存 channel 又叫同步channel。发送数据之后 发送者的 goroutine 是阻塞的，直到这个发送的数据被另一个端接收完成才继续执行发送者发送数据之后的goroutine
	// 接收者也是如此
	<-ch

	syncAndNoBufferChannel()

	syncChannel()
}

// 无缓冲 Channel，即同步 Channel
func syncAndNoBufferChannel() {
	messages := make(chan string)

	go func() {
		messages <- "ping"
	}()

	fmt.Println(<-messages)
}

// 缓冲 Channel
func bufferChannel() {
	messages := make(chan string, 2)
	messages <- "buffered"
	messages <- "channel"
	fmt.Println(<-messages)
	fmt.Println(<-messages)
}

func syncChannel() {
	done := make(chan bool, 1)
	go worker(done)
	fmt.Println("主线程等待...")
	fmt.Println("结束等待，done")
}

func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second * 5)
	fmt.Print("done")

	done <- true
}

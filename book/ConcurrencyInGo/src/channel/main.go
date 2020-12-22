package main

import (
	"fmt"
	"sync"
)

func main() {
	// singleChannelExample()
	// channelExample()
	channelExample2()
}

// 单向甬道数据流
// 单向通道常被用做入参和出参
// 可以将双向通道转换成单向通道
func singleChannelExample() {
	// 申明一个只能接收的数据流
	var receiveDataStream <-chan interface{}
	receiveDataStream = make(<-chan interface{})
	// 申明一个只能发送的数据流
	var sendDataStream chan<- interface{}
	sendDataStream = make(chan<- interface{})
	// 双向通道
	dataStream := make(chan interface{})
	// 双变单
	receiveDataStream = dataStream
	sendDataStream = dataStream

	<-receiveDataStream
	sendDataStream <- 3
}

func channelExample() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hello channels."
	}()
	fmt.Println(<-stringStream)
}

// 在读取接收channel的值之前，如果通道内没有元素发送就会一直阻塞
func channelDeatlockExample() {
	stringStream := make(chan string)
	// 关闭通道会节省资源，多个 goroutine 执行也会越快
	defer close(stringStream)
	go func() {
		for 0 != 1 {
			return
		}
		stringStream <- "Hello channels"
	}()
	fmt.Println(<-stringStream)
}

// channel 可以与 for range 模式一起使用
// 会在通道关闭自动结束循环
func channelExample2() {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 1; i < 5; i++ {
			intStream <- i
		}
	}()
	for intger := range intStream {
		fmt.Printf("%v ", intger)
	}
}

func channelExample3() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%v has begun\n", i)
		}(i)
	}
	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Done()
}

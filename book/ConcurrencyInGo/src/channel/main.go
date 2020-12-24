package main

import (
	"bytes"
	"fmt"
	"os"
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

// 缓冲通道如果是空的，但是此时有接收器正在接收，那么这个缓冲器就会被绕过，发送器直接将数据发送给接收器以避免阻塞。
// channel 读取注意事项:
// 读：channel 为 nil 时，会阻塞，在只写（发送）状态的 channel 时会编译失败
// 写：channel 为 nil 时，会阻塞，在关闭时，会报错；在接收状态的 channel 会编译失败
// 关闭：nil 时，会报错；在关闭时会报错；在只接收状态下，编译失败
func channelExample4() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()
	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %v.\n", integer)
	}
}

// 如何判断一个 channel 是否被关闭？
// 如何才能知道是因为什么被阻塞？
func channelExample5() {
	//强烈建议在自己的程序中尽可能做到保持通道覆盖范围最小，以便这些事情保持明显。如果你将一个通道作为一个结构体的成员变量，并且有很多方法，它很快就会把你自己给绕进去
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5) // 1
		go func() {                       // 2
			defer close(resultStream) // 3
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream // 4
	}
	resultStream := chanOwner()
	for result := range resultStream { // 5
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
}

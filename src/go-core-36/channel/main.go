package main

import (
	"fmt"
	"time"
)

func main() {
	sigal := make(chan int, 1)
	channelRule()

	// ch := make(chan int, 2)

	// go func() {
	// 	for {
	// 		select {
	// 		case <-ch:
	// 			fmt.Println("接收数据...")
	// 		default:
	// 			fmt.Println("通道元素为空，等待发送端发送数据")
	// 		}
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		select {
	// 		case ch <- 0:
	// 			fmt.Println("发送端已发送数据")
	// 		default:
	// 			fmt.Println("通道元素已满，等待接收端接收数据")
	// 		}
	// 	}
	// }()

	<-sigal
}

func timeout() {
	ch := make(chan int, 1)
	time.AfterFunc(time.Second, func() {
		close(ch)
	})
	select {
	case _, ok := <-ch:
		if !ok {
			fmt.Println("The candidate case is closed.")
			break
		}
		fmt.Println("The candidate case is selected.")
	}
}

var channels = [3]chan int{
	nil,
	make(chan int),
	nil,
}

var numbers = []int{1, 2, 3}

func channelRule() {
	select {
	case getChan(0) <- getNumber(0):
		fmt.Println("The first candidate case is selected.")
	case getChan(1) <- getNumber(1):
		fmt.Println("The second candidate case is selected.")
	case getChan(2) <- getNumber(2):
		fmt.Println("The third candidate case is selected")
	default:
		fmt.Println("No candidate case is selected!")
	}
}
func getNumber(i int) int {
	fmt.Printf("numbers[%d]\n", i)
	return numbers[i]
}

func getChan(i int) chan int {
	fmt.Printf("channels[%d]\n", i)
	return channels[i]
}

package main

import (
	"fmt"
	"time"
)

func main() {
	// selectExample()
	// selectExample2()
	selectExample4()
}

// select 如果满足 case 条件，就会执行；否则会一直阻塞
func selectExample() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()

	fmt.Println("Blocking on read...")
	select {
	case <-c: // 2 因为没有发送值，所以也就不存在接收，即一直阻塞
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

// 多个通道读取
func selectExample2() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)
	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
	// 每次运行的结构几乎都是各占一半，这是因为 Go 运行时对此做了伪随机操作
	// 在同时都满足 case 的情况会伪随机执行
}

// 通道没有初始化都开始读取，所以会无限阻塞下去
// 所以可以设置一个超时时间
func selectExample3() {
	var c <-chan int
	select {
	case <-c: //1 读取 nil，会无限阻塞
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}
}

// 如果想给 select 一个默认处理的条件
func selectExample4() {
	start := time.Now()
	var c1, c2 <-chan int
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}
}

// for + select
// 实现在逻辑中可以一边等待其它 goroutine 的结果
// 也可以执行自己的逻辑
func selectExample5() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()
	workCounter := 0

loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		workCounter++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
}

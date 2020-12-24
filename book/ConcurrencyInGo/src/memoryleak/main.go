package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	memoryleak2()
	notmemoryleak2()
}

func memoryleak() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings {
				fmt.Println(s)
			}
		}()
		return completed
	}
	// 传递一个 nil 的 channel
	// 所以读取这个 channel 会一直阻塞
	// 而且包含 doWork 的 goroutine 将在这个过程的整个生命周期中保留在内存中
	doWork(nil)
	// 这里还有其他任务执行
	fmt.Println("Done.")
}

func notmemoryleak() {
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		//1
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					// Do something interesting
					fmt.Println(s)
				case <-done: //2
					return
				}
			}
		}()
		return terminated
	}
	done := make(chan interface{})
	terminated := doWork(done, nil)
	go func() { //3
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()
	<-terminated //4
	fmt.Println("Done.")
}

// 给正在阻塞的 channel 写入值的情况
func memoryleak2() {
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for { // 这里一旦触发 newRannStream 就会无限循环写入数据至通道
				randStream <- rand.Int()
			}
		}()
		return randStream
	}

	randStream := newRandStream()
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ { // 这里只接收三次数据，所以 randStream 函数内的写入通道的 goroutine 还是一直会存在，所占用的内存也不会被回收
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
}

// 这里就用另一个只接收的通道来通知取消，来防止内存泄漏
func notmemoryleak2() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for { // 这里一旦触发 newRannStream 就会无限循环写入数据至通道
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done) // 通知流程已经结束
	// 模式耗时
	time.Sleep(1 * time.Second)
}

package main

import (
	"fmt"
	"time"
)

func main() {

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

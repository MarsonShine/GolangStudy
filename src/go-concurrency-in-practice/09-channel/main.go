package main

import (
	"fmt"
	"go-cip/09-channel/mutex"
	"os"
	"os/signal"
	"reflect"
	_ "runtime"
	"syscall"
	"time"
)

func main() {
	// 4个goroutine，分别按顺序输出 1，2，3，4
	method1()
	method2()

	// 优雅退出
	gracefulStop()

	m := mutex.NewMutex()
	ok := m.TryLock()
	fmt.Printf("locked v %v\n", ok)
	ok = m.TryLock()
	fmt.Printf("locked %v\n", ok)
}

func method1() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	ch3 := make(chan struct{})
	ch4 := make(chan struct{})
	go func(n int) {
		for range ch1 {
			print(n)
			ch2 <- struct{}{}
		}
	}(1)
	go func(n int) {
		for range ch2 {
			print(n)
			ch3 <- struct{}{}
		}
	}(2)
	go func(n int) {
		for range ch3 {
			print(n)
			ch4 <- struct{}{}
		}
	}(3)
	go func(n int) {
		for range ch4 {
			print(n)
			ch1 <- struct{}{}
		}
	}(4)
	ch1 <- struct{}{}
	time.Sleep(10 * time.Second) // 注意，会内存泄漏，只做demo演示
	close(ch1)
	close(ch2)
	close(ch3)
	close(ch4)
}
func method2() {
	chs := []chan Token{make(chan Token), make(chan Token), make(chan Token), make(chan Token)}
	// 创建4个worker
	for i := 0; i < 4; i++ {
		go newWorker(i, chs[i], chs[(i+1)%4])
	}
	// 把令牌发送给第一个worker
	chs[0] <- struct{}{}
	select {}
}

func gracefulStop() {
	// 处理程序中断
	terminalChan := make(chan os.Signal)
	signal.Notify(terminalChan, syscall.SIGINT, syscall.SIGTERM)
	<-terminalChan

	// 执行退出之前的清理动作
	cleanup()
	fmt.Println("优雅退出")
}
func cleanup() {

}

type Token struct{}

func newWorker(id int, ch chan Token, nextCh chan Token) {
	for {
		token := <-ch
		fmt.Println(id + 1)
		time.Sleep(time.Second)
		nextCh <- token
	}
}

func print(n int) {
	fmt.Println(n)
	time.Sleep(time.Second)
}

func multiChan() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	select {
	case v := <-ch1:
		fmt.Println(v)
	case v := <-ch2:
		fmt.Println(v)
		// case ...
	}
}

// 要点1
func multiChanHanlder() {
	var ch1 = make(chan int, 10)
	var ch2 = make(chan int, 10)
	// 创建SelectCase
	var cases = createCases(ch1, ch2)
	// 执行10次
	for i := 0; i < 10; i++ {
		chosen, recv, ok := reflect.Select(cases)
		if recv.IsValid() {
			fmt.Println("recv:", cases[chosen].Dir, recv, ok)
		} else {
			fmt.Println("send:", cases[chosen].Dir, ok)
		}
	}
}

func createCases(chs ...chan int) []reflect.SelectCase {
	var cases []reflect.SelectCase
	// 创建recv case
	for _, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}
	// 创建send case
	for _, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
		})
	}
	return cases
}

/*
要点1：如何动态的处理多个chan？
首先我们可以通过 select 语句来除了多个chan，有多少个chan就可以有对应数量的case节点。但是如果程序中有几十上百的chan，难道也要一个个写case么？
要解决这个问题，即动态处理多个chan，就可以用到反射了。
*/

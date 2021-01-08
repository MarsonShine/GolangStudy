package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 基于固定时间进行心跳检查
func heartbeatExample() {
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{}) // 1 创建一个心跳的通道
		results := make(chan time.Time)

		go func() {
			defer close(heartbeat)
			defer close(results)

			pulse := time.Tick(pulseInterval)       // 2, 设置时间定时发送心跳
			workGen := time.Tick(2 * pulseInterval) // 3, 这是模拟另一处正在工作的代码，模拟接收心跳信息

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default: // 4, 如果没有接受到心跳也能保证不会阻塞
				}
			}

			sendResult := func(r time.Time) {
				select {
				case <-done:
					return
				case <-pulse: // 5, 接受到心跳请求，发送心跳请求
					sendPulse()
				case results <- r:
					return
				}
			}

			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()

		return heartbeat, results
	}

	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() {
		close(done) // 1, 10秒后触发取消操作
	})
	const timeout = 2 * time.Second               // 2, 设置超时时间
	heartbeat, results := doWork(done, timeout/2) // 3, 传递事件并同时返回心跳和结果

	for {
		select {
		case _, ok := <-heartbeat: // 4 每隔 timeout/2 就获取来自心跳发来的消息
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results: // 5, 接收数据
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout): // 6, 超时，说明没有在规定的时间内响应心跳和发送结果
			return
		}
	}
}

// 循环中断
func heartbeatExample2() {
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		go func() {
			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}
			sendResult := func(r time.Time) {
				for {
					select {
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for i := 0; i < 2; i++ { // 1, 循环两次中断
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}

	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r)
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}

// 在工作任务开始时开启心跳
func heartbeatExample3() {
	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		heartbeatStream := make(chan interface{}, 1) //1
		workStream := make(chan int)
		go func() {
			defer close(heartbeatStream)
			defer close(workStream)
			for i := 0; i < 10; i++ {
				select { //2
				case heartbeatStream <- struct{}{}:
				default: //3
				}
				select {
				case <-done:
					return
				case workStream <- rand.Intn(10):
				}
			}
		}()
		return heartbeatStream, workStream
	}
	done := make(chan interface{})
	defer close(done)
	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}

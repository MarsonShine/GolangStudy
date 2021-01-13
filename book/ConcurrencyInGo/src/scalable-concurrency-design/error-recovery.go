package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// 错误异常恢复
type startGoroutineFn func(done <-chan interface{}, pulseInternal time.Duration) (heartbeat <-chan interface{}) //1 定一个可以监控和重新启动的goroutine的函数

var orDone = func(done, c <-chan interface{}) <-chan interface{} {
	varStream := make(chan interface{})
	go func() {
		defer close(varStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case varStream <- v:
				case <-done:
				}

			}
		}
	}()
	return varStream
}

var bridge = func(
	done <-chan interface{},
	chanStream <-chan <-chan interface{},
) <-chan interface{} {
	valStream := make(chan interface{}) // <1>
	go func() {
		defer close(valStream)
		for { // <2>
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range orDone(done, stream) { // <3>
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

var or func(channels ...<-chan interface{}) <-chan interface{}

// 这个函数会从其传入的valueStream中取出第一个元素然后退出
var take = func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

func errorRecoverySample() {
	or = func(channels ...<-chan interface{}) <-chan interface{} { // <1>
		switch len(channels) {
		case 0: // <2>
			return nil
		case 1: // <3>
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() { // <4>
			defer close(orDone)

			switch len(channels) {
			case 2: // <5>
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default: // <6>
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...): // <6>
				}
			}
		}()
		return orDone
	}

	newSteward := func(timeout time.Duration, startGoroutine startGoroutineFn) startGoroutineFn { // 2 设置超时，并用 startGoroutineFn 来监控，返回也是一个监控函数，就是说这个方法本身也是可以监控的
		return func(done <-chan interface{}, pulseInternal time.Duration) <-chan interface{} {
			heartbeat := make(chan interface{})
			go func() {
				defer close(heartbeat)
				var wardDone chan interface{}
				var wardHeartbeat <-chan interface{}
				startWard := func() { // 3
					wardDone = make(chan interface{})                             // 4
					wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2) // 5
				}
				startWard()
				pulse := time.Tick(pulseInternal)

			monitorLoop:
				for { // 6
					timeoutSignal := time.After(timeout)
					for {
						select {
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						case <-wardHeartbeat: // 7 表示接收到监控着的心跳，说明处于正常工作状态，继续循环监控
							continue monitorLoop
						case <-timeoutSignal: // 8 这里如果我们发现监控者超时，我们要求监控者停下来，并开始一个新的goroutine。然后开始新的检测
							log.Println("steward: ward unhealthy; restarting")
							close(wardDone)
							startWard()
							continue monitorLoop
						case <-done:
							return
						}
					}
				}
			}()
			return heartbeat
		}
	}

	doWorkFn := func(done <-chan interface{}, intList ...int) (startGoroutineFn, <-chan interface{}) { // 1
		intChanStream := make(chan (<-chan interface{})) // 2 创建一个通道的通道
		intStream := bridge(done, intChanStream)

		doWork := func(done <-chan interface{}, pulseInternal time.Duration) <-chan interface{} { // 3 建立闭包控制器的开启和关闭
			intStream := make(chan interface{}) // 4 数据
			heartbeat := make(chan interface{})

			go func() {
				defer close(intStream)
				select {
				case intChanStream <- intStream: // 5 向通道传入数据
				case <-done:
					return
				}

				pulse := time.Tick(pulseInternal)

				for {
				valueLoop:
					for _, intVal := range intList {
						if intVal < 0 {
							log.Printf("negative value: %v\n", intVal) // 6 返回负数并从goroutine返回以模拟不正常的工作状态
							return
						}

						for {
							select {
							case <-pulse:
								select {
								case heartbeat <- struct{}{}:
								default:
								}
							case intStream <- intVal:
								continue valueLoop
							case <-done:
								return
							}
						}
					}
				}
			}()
			return heartbeat
		}
		return doWork, intStream
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	done := make(chan interface{})
	defer close(done)

	doWork, intStream := doWorkFn(done, 1, 2, -1, 3, 4, 5)      // 1
	doWorkWithSteward := newSteward(1*time.Microsecond, doWork) // 2
	doWorkWithSteward(done, 1*time.Hour)                        // 3
	for intVal := range take(done, intStream, 6) {              // 4
		fmt.Printf("Received: %v\n", intVal)
	}

	// doWork := func(done <-chan interface{}, _ time.Duration) <-chan interface{} {
	// 	log.Println("ward: Hello, I'm irresponsible!")
	// 	go func() {
	// 		<-done // 1 持久等待 done 触发
	// 		log.Println("ward: I am halting.")
	// 	}()
	// 	return nil
	// }

	// doWorkWithSteward := newSteward(4*time.Second, doWork) // 2 开启监控程序，4秒后超时
	// done := make(chan interface{})
	// time.AfterFunc(9*time.Second, func() { // 3 9秒超时然后关闭 done
	// 	log.Println("main: halting steward and ward.")
	// 	close(done)
	// })

	// for range doWorkWithSteward(done, 4*time.Second) { // 4

	// }
	log.Println("Done")
}

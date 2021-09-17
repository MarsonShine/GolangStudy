package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var once sync.Once

func main() {
	pool()

	for i := 0; i < 10; i++ {
		once.Do(func() {
			fmt.Println("execed...")
		})
	}

	duplicate(once)
	atomicCounter()
	mockMutexByChannel()
}

func duplicate(once sync.Once) {
	for i := 0; i < 10; i++ {
		once.Do(func() {
			fmt.Println("execed2...")
		})
	}
}

func pool() {
	myPool := &sync.Pool{
		// New: func() interface{} {
		// 	fmt.Println("Creating new instance.")
		// 	return struct{}{}
		// },
	}
	o := myPool.Get() //1
	myPool.Put(1)
	fmt.Println(o)
	// instance := myPool.Get() //1
	// myPool.Put(instance)     //2
	// myPool.Get()             //3
}

func mutexLock() {
	// 共享状态
	var state = make(map[int]int)
	var mutex = &sync.Mutex{}

	var readOps uint64
	var writeOps uint64

	// 开启100个线程并发对共享状态执行写入操作
	for i := 0; i < 100; i++ {
		go func() {
			total := 0
			for {
				key := rand.Intn(5)
				mutex.Lock()
				total += state[key]
				mutex.Unlock()
				atomic.AddUint64(&readOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}
	// 开启10个线程并发循环写入
	for w := 0; w < 10; w++ {
		go func() {
			for {
				key := rand.Intn(5)
				val := rand.Intn(100)
				mutex.Lock()
				state[key] = val
				mutex.Lock()
				atomic.AddUint64(&writeOps, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)

	readOpsFinal := atomic.LoadUint64(&readOps)
	fmt.Println("readOps:", readOpsFinal)
	writeOpsFinal := atomic.LoadUint64(&writeOps)
	fmt.Println("writeOps:", writeOpsFinal)

	mutex.Lock()
	fmt.Println("state:", state)
	mutex.Unlock()
}

func mockMutexByChannel() {
	var state int = 0
	var notify = make(chan int)
	var done = make(chan bool)
	go func() {
		for {
			select {
			case i := <-notify:
				state += i
			case <-done:
				fmt.Println("stop goroutine...")
				return
			}
		}
	}()

	for i := 0; i < 50; i++ {
		go func(k int) {
			notify <- k
		}(1)
	}
	time.Sleep(time.Second)
	done <- true
	fmt.Println("sum:", state)
}

package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"text/tabwriter"
	"time"
)

var wg sync.WaitGroup

func main() {

	// sayHello := func() {
	// 	defer wg.Done()
	// 	fmt.Println("hello")
	// }
	// wg.Add(1)
	// go sayHello()
	// wg.Wait()

	// example2()
	// wg.Wait()	// 判断计数是否等于0，否则继续等待

	// memoryMeasurement()
	// bulkWait()
	// mutexExample()
	// rwmutexExample()
	// condExample()
	// coodExample2()
	// onceDoExample()
	poolExample()
}

func example2() {
	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1) // 增加计数
		go func(sal string) {
			defer wg.Done() // 减少计数
			fmt.Println(sal)
		}(salutation)
	}

}

// 测量一个goroutine的实际内存占用大小。
func memoryMeasurement() {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	noop := func() { wg.Done(); <-c } // 1

	const numGoroutines = 1e4 // 2
	wg.Add(int(numGoroutines))
	before := memConsumed() // 3
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed() // 4
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}

func bulkWait() {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello from %v!\n", id)
	}

	const numGreeters = 5
	var wg sync.WaitGroup
	wg.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg, i+1)
	}
	wg.Wait()
}

func mutexExample() {
	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("Incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("Decrementing: %d\n", count)
	}

	var arithmetic sync.WaitGroup
	for i := 0; i < 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}

	for i := 0; i < 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	arithmetic.Wait()
	fmt.Println("Arithmetic complete.")
}

// 优化 可以用读写锁，只要没有别的 goroutine 对变量进行写操作，就可以任意的读
func rwmutexExample() {
	producer := func(wg *sync.WaitGroup, l sync.Locker) { // 1
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			l.Unlock()
			time.Sleep(1) //2
		}
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}

		wg.Wait()
		return time.Since(beginTestTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutext\tMutex\n")
	for i := 0; i < 25; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw, "%d\t%v\t%v\n", count,
			test(count, &m, m.RLocker()), test(count, &m, &m),
		)
	}
}

// 多个 goroutine 协作，通知其它 goroutine 任务已经完毕
// 通知让各自满足条件的 goroutine 继而往下执行逻辑
// 在没有 cond 之前，想要实现这个目的，需要无限循环判断是否条件：for conditionTrue() == false {}
// 这样会浪费 CPU 的资源（一个循环就是一个 CPU 核芯）
// 作为优化可以显式选择休眠时间，但是同样这个具体数值是不好掌握的：for conditionTrue() == false { time.sleep(1 * time.Millisecond) }
// cond 就能很轻易的做到
func condExample() {
	c := sync.NewCond(&sync.Mutex{}) // 1
	// c.L.Lock()
	// for conditionTrue() == false {
	// 	c.Wait()
	// }
	// c.L.Unlock()

	queue := make([]interface{}, 0, 10) // 2

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()        // 8 锁定 c.L 来进行删除元素操作
		queue = queue[1:] // 9 删除第一个元素
		fmt.Println("Removed from queue")
		c.L.Unlock() // 10 成功删除并解锁
		c.Signal()   // 11 通知其它 goroutine
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()            // 3 锁定 c.L
		for len(queue) == 2 { // 4 循环判断，不满足条件时跳出。因为 removeFromQueue 是异步的，如果把 for 换成 if，则做不到重复判断，只能等下一次迭代在判断，而 for 则是一直重复判断，效率上 for 要更占优
			c.Wait() // 5 调用 wait 阻塞 main goroutine，直到接受信号
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second) // 6 开启新的 goroutine 来执行移除队列元素
		c.L.Unlock()                        // 7 解锁 c.L，因为已经添加元素成功
	}
}

func conditionTrue() bool {
	time.Sleep(1)
	return false
}

func coodExample2() {
	type Button struct {
		// 1
		Clicked *sync.Cond
	}

	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) { // 2
		var tempwg sync.WaitGroup
		tempwg.Add(1)
		go func() {
			tempwg.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		tempwg.Wait()
	}

	var wg sync.WaitGroup // 3
	wg.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window")
		wg.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		wg.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		wg.Done()
	})
	// 在内部，运行时维护一个等待信号发送的goroutines的FIFO列表
	button.Clicked.Broadcast()
	wg.Wait()
}

func onceDoExample() {
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }

	var once sync.Once
	once.Do(increment)
	once.Do(decrement)
	fmt.Println(count)
	// 输出 1 而不是 0
	// 这是因为 once.Do 记录的是调用的次数，而不是传入不同函数的调用次数。
}

func onceDoExample2() {
	var onceA, onceB sync.Once
	var initB func()
	initA := func() { onceB.Do(initB) }
	initB = func() { onceA.Do(initA) }
	onceA.Do(initA)
	// 会发生死锁
}

// 池模式是一种创建和提供固定可用数量对象的方式
// sync.Pool 可以被多个例程安全使用
func poolExample() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance")
			return struct{}{}
		},
	}

	myPool.Get()             // 1 Get 方法会在内部查询是否有可用对象，有直接返回并移除该对象，没有则调用初始化 New
	instance := myPool.Get() // 1
	myPool.Put(instance)     // 2 Put 方法则是将对象返回给对象池
	myPool.Get()             // 3 这个时候再 Get 则不会执行 New 方法函数，因为池中存在一个可用的对象
}

func poolExmaple2() {
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated++
			mem := make([]byte, 1024)
			return &mem
		},
	}
	// 将池扩充到 4 kb
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
		}()
	}
	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}

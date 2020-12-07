// sync.Cond 一个条件变量，它可以让一系列的 Goroutine 都在满足特定条件时被唤醒。
// sync.Cond.Signal 和 sync.Cond.Broadcast 都是用来唤醒调用 sync.Cond.Wait 休眠的 Goroutine
// Cond.Signal：会唤醒队列最前面的 Goroutine
// Cond.Broadcast: 方法会唤醒队列全部的 Goroutine
// 使用注意事项:
// sync.Cond.Wait: 方法在调用之前一定要使用获取互斥锁，否则会触发程序崩溃
// sync.Cond.Signal: 唤醒最先陷入休眠的 goroutine
// sync.Cond.Broadcast: 会按照一定顺序广播通知所有等待的 goroutine
package locks

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

var status int64

func MethodCond() {
	// 每次调用都要传一个互斥锁
	c := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 10; i++ {
		go listen(c)
	}
	time.Sleep(1 * time.Second)
	// go broadcast(c)
	go signalOnce(c)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func broadcast(c *sync.Cond) {
	c.L.Lock()
	atomic.StoreInt64(&status, 1)
	c.Broadcast()
	c.L.Unlock()
}

func signalOnce(c *sync.Cond) {
	c.L.Lock()
	atomic.StoreInt64(&status, 1)
	c.Signal()
	c.L.Unlock()
}

func listen(c *sync.Cond) {
	c.L.Lock()
	for atomic.LoadInt64(&status) != 1 {
		c.Wait() // 将 Goroutine 进入休眠状态
	}
	fmt.Println("listen")
	c.L.Unlock()
}

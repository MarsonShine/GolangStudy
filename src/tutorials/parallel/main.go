package main

import (
	"fmt"
	"sync"
)

// 如果存在多读单写的情况，我们可以用“多读单写”（multiple readers, single writer lock）锁 —— sync.WRMutex

var mu sync.RWMutex // 读写锁
var balance int

func Balance() int {
	mu.RLock()
	defer mu.RUnlock()
	return balance
}

// 内存同步
// 在现代计算机中，每个处理器都会有自己的本地缓存（local cache）。为了效率，对内存的写入一般会缓冲到每个处理器的本地缓存中，在必要时一次性flush到主存。这种情况下这些数据可能会以与当初goroutine写入顺序不同的顺序被提交到主存。像channel通信或者互斥量操作这样的原语会使处理器将其聚集的写入flush并commit，这样goroutine在某个时间点上的执行结果才能被其它处理器上运行的goroutine得到。

func main() {
	var x, y int
	go func() {
		x = 1
		fmt.Print("y:", y, " ") //A2
	}()

	go func() {
		y = 1
		fmt.Print("x:", x, " ")
	}()
}

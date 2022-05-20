package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	c := sync.NewCond(&sync.Mutex{})
	var ready int
	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
			// 加锁等待更改条件
			c.L.Lock()
			ready++
			c.L.Unlock()
			log.Printf("#%d 准备就绪", i)
			c.Broadcast()
		}(i)
	}

	c.L.Lock()        // 这一步非常重要
	for ready != 10 { // 一定要 for condition
		c.Wait()
		log.Println("Broadcast，由于不满足条件，继续等待")
	}
	c.L.Unlock() // 这一步非常重哟啊

	fmt.Printf("全部准备完成")

}

/*
sync.Cond
和某个条件有关，这个条件需要一组goroutine共同协作完成。
Cond通常应用于等待某个条件的一组goroutine。
在满足条件之前，这一组goroutine是阻塞等待的

**** 使用sync.Cond时，再调用Wait前一定要上锁，Wait方法内部是先Unlock ****
*/

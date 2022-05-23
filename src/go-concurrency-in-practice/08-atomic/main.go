package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// 无符号数
	// AddUint32()
	Value()
}

func AddUint32() {
	var x uint32 = 1
	c := 2
	// 利用计算机补码的规则，把减法变成加法
	atomic.AddUint32(&x, ^uint32(c-1))
	// 减一操作可以用下面操作
	atomic.AddUint32(&x, ^uint32(0))
}

func Cas() { // 比较交换
	var x int32 = 0
	ok := atomic.CompareAndSwapInt32(&x, 0, 1)
	if ok {
		fmt.Println("比较交换成功")
	}

	// 上面代码相当于下面代码
	// if *addr == old {
	// 	*addr = new
	// 	return true
	// }
	// return false
}

type Config struct {
	NodeName string
	Addr     string
	Count    int32
}

func loadNewConfig() Config {
	return Config{
		NodeName: "深圳",
		Addr:     "192.168.3.67",
		Count:    rand.Int31(),
	}
}

func Value() { // 原子操作下的存取任意值，常用于配置更新场景
	var config atomic.Value
	config.Store(loadNewConfig())
	var cond = sync.NewCond(&sync.Mutex{})

	// 设置新的config
	go func() {
		for {
			time.Sleep(time.Duration(5+rand.Int63n(5)) * time.Second)
			config.Store(loadNewConfig())
			cond.Broadcast() // 通知等待配置更新的协程
		}
	}()

	go func() {
		for {
			cond.L.Lock()
			cond.Wait()
			c := config.Load().(Config)
			fmt.Printf("新配置：%+v\n", c)
			cond.L.Unlock()
		}
	}()

	select {}
}

/*
atomic 操作的对象是一个地址，你需要把可寻址的变量的地址作为参数传递给方法，而不是把变量的值传递给方法。
Tips:
如何给无符号数做减法？
利用补码的规则
*/

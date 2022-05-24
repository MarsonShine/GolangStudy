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

直接操作地址赋值，这个操作是原子的么？ 为什么都建议地址赋值要用aotmic
在现在的系统中，write的地址基本上都是对齐（aligned）的。比如32位系统的、CPU以及编译器，write的地址总是4的倍数；64位的地址总是8的倍数。
对齐地址的写，不会导致其他人看到只写了一半的数据，因为它通过一个指令就可以实现对地址的操作。
如果地址不是对齐的，那么处理器就需要分成两个指令去处理，如果执行了一个指令，其他人就会看到这个更新了一半的错误的数据，这种就被称为“撕裂写（torn write）”。
所以一般的认为赋值操作是一个原子操作，这个原子操作是为了保证数据的完整性。

对于多核多CPU执行的系统来说，这个就复杂了。由于cache、指令重排、可见性等问题，对原子操作的意义有了更多的目的。
在多核系统中，一个核对地址的值更改，在更新到主内存中之前，是在多级缓存中存放的。这时，看到的数据可能是不一致的，因为其它核还没有看到最新的数据，还在因没有及时更新内存（或缓存失效）使用旧的数据。
为了解决这个问题，系统就使用了一种内存屏障（memory barrier）的技术。写内存屏障会告诉处理，必须要等到它管道内未完成的操作（特别是写操作）都被刷新到内存中，再进行操作。此操作还会让相关处理器的CPU缓存失效，这样就会让其它核从主存中拉取最新数据，从而保障了数据完整性与一致性

而aotmic下的操作就提供了内存屏障操作，保障了数据的完整性与一致性。
*/

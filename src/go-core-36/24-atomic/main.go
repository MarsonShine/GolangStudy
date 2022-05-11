package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// atomic提供了哪些原子操作？
// add/sub load/store swap cas(compare and swap)

func main() {
	var inf int32
	_ = atomic.AddInt32(&inf, 1) // Q：为什么传递的参数是指针，而不是数值本身 A：因为值传递会发生复制，那么就与inf本身没什么关系了，因此也没有意义了。
	fmt.Printf("inf = %d", inf)  // inf = 1

	// 无符号的原子减法怎么实现？
	var uinf uint32
	// atomic.AddUint32(&uinf, -5) // 报错，因为-5不是uint32类型
	// 类型转换
	delta := int32(-5) // 报错，同理 int32 的范围与uint32不一致
	// 通过delta类型转换
	atomic.AddUint32(&uinf, uint32(delta))
	// 或者通过补码
	atomic.AddUint32(&uinf, ^uint32(-(-5)-1)) // 与-5的补码相同

	// cas模拟自旋锁
	spinLock()
}

func spinLock() {
	var num int32 = 0
	for {
		if atomic.CompareAndSwapInt32(&num, 10, 0) {
			fmt.Printf("The second number has gone to zero.")
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
}

// 模拟简单的自旋锁
func forAndCAS() {
	sign := make(chan struct{}, 2)
	num := int32(0)
	fmt.Printf("The number: %d\n", num)

	go func() {
		defer func() {
			sign <- struct{}{}
		}()

		for {
			time.Sleep(time.Millisecond * 500)
			newNum := atomic.AddInt32(&num, 2) // 原子操作
			fmt.Printf("The number: %d\n", newNum)
			if newNum == 10 {
				break
			}
		}
	}()

	go func() {
		defer func() {
			sign <- struct{}{}
		}()

		for {
			if atomic.CompareAndSwapInt32(&num, 10, 0) {
				fmt.Println("The number has gone to zero.")
				break
			}
		}
		time.Sleep(time.Millisecond * 500)
	}()

	<-sign
	<-sign
}

// 通过cas模拟互斥锁
func forAndCAS2() {
	sign := make(chan struct{}, 2)
	num := int32(0)
	fmt.Printf("The number: %d\n", num)
	max := int32(20)

	go func(id int, max int32) {
		defer func() {
			sign <- struct{}{}
		}()
		for i := 0; ; i++ {
			currNum := atomic.LoadInt32(&num)
			if currNum >= max {
				break
			}
			newNum := currNum + 2
			time.Sleep(time.Millisecond * 200)
			if atomic.CompareAndSwapInt32(&num, currNum, newNum) {
				fmt.Printf("The number: %d [%d-%d]\n", newNum, id, i)
			} else {
				fmt.Printf("The CAS operation failed. [%d-%d]\n", id, i)
			}
		}
	}(1, max)

	go func(id int, max int32) { // 定时增加num的值
		defer func() {
			sign <- struct{}{}
		}()
		for j := 0; ; j++ {
			currNum := atomic.LoadInt32(&num)
			if currNum >= max {
				break
			}
		}
	}(2, max)

	<-sign
	<-sign
}

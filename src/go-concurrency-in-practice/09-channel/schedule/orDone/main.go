package main

import (
	"fmt"
	"reflect"
	"time"
)

// channel 编排模式之orDone：只要任意一个任务执行完，就返回；
// 信号通知模式
func main() {
	start := time.Now()
	<-orWithoutRecurison(
		sig(10 * time.Second),
		// sig(20*time.Second),
		// sig(30*time.Second),
		// sig(40*time.Second),
		// sig(50*time.Second),
		// sig(01*time.Minute),
	)
	// 执行到这里，说明orDone已经被close掉了
	fmt.Printf("done after %v", time.Since(start))
}

// 先定义一个chan struct 的done变量，起到当任务完成时触发通知作用
// 具体做法是，当任务完成时，close chan，这时 receiver 就会收到这个通知
func or(channels ...<-chan interface{}) <-chan interface{} {
	// 特殊情况，只有零个或者1个chan
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2: //2也是特殊情况
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: //超过2个，2分法递归处理
			m := len(channels) / 2
			select {
			case <-or(channels[:m]...):
			case <-or(channels[m:]...):
			}
		}
	}()
	return orDone
}

func orWithoutRecurison(channels ...<-chan interface{}) <-chan interface{} {
	// 特殊情况，只有零个或者1个chan
	switch len(channels) {
	case 0:
		fmt.Println(0)
		return nil
	case 1:
		fmt.Println(1)
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer func() {
			fmt.Println("close orDone...")
			close(orDone)
		}()
		// 利用反射
		var cases []reflect.SelectCase
		for _, c := range channels {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		// 随机选择一个可用的case
		c, _, ok := reflect.Select(cases)
		// ch := r.Interface().(chan interface{})
		if ok {
			fmt.Println(c)
		} else {
			fmt.Println("not ok")
		}
	}()
	return orDone
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

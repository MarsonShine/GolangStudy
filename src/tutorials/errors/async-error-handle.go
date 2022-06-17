package errors

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 在异步情况下的错误异常处理
// 采用类Javascript的Promise错误处理模式
type Promise struct {
	wg  sync.WaitGroup
	err error
	msg string
}

func NewPromise(f func() (string, error)) *Promise {
	p := &Promise{}
	p.wg.Add(1)
	go func() {
		p.msg, p.err = f()
		p.wg.Done()
	}()
	return p
}

func (p *Promise) Then(r func(string), e func(error)) *Promise {
	go func() {
		p.wg.Wait()
		if p.err != nil {
			e(p.err)
			return
		}
		r(p.msg)
	}()
	return p
}

func example() (string, error) {
	for i := 0; i < 3; i++ {
		fmt.Println(i)
		<-time.Tick(time.Second * 1)

	}
	rand.Seed(time.Now().UTC().UnixNano())
	r := rand.Intn(100) % 2
	fmt.Println(r)
	if r != 0 {
		return "hello, world", nil
	} else {
		return "", fmt.Errorf("error")
	}
}

func errorHandle() {
	doneChan := make(chan int)
	var p = NewPromise(example)
	p.Then(func(s string) { fmt.Println(s); doneChan <- 1 },
		func(err error) { fmt.Println(err); doneChan <- 1 })
	<-doneChan
}

// 为了达到通用化，我们需要泛化函数参数，go 1.18提供了泛型，所以我们可以利用泛型实现泛化通用的promise异常处理模式

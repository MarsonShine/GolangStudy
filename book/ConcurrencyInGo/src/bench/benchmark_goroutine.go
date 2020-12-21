package main

import (
	sync "sync"
	"testing"
)

func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var token struct{}
	sender := func() {
		defer wg.Done()
		<-begin // 1 这里会被阻塞，直到接受到数据。我们不希望设置和启动goroutine影响上下文切换的度量。
		for i := 0; i < b.N; i++ {
			c <- token // 2 向接收者发送数据。struct{}{}是空结构体且不占用内存；这样我们就可以做到只测量发送信息所需要的时间。
		}
	}

	receiver := func() {
		defer wg.Done()
		<-begin // 1
		for i := 0; i < b.N; i++ {
			<-c //3 接收传递过来的数据，但不做任何事。
		}
	}

	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer() //4 启动计时器。
	close(begin)   //5 通知发送和接收的goroutine启动。
	wg.Wait()
}

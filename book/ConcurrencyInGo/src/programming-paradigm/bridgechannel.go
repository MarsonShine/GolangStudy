package main

import "fmt"

// 通道桥接，定义一个函数将多个channel拆解为一个简单的channel
var bridge = func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{}) // 这个通道会返回所有传入 bridge 的通道
	go func() {
		defer close(valStream)
		for { // 该循环负责从 chanStream 中提取并将其值提供给嵌套循环使用
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range orDone(done, stream) { // 该循环负责读取已经给出的通道的值，并将这个值发送给 valStream
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func brideChannelExample() {
	// 创建10个通道，每个通道写入一个元素
	// 并将这些通道传入给bridge
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}

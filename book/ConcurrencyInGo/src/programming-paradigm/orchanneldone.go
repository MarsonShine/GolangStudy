package main

// 要与来自系统的不同部分的channel交互，如果其中一个goroutine取消操作，那么我们没有更进一步的信息进行下一步操作
// 也就是说我们不知道这些 goroutine 处于什么状态，所以处于对 goroutine 泄露的考虑，需要加一个 done 通道来封装这些 goroutine 操作

var orDone = func(done, c <-chan interface{}) <-chan interface{} {
	varStream := make(chan interface{})
	go func() {
		defer close(varStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case varStream <- v:
				case <-done:
				}

			}
		}
	}()
	return varStream
}

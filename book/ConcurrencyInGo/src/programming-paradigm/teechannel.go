package main

import "fmt"

// 想要将多个channel分割，发送到两个独立的区域
// 比如：在一个通道上接受一系列指令，将它们发送给执行者，同时记录日志

// tee 名字的来源：unix 的 tee 命令
var tee = func(done <-chan interface{},
	in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range orDone(done, in) {
			var out1, out2 = out1, out2 // 拷贝至本地变量
			for i := 0; i < 2; i++ {    // for select 模式可以让写入 out1,out2 不会阻塞
				select {
				case <-done:
				case out1 <- val: // 一旦我们写入了通道，我们将其副本设置为零，这样继续写入将阻塞，而另一个通道可以继续执行。
					fmt.Println("out1 <- val")
					out1 = nil
				case out2 <- val: // 一旦我们写入了通道，我们将其副本设置为零，这样继续写入将阻塞，而另一个通道可以继续执行。
					fmt.Println("out2 <- val")
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

func teeExample() {
	done := make(chan interface{})
	defer close(done)
	// 利用这种模式，很容易使用通道作为系统数据的连接点。
	out1, out2 := tee(done, take(done, repeat(done, 1, 2, 3), 9))

	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val1, <-out2)
	}
}

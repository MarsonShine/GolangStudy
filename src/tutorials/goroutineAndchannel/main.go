package main

// 每一个并发的执行单元叫作一个goroutine
// 当一个程序启动时，其主函数即在一个单独的goroutine中运行，我们叫它main goroutine。新的goroutine会用go语句来创建
// f()
// go f() // go 语句后面就开始一个 goroutine 执行 f() 方法
// channel 还能实现 pipeline，如一个channel的输出可以作为下一个channel的输入
// ch = make(chan string, 3) 这是表示一个带缓存的channel，内部维护一个有三个元素的队列。
// 这样我们就能无阻塞的往队列里面发送三个消息了，超过3个就要阻塞了
// 可以利用 channel 做到多路复用（multiplex） select { case}，一个空的多路复用即 select{} 这样会永远等待下去
// select {
// case <-ch1:
// 	//...
// case x:= <-ch2:
// 	// ... use x ...
// case ch3 <- y:
// 	// ...
// default:
// 	//
// }

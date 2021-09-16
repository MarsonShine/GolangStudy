package main

import (
	"fmt"
	"time"
)

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
// 不要使用共享数据来通信；使用通信来共享数据

func main() {
	// select + default 实现异步阻塞
	noHappen()
	happen()
	workPool()
}

func noHappen() {
	messages := make(chan string)
	signals := make(chan bool)

	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}

	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("send message", msg)
	default:
		fmt.Println("no message send")
	}

	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

func happen() {
	messages := make(chan string)
	signals := make(chan bool)
	go func(msg string) {
		messages <- msg
	}("marsonshine")
	time.Sleep(time.Second * 1)
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}

	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("send message", msg)
	default:
		fmt.Println("no message send")
	}

	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

func closeChanAfterJobFinished() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done <- true
				return
			}
		}
	}()

	for i := 0; i < 3; i++ {
		jobs <- i
		fmt.Println("send job", i)
	}
	close(jobs)
	fmt.Println("send all jobs")

	<-done
}

func closeWhenChanHasData() {
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)
	// 关闭之后读取不会报错，而是将剩下的消费完
	for elem := range queue {
		fmt.Println(elem)
	}
}

func workPool() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

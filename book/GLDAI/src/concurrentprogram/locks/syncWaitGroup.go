// sync.waitGroup 等待一组 Goroutine 的返回
// waitGroup 内部有一个计数器，初始化从 0 开始，有三个方法：Add(),Done(),Wait() 可以控制这个计数器
// Add(n) 将计数器直接设置为 n
// Done() 将计数器减 1
// Wait() 判断计数器是否等于 0，不等于 0 阻塞当前 Goroutine

// x/sync/errgroup.Group 在一组 goroutine 中提供了同步、错误传播以及上下文取消的功能
// x/sync/errgroup.Go 一个 goroutine 进入的执行体
// x/sync/errgroup.Wait 等待所有 goroutine，返回结果如果为空，表示所有 goroutine 成功；如果返回错误，表示至少有一个 goroutine 错误
// 只有第一个错误出现时就返回，剩下的错误会直接被抛弃
package locks

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

func Start() {
	// 有 100 个 Goroutine，要等待这 100 个线程全部完成
	// 方法1：用 Channel 来实现同步
	method1()
}

func ErrGroup() {
	var g errgroup.Group
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}
	for i := range urls {
		url := urls[i]
		g.Go(func() error {
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			return err
		})
	}
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully fetched all URLs.")
	}
}

func method2() {
	wg := sync.WaitGroup{}
	wg.Add(100) // 初始化计数器 100
	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done() // 完成一个减 1
		}(i)
	}
	wg.Wait() // 等待所有 goroutine 执行完毕
}

func method1() {
	c := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			fmt.Println(i)
			c <- true
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-c // 一直阻塞，知道 c 接收到发送的值
	}
}

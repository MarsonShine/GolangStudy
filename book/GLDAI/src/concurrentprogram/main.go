// context.Context 是用来设置截止日期、同步信号，传递请求相关值的结构体
// 每个请求就是一个 goroutine
// context.Context 的作用就是在不同 Goroutine 之间同步请求特定数据、取消信号以及处理请求的截止日期。
// 每一个 context.Context 都会从最顶层的 Goroutine 一层一层传递到最下层
// 在上层出现异常时，会发信号以及携带错误信息传递下一层，所以下层接收到信号应该要及时停止无用的工作以减少额外资源的浪费
// context.Background() 返回一个默认的上下文，如果程序没有其它上下文信息的化，一般用这个来传递
// context.TODO() 这个是一个空的上下文，什么也不做，这个是因为前期可能不确定用什么上下文，然后利用多态来替换
// context.WithTimeout() 会生成一个子上下文并返回用于取消该上下文的函数。所以父协程提前结束或是这个字协程主动取消都会引起该上下文的函数取消
// context.WithValue 函数能从父上下文中创建一个子上下文

package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// context.Context 用来对信号进行同步
	// 设置一个过期时间为1s的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go handle(ctx, 100*time.Millisecond) // 将超时1秒的上下文传递给函数，并用指定的500ms处理
	// 因为设置的超时时间是大于处理运行的时间
	// 所以程序会正常运行结束，所以handle程序会输出 process request with 500ms
	// 但是 main 函数的 select 会等待 context.Context 的超时
	// 所以会输出 main context deadline exceeded
	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
		// 加了 defaul，有多个分支，主协程select不会等待而直接执行 default 块
		// 这里主协程就会结束，那么子协程没有运行完就会取消，所以 handle 就会输出 handle context canceled
		// default:
		// 	fmt.Println("main done", ctx.Err())
	}
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}

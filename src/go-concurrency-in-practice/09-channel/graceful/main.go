package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var closing = make(chan struct{}) // 表示非正常关闭，如因超时
	var closed = make(chan struct{})  // 正常关闭

	go func() {
		for {
			select {
			case <-closing:
				return
			case <-closed:
				//...
			default:
				// ...业务计算
				time.Sleep(100 * time.Microsecond)
			}
		}
	}()
	terminalChan := make(chan os.Signal)
	signal.Notify(terminalChan, syscall.SIGINT, syscall.SIGTERM)
	<-terminalChan

	// 执行退出之前的清理动作
	close(closing)
	go cleanup(closed)
	select {
	case <-closed:
	case <-time.After(time.Second):
		fmt.Println("清理超时，强制退出")
	}
	fmt.Println("优雅退出")

}

func cleanup(ch chan struct{}) {
	// clean
	time.Sleep(time.Second) // 模拟清理耗时
	close(ch)
}

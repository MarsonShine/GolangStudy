package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 检测操作系统的信号，进行额外的操作
	// 此方法可以用作给程序进行优雅的退出，如退出前通知应用程序

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}

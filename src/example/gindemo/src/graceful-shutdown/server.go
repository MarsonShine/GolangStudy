package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(200, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	// 开启一个协程初始化 server
	// 这样不会阻塞下面优雅的关机处理
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 使用 5 秒超时处理等待中断信号来优雅的处理关机
	quit := make(chan os.Signal)
	// kill 无参默认发送系统调用：syscall.SIGTERM
	// kill -2 syscall.SIGINT
	// kill -9 syscall.SIGKILL，但是无法捕捉异常，所以没必要添加它
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, calcel := context.WithTimeout(context.Background(), 5*time.Second)
	defer calcel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}

package main

import (
	myhttp "basicMiddlewareHttpServer/application/http"
	"basicMiddlewareHttpServer/persistents/middleware"
	"log"
	"net/http"
	"os"
)

// 不依靠第三路由包，就需要自己写解析路由算法
func main() {
	mux := http.NewServeMux()
	server := myhttp.NewTaskHttpServer()
	mux.HandleFunc("/task/", server.TaskHandler)

	// 注册中间件
	handler := middleware.Logging(mux) // 全局中间件
	// 如果只针对某个路由添加中间件，则可以
	// handler := middleware.Logging(http.HandlerFunc(server.TaskHandler))
	// mux.Handle("/task/", handler)

	// 注册多个
	// handler := middleware.Logging(mux)
	// handler = middleware.PanicRecovery(handler)

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), handler))

	// 中间件执行顺序
	// request --> [Logging] --> [Mux] --> [Handler]

	// 注册多个的执行顺序
	// request --> [Panic Recovery] --> [Logging] --> [Mux] --> [tagHandler]
}

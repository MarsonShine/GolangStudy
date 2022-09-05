package main

import (
	myhttp "basicHttpServer/application/http"
	"log"
	"net/http"
	"os"
)

// 不依靠第三路由包，就需要自己写解析路由算法
func main() {
	mux := http.NewServeMux()
	server := myhttp.NewTaskHttpServer()
	mux.HandleFunc("/task/", server.TaskHandler)
	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))
}

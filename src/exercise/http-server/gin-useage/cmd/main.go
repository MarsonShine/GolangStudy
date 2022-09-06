package main

import (
	myhttp "ginHttpServer/application/http"

	"github.com/gin-gonic/gin"
)

// 不依靠第三路由包，就需要自己写解析路由算法
func main() {
	router := gin.Default()

	server := myhttp.NewTaskHttpServer()

	router.POST("/task/", server.CreateTaskHandler)
	router.GET("/task/", server.GetAllTasksHandler)
	router.GET("/task/:id", server.GetTaskHandler)
	router.Run("localhost:5000")
}

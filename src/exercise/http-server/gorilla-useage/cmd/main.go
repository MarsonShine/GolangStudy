package main

import (
	myhttp "gorillaHttpServer/application/http"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// 不依靠第三路由包，就需要自己写解析路由算法
func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := myhttp.NewTaskHttpServer()

	router.HandleFunc("/task/", server.CreateTaskHandler).Methods("POST")
	router.HandleFunc("/task/", server.GetAllTasksHandler).Methods("GET")
	router.HandleFunc("/task/{id:[0-9]+}/", server.GetTaskHandler).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:5000", router))
}

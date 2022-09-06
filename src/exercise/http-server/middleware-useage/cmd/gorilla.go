package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func middlewareWithGorilla() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	// server := myhttp.NewTaskHttpServer()
	// router.HandleFunc("/task/", server.createTaskHandler).Methods("POST")
	// router.HandleFunc("/task/", server.getAllTasksHandler).Methods("GET")
	// router.HandleFunc("/task/", server.deleteAllTasksHandler).Methods("DELETE")
	// router.HandleFunc("/task/{id:[0-9]+}/", server.getTaskHandler).Methods("GET")
	// 注册中间件
	router.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	})
	router.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	log.Fatal(http.ListenAndServe("localhost:5000", router))
}

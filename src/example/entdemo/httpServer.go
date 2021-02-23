package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

const stopTimeout = time.Second * 10

func startHTTPServer() *http.Server {
	router := mux.NewRouter().StrictSlash(true)
	srv := &http.Server{Addr: ":8081", Handler: router}
	sigs := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(sigs, os.Interrupt)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world\n")
	})

	router.HandleFunc("/user/{id:[0-9]+}", getUserHandler)

	router.HandleFunc("/user/create", createUserHandler)
	router.HandleFunc("/user/delete/{id:[0-9]+}", deleteUserHandler)
	router.HandleFunc("/user/delete/{name:[a-zA-Z]+}", deleteUserByNameHandler)
	go func() {
		<-sigs
		ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Server shutdown with error: ", err)
		}
		close(done)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("Http server start failed", err)
	}
	<-done

	return srv
}

func writeBackStream(w http.ResponseWriter, data interface{}) {
	jsonResponse, _ := json.Marshal(data)
	w.Header().Set("content-type", "application/json")
	w.Write(jsonResponse)
}

type DataResponse struct {
	Success bool
	Message string
	Data    interface{}
}

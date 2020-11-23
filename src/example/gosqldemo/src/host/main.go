package main

import (
	"context"
	"encoding/json"
	"gosqldemo/src/app/userservice"
	"gosqldemo/src/domain"
	"gosqldemo/src/dto"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const stopTimeout = time.Second * 10

func main() {

	startHTTPServer()
}

func startHTTPServer() *http.Server {
	router := mux.NewRouter().StrictSlash(true)
	srv := &http.Server{Addr: ":8080", Handler: router}
	sigs := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(sigs, os.Interrupt)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world\n")
	})

	router.HandleFunc("/user", getUserListHandler)
	router.HandleFunc("/user/{id:[0-9]+}", getUserHandler)

	router.HandleFunc("/user/create", createUserHandler)

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

func getUserListHandler(w http.ResponseWriter, r *http.Request) {
	data := &dto.DataResponse{Success: true, Message: "success", Data: []string{}}
	// var jsonResponse = []byte(data)
	var users = userservice.NewUserService().GetAllUsers()
	data.Data = users
	jsonResponse, _ := json.Marshal(data)
	w.Header().Set("content-type", "text/json")
	w.Write(jsonResponse)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	data := &dto.DataResponse{Success: false}
	if err != nil {
		data.Message = "false"
	} else {
		var user domain.User
		user = userservice.NewUserService().GetUserById(userID)
		data.Success = true
		data.Data = user
	}
	jsonResponse, _ := json.Marshal(data)
	w.Header().Set("content-type", "text/json")
	w.Write(jsonResponse)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u struct {
		Name     string
		Email    *string
		Age      uint8
		Birthday *domain.JSONTime
	}
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}
	t := time.Time(*u.Birthday)

	err = userservice.NewUserService().CreateUser(domain.User{Name: u.Name, Email: u.Email, Age: u.Age, Birthday: &t})
	var jsonResponse []byte
	if err != nil {
		jsonResponse = []byte(`{"sucess":false, "message": "` + err.Error() + `"}`)
	} else {
		jsonResponse = []byte(`{"sucess":true, "message": "success!"}`)
	}

	w.Header().Set("content-type", "text/json")
	w.Write(jsonResponse)
}

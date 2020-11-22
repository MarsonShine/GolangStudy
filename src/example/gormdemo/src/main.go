package main

import (
	"context"
	"encoding/json"
	appservice "gormdemo/src/app"
	"gormdemo/src/models"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const stopTimeout = time.Second * 10
const dsn = "root:123456@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	// db, err := gorm.Open(sqlite.Open("./src/test.db"), &gorm.Config{})
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败！")
	}

	// migration
	db.AutoMigrate(&models.Product{}, &models.User{})

	// // Create
	// db.Create(&models.Product{Code: "D42", Price: 100})
	// // Read
	// var product models.Product
	// db.First(&product, 1)
	// db.First(&product, "code = ?", "D42") // 查找 code 字段为 D42

	// // Upload
	// db.Model(&product).Update("Price", 200)
	// // 更新多个字段
	// db.Model(&product).Updates(models.Product{Price: 200, Code: "F42"})
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	// db.Delete(&product, 1)
	startHTTPServer()
}

// 公共返回体
type DataResponse struct {
	Success bool
	Message string
	Data    interface{}
}

func getUserListHandler(w http.ResponseWriter, r *http.Request) {
	users := appservice.GetUserAll()
	data := &DataResponse{Success: true, Message: "success", Data: []string{}}
	if users != nil {
		data.Data = users
	}
	// var jsonResponse = []byte(data)
	jsonResponse, _ := json.Marshal(data)
	w.Header().Set("content-type", "text/json")
	w.Write(jsonResponse)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	data := &DataResponse{Success: false}
	if err != nil {
		data.Message = "false"
	} else {
		user := appservice.GetUser(id)
		data.Success = true
		data.Data = user
	}
	jsonResponse, _ := json.Marshal(data)
	w.Header().Set("content-type", "text/json")
	w.Write(jsonResponse)
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
	router.HandleFunc("/user/{id}", getUserHandler)

	router.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var u struct {
			Name     string
			Email    *string
			Age      uint8
			Birthday *models.JSONTime
		}
		err := decoder.Decode(&u)
		if err != nil {
			panic(err)
		}
		sysTime := time.Time(*u.Birthday)
		appservice.CreateUserService(u.Name, u.Email, u.Age, &sysTime)
		var jsonResponse = []byte(`{"sucess":true, "message": "success!"}`)
		w.Header().Set("content-type", "text/json")
		w.Write(jsonResponse)
	})

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

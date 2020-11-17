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
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const stopTimeout = time.Second * 10

func main() {
	db, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
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

func startHTTPServer() *http.Server {
	srv := &http.Server{Addr: ":8080"}
	sigs := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(sigs, os.Interrupt)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world\n")
	})

	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
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

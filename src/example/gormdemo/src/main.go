package main

import (
	"context"
	"database/sql"
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
	"gorm.io/gorm/logger"
)

const stopTimeout = time.Second * 10

const dsn = "root:123456@tcp(192.168.3.10:3306)/go_testdb?charset=utf8mb4&parseTime=True&loc=Local"

// const dsn = "root:123456@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

var gormDB *gorm.DB

func main() {
	// db, err := gorm.Open(sqlite.Open("./src/test.db"), &gorm.Config{})
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	panic("连接数据库失败！")
	// }
	// sqlDB, _ := db.DB()
	// defer sqlDB.Close()
	// // migration
	// db.AutoMigrate(&models.Product{}, &models.User{})

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
	initializeDataBase()
	startHTTPServer()
}

// 公共返回体
type DataResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func initializeDataBase() {
	gormDB = openDbConnection()
}

func NewDataResponse() DataResponse {
	dr := DataResponse{}
	dr.Success = dr.Message == ""
	return dr
}

func (dr DataResponse) SetMessage(msg string) {
	dr.Message = msg
	dr.Success = false
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

func getUserListSimpleHandler(w http.ResponseWriter, r *http.Request) {
	data := &DataResponse{Success: true, Message: "success", Data: []string{}}
	// var mydb = openDbConnection()
	var users []models.User
	_, _ = gormDB.Select([]string{}).Find(&users).DB()

	data.Data = users
	writeBackStream(w, data)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	data := &DataResponse{Success: false}
	if err != nil {
		data.Message = "false"
	} else {
		user := appservice.NewUserService(gormDB).Get(uint(id))
		data.Success = true
		data.Data = user
	}
	writeBackStream(w, data)
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

	router.HandleFunc("/user", getUserListSimpleHandler)
	router.HandleFunc("/user/{id:[0-9]+}", getUserHandler)

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

		appservice.CreateUserService(gormDB, u.Name, u.Email, u.Age, &sysTime)
		var jsonResponse = []byte(`{"sucess":true, "message": "success!"}`)
		w.Header().Set("content-type", "text/json")
		w.Write(jsonResponse)
	})
	router.HandleFunc("/user/delete/{id:[0-9]+}", deleteUserHandler)

	productRouterInitialize(router)

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

func productRouterInitialize(router *mux.Router) {
	router.HandleFunc("/product/create", createProductHandler)
	router.HandleFunc("/product/{id:[0-9]+}", productDetailHandler)
	router.HandleFunc("/product/update", productUpdateHandler)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p models.Product
	data := NewDataResponse()
	err := decoder.Decode(&p)
	if err != nil {
		data.SetMessage(err.Error())
	} else {
		err := appservice.NewProductService(gormDB).CreateProduct(&p)
		if err != nil {
			data.SetMessage(err.Error())
		}
	}
	jsonResponse, _ := json.Marshal(data)
	w.Header().Set("content-type", "text/json")
	w.Write(jsonResponse)
}

func productDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := NewDataResponse()
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		data.Message = "false"
	} else {
		productDetail := appservice.NewProductService(gormDB).GetProductDetail(uint(productID))
		data.Data = productDetail
		data.Success = true
	}
	writeBackStream(w, data)
}

func productUpdateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p appservice.ProductUpdated
	data := NewDataResponse()
	err := decoder.Decode(&p)
	if err != nil {
		data.SetMessage(err.Error())
	} else {
		r, err := appservice.NewProductService(gormDB).UpdateProductAndUser(&p)
		if err != nil {
			data.SetMessage(err.Error())
		} else {
			data.Success = r
		}
	}
	writeBackStream(w, data)
}

func writeBackStream(w http.ResponseWriter, data interface{}) {
	jsonResponse, _ := json.Marshal(data)
	w.Header().Set("content-type", "text/json")
	w.Write(jsonResponse)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	data := NewDataResponse()
	if err != nil {
		data.Message = "false"
	} else {
		data.Success = appservice.NewUserService(gormDB).Delete(uint(userID))
	}
	writeBackStream(w, data)
}

func openDbConnection() *gorm.DB {
	// db, err := gorm.Open(sqlite.Open("./src/test.db"), &gorm.Config{})
	sqlDB, err := sql.Open("mysql", dsn)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(time.Millisecond * 200)
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		sqlDB.Close()
		panic("连接数据库失败！")
	}
	return db
}

package main

import (
	"context"
	"entdemo/ent"
	"entdemo/ent/car"
	"entdemo/ent/group"
	"entdemo/ent/user"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	ConnectionString string = "root:123456@tcp(192.168.3.10:3306)/go_testdb?charset=utf8mb4&parseTime=True&loc=Local"
)

var client *ent.Client

func main() {
	drv, err := sql.Open("mysql", ConnectionString)
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}

	sqlDB := drv.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(152)
	sqlDB.SetConnMaxLifetime(time.Millisecond * 100)

	// c, err := ent.Open(dialect.MySQL, ConnectionString)
	client = ent.NewClient(ent.Driver(drv), ent.Debug(), ent.Log(sqlLogging))
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}
	defer client.Close()
	// 运行自动迁移工具
	if err := client.Schema.
		Create(context.Background()); err != nil {
		log.Fatalf("添加实体表失败：%v", err)
	}
	startHTTPServer()
}

func sqlLogging(opts ...interface{}) {
	log.Print(opts...)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := QueryUser(r.Context(), client)
	resp := &DataResponse{Success: false}
	if err != nil {
		resp.Message = err.Error()
		writeBackStream(w, resp)
	} else {
		resp.Success = true
		resp.Data = user
		writeBackStream(w, resp)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := CreateUser(context.Background(), client)
	resp := &DataResponse{Success: false}
	if err != nil {
		resp.Message = err.Error()
		resp.Success = false
		writeBackStream(w, resp)
	} else {
		resp.Data = user
		resp.Success = true
		writeBackStream(w, resp)
	}
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vas := mux.Vars(r)
	resp := &DataResponse{Success: false}
	id, err := strconv.Atoi(vas["id"])
	if err != nil {
		resp.Message = "name 参数错误"
		writeBackStream(w, resp)
	} else {
		b, err := DeleteUser(r.Context(), client, id)
		if err != nil {
			resp.Message = err.Error()
		} else {
			resp.Success = b
		}
		writeBackStream(w, resp)
	}
}

func deleteUserByNameHandler(w http.ResponseWriter, r *http.Request) {
	vas := mux.Vars(r)
	resp := &DataResponse{Success: false}
	b, err := DeleteUserByName(r.Context(), client, vas["name"])
	if err != nil {
		resp.Message = err.Error()
	} else {
		resp.Success = b
	}
	writeBackStream(w, resp)
}

func QueryGithub(ctx context.Context, client *ent.Client) error {
	cars, err := client.Group.
		Query().
		Where(group.Name("GitHub")). // (Group(Name=GitHub),)
		QueryUsers().                // (User(Name=Ariel, Age=30),)
		QueryCars().                 // (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Mazda, RegisteredAt=<Time>),)
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %v", err)
	}
	log.Println("cars returned:", cars)
	// Output: (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Mazda, RegisteredAt=<Time>),)
	return nil
}

func QueryArielCars(ctx context.Context, client *ent.Client) error {
	// Get "Ariel" from previous steps.
	a8m := client.User.
		Query().
		Where(
			user.HasCars(),
			user.Name("Ariel"),
		).
		OnlyX(ctx)
	cars, err := a8m. // Get the groups, that a8m is connected to:
				QueryGroups(). // (Group(Name=GitHub), Group(Name=GitLab),)
				QueryUsers().  // (User(Name=Ariel, Age=30), User(Name=Neta, Age=28),)
				QueryCars().   //
				Where(         //
			car.Not( //  Get Neta and Ariel cars, but filter out
				car.ModelEQ("Mazda"), //  those who named "Mazda"
			), //
		). //
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %v", err)
	}
	log.Println("cars returned:", cars)
	// Output: (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Ford, RegisteredAt=<Time>),)
	return nil
}

func QueryGroupWithUsers(ctx context.Context, client *ent.Client) error {
	groups, err := client.Group.
		Query().
		Where(group.HasUsers()).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting groups: %v", err)
	}
	log.Println("groups returned:", groups)
	// Output: (Group(Name=GitHub), Group(Name=GitLab),)
	return nil
}

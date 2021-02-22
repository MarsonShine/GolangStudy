package main

import (
	"context"
	"entdemo/ent"
	"entdemo/ent/user"
	"fmt"
	"log"
	"net/http"
	"time"

	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
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
	sqlDB.SetConnMaxLifetime(time.Millisecond * 200)

	// c, err := ent.Open(dialect.MySQL, ConnectionString)
	client = ent.NewClient(ent.Driver(drv))
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

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := QueryUser(r.Context(), client)
	if err != nil {
		writeBackStream(w, struct {
			message string
			success bool
		}{
			err.Error(),
			false,
		})
	} else {
		writeBackStream(w, struct {
			message string
			success bool
			data    interface{}
		}{
			"",
			true,
			user,
		})
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

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	sex := false
	u, err := client.User.
		Create().
		SetName("marsonshine").
		SetAge(27).
		SetAddress("深圳市南山区桃园街道创新大厦").
		SetNillableSex(&sex).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("添加用户失败：%v", err)
	}
	log.Printf("创建用户成功：%v", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.NameEQ("marsonshine")).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %v", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

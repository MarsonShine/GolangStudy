package repository

import (
	"sync"

	"github.com/jmoiron/sqlx"
	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

const dsn = ""

var db *sqlx.DB
var once sync.Once
var mu sync.Mutex

func OpenDbConnection() *sqlx.DB {
	// db, err := sqlx.Open("sqlite3", "./src/test.db")
	db, err := sqlx.Connect("mysql", "root:123456@tcp(192.168.3.10:3306)/go_testdb?charset=utf8mb4&parseTime=True&loc=Local")
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(130)
	// db.SetConnMaxLifetime(time.Millisecond * 200)
	if err != nil {
		panic("连接数据库失败！")
	}
	return db
}

func singletonInstance() *sqlx.DB {
	once.Do(func() {
		db = OpenDbConnection()
	})
	return db
}

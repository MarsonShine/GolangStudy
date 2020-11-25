package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

const dsn = ""

func OpenDbConnection() *sqlx.DB {
	// db, err := sqlx.Open("sqlite3", "./src/test.db")
	db, err := sqlx.Connect("mysql", "root:123456@tcp(192.168.3.10:3306)/go_testdb?charset=utf8mb4&parseTime=True&loc=Local")
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(130)
	db.SetConnMaxLifetime(time.Millisecond * 200)
	if err != nil {
		panic("连接数据库失败！")
	}
	return db
}

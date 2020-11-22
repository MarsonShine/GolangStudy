package repository

import (
	"github.com/jmoiron/sqlx"
	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

const dsn = ""

func OpenDbConnection() *sqlx.DB {
	// db, err := sqlx.Open("sqlite3", "./src/test.db")
	db, err := sqlx.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("连接数据库失败！")
	}
	err = db.Ping()
	return db
}

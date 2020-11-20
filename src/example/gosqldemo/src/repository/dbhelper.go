package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDbConnection() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "./src/test.db")
	if err != nil {
		panic("连接数据库失败！")
	}
	err = db.Ping()
	return db
}

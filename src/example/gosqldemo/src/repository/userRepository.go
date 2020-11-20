package repository

import (
	"fmt"
	"gosqldemo/src/domain"
	"sync"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
}

var db *sqlx.DB

var mu sync.Mutex

func (repository *UserRepository) New() {
	mu.Lock()
	defer mu.Unlock()
	db = OpenDbConnection()
}

// 获取所有用户
func (repository UserRepository) GetUserAll() *[]domain.User {
	rows, err := db.Query("SELECT * FROM users")
	var users *[]domain.User
	if err == nil {
		for rows.Next() {
			var user domain.User
			err = rows.Scan(&user)
			*users = append(*users, user)
		}
	}
	return users
}

// 获取指定用户
func (repository UserRepository) GetUser(id int) domain.User {
	var user domain.User
	err := db.Get(&user, "select * from users where id = ?", id)
	if err != nil {
		panic("用户不存在")
	}
	return user
}

// 创建用户
func (repository UserRepository) CreateUser(u domain.User) {
	rows, _ := db.Query("select * from users where email = ?", u.Email)
	if !rows.Next() {
		panic(fmt.Sprintf("用户已存在: %s\n", *u.Email))
	}
	_, err := db.Exec("insert into users (name,email,age,birthday,createdAt) values (?, ?, ?, ?, ?)")
	if err != nil {
		panic(fmt.Sprintf("用户添加失败: %s\n", err.Error()))
	}
}

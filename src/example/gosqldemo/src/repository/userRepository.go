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
func (repository UserRepository) GetUserAll() []domain.User {
	users := []domain.User{}
	err := db.Select(&users, "SELECT id,name,email,age,birthday FROM users")
	db.Close()
	if err != nil {
		db.Close()
		panic(err)
	}
	return users
}

// 获取指定用户
func (repository UserRepository) GetUser(id int) domain.User {
	var user domain.User
	err := db.Get(&user, "select id,name,email,age,birthday,member_number,actived_at,created_at,updated_at,deleted_at from users where id = ?", id)
	defer db.Close()
	if err != nil {
		db.Close()
		panic("用户不存在")
	}
	return user
}

// 创建用户
func (repository UserRepository) CreateUser(u domain.User) error {
	_, err := db.Query("select * from users where email = ?", u.Email)
	db.Close()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	// if rows.Next() {
	// 	return fmt.Errorf("用户已存在: %s\n", *u.Email)
	// }
	_, err = db.Exec("insert into users (name,email,age,birthday,created_at) values (?, ?, ?, ?, NOW())", u.Name, u.Email, u.Age, u.Birthday)
	if err != nil {
		return fmt.Errorf("用户添加失败: %s\n", err.Error())
	}
	return nil
}

func (repository UserRepository) DeleteUser(id int) (bool, error) {
	defer func() {
		if r := recover(); r != nil {

		}
	}()
	user := repository.GetUser(id)
	if user.IsEmpty() {
		return true, nil
	}
	_, err := db.Exec("delete from users where id = ?", id)
	return err == nil, err
}

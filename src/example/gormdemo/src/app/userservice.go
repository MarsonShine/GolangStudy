package appservice

import (
	"fmt"
	"gormdemo/src/contracts"
	"gormdemo/src/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 代表用户操作
type UserService struct {
	db *gorm.DB
}

// 创建用户
func (us UserService) Create(user *models.User) {
	if user == nil {
		panic("用户不存在！")
	} else {
		result := us.db.Create(user)
		if result.Error != nil && result.RowsAffected == 0 {
			panic("添加失败：" + result.Error.Error())
		}
	}
}

// 更新用户
func (us UserService) Update(user *models.User) bool {
	if user == nil {
		panic("用户不存在！")
	} else {
		// 查最新的用户
		var existUser models.User
		result := us.db.Where("email = ?", *user.Email).First(&existUser)
		if result.RowsAffected > 0 {
			existUser.Name = user.Name
			existUser.Age = user.Age
			existUser.Email = user.Email
			existUser.Birthday = user.Birthday
			us.db.Save(&existUser)
		} else {
			panic("用户不存在！")
		}

	}
	return true
}

// 获取用户
func (us UserService) Get(id uint) *models.User {
	var user models.User
	us.db.First(&user, id)
	return &user
}

// 删除用户
func (us UserService) Delete(id uint) bool {
	us.db.Delete(&models.User{}, id)
	return true
}

func CreateUserService(name string, email *string, age uint8, birthday *time.Time) {
	var mydb = openDbConnection()
	var user models.User
	result := mydb.First(&user, "email = ?", *email)
	if result.RowsAffected > 0 {
		panic(fmt.Sprintf("已经存在重复的用户 email:%s", *email))
	} else {
		createUser(UserService{db: mydb}, models.User{Name: name, Email: email, Age: age, Birthday: user.Birthday})
	}
}

// 获取所有用户
func GetUserAll() *[]models.User {
	var mydb = openDbConnection()
	var users []models.User
	result := mydb.Select([]string{}).Find(&users)
	if result.RowsAffected > 0 {
		return &users
	}
	return nil
}

// 查询指定用户
func GetUser(id int) models.User {
	var mydb = openDbConnection()
	return *(UserService{mydb}.Get(uint(id)))
}

func createUser(uo contracts.UserOperation, user models.User) {
	uo.Create(&user)
}

func openDbConnection() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./src/test.db"), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败！")
	}
	return db
}

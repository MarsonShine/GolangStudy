package appservice

import (
	"fmt"
	"gormdemo/src/contracts"
	"gormdemo/src/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateUserService(name string, email *string, age uint, birthday time.Time) {
	var db = openDbConnection()
	var user models.User
	result := db.First(&user, "email = ?", *email)
	if result.Error != nil {
		panic(fmt.Sprint("已经存在重复的用户 email:%s", *email))
	} else {
		contracts.UserService.Create(&models.User{Name: name, Email: email, Age: age, Birthday: &birthday})
	}
}

func openDbConnection() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./src/test.db"), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败！")
	}
	return db
}

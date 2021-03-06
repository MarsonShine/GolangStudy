package appservice

import (
	"gormdemo/src/contracts"
	"gormdemo/src/models"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 代表用户操作
type UserService struct {
	db *gorm.DB
}

func NewUserService(gormDB *gorm.DB) UserService {
	return UserService{db: gormDB}
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
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		us.db.Save(user)
		// 查最新的用户
		// var existUser models.User
		// result := us.db.Where("email = ?", *user.Email).First(&existUser)
		// if result.RowsAffected > 0 {
		// 	existUser.Name = user.Name
		// 	existUser.Age = user.Age
		// 	existUser.Email = user.Email
		// 	existUser.Birthday = user.Birthday
		// 	us.db.Save(&existUser)
		// } else {
		// 	panic("用户不存在！")
		// }

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

func CreateUserService(gormDB *gorm.DB, name string, email *string, age uint8, birthday *time.Time) {
	// var mydb = openDbConnection()
	// sqlDB, _ := mydb.DB()
	// defer sqlDB.Close()
	var user models.User
	// result := mydb.First(&user, "email = ?", *email)
	// if result.RowsAffected > 0 {
	// 	panic(fmt.Sprintf("已经存在重复的用户 email:%s", *email))
	// } else {
	// 	createUser(UserService{db: mydb}, models.User{Name: name, Email: email, Age: age, Birthday: user.Birthday})
	// }
	createUser(UserService{db: gormDB}, models.User{Name: name, Email: email, Age: age, Birthday: user.Birthday})
}

func UpdateUserService(gormDB *gorm.DB, id uint, name string, email *string, age uint8, birthday *time.Time) {
	updateUser(UserService{db: gormDB}, models.User{ID: id, Name: name, Email: email, Age: age, Birthday: birthday})
}

// 获取所有用户
func GetUserAll() *[]models.User {
	var mydb = openDbConnection()
	sqlDB, _ := mydb.DB()
	defer sqlDB.Close()
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

func DeleteUserByUserID(id uint) bool {
	var mydb = openDbConnection()
	// sqlDB, _ := mydb.DB()
	// defer sqlDB.Close()
	return UserService{mydb}.Delete(id)
}

func createUser(uo contracts.UserOperation, user models.User) {
	uo.Create(&user)
}

func updateUser(uo contracts.UserOperation, user models.User) {
	uo.Update(&user)
}

func openDbConnection() *gorm.DB {
	dsn := "root:123456@tcp(192.168.3.10:3306)/go_testdb?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(sqlite.Open("./src/test.db"), &gorm.Config{})
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(200)
	if err != nil {
		panic("连接数据库失败！")
	}
	return db
}

var (
	once       sync.Once
	dbInstance *gorm.DB
)

func GetDbConnection() *gorm.DB {
	once.Do(func() {
		dsn := "root:123456@tcp(192.168.3.10:3306)/go_testdb?charset=utf8mb4&parseTime=True&loc=Local"
		// db, err := gorm.Open(sqlite.Open("./src/test.db"), &gorm.Config{})
		var err error
		dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("连接数据库失败！")
		}
	})
	return dbInstance
}

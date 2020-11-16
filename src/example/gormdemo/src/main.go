package main

import (
	"encoding/json"
	appservice "gormdemo/src/app"
	"gormdemo/src/models"
	"net/http"
)

func main() {
	// db, err := gorm.Open(sqlite.Open("./src/test.db"), &gorm.Config{})
	// if err != nil {
	// 	panic("连接数据库失败！")
	// }

	// // migration
	// db.AutoMigrate(&models.Product{})

	// // Create
	// db.Create(&models.Product{Code: "D42", Price: 100})
	// // Read
	// var product models.Product
	// db.First(&product, 1)
	// db.First(&product, "code = ?", "D42") // 查找 code 字段为 D42

	// // Upload
	// db.Model(&product).Update("Price", 200)
	// // 更新多个字段
	// db.Model(&product).Updates(models.Product{Price: 200, Code: "F42"})
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	// db.Delete(&product, 1)

	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var u models.User
		err := decoder.Decode(&u)
		if err != nil {
			panic(err)
		}
		appservice.CreateUserService(u.Name, u.Email, u.Age, *u.Birthday)
		var jsonResponse = []byte(`{"sucess":true, "message": "success!"}`)
		w.Header().Set("content-type", "text/json")
		w.Write(jsonResponse)
	})
}

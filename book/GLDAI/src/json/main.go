package main

import "encoding/json"

type User struct {
	UserName string
	Age      int16
	Sex      bool
}

// json 序列化映射
type UserMap struct {
	UserName string `json:"userName"`
	Age      int16  `json:"age,string,omitempty"` //后面的 string,omitempty 是常用的 json 标签选项。分别标识序列化时以字符串形式标识，omitempty 表示在字段值为空时的一种默认形式
	Sex      bool
}

func main() {
	var user = User{
		UserName: "marsonshine",
		Age:      28,
		Sex:      true,
	}
	buffer, err := json.Marshal(user)
	if err == nil {
		println(string(buffer))
	}

	var userMap = UserMap{
		UserName: "marsonshine",
		Age:      28,
		Sex:      true,
	}

	buffer, err = json.Marshal(userMap)
	if err == nil {
		println(string(buffer))
	}

	var userMap2 = UserMap{
		UserName: "marsonshine",
	}

	buffer, err = json.Marshal(userMap2)
	if err == nil {
		println(string(buffer))
	}
}

package main

import (
	"fmt"
	"reflect"
)

/*
通过反射来获取对象的自定义tag内容，丰富反射的行为
json,orm等都用到了tag
*/

const tagName string = "validate"

type User struct {
	Id    int    `validate:"-"`
	Name  string `validate:"string,min=2,max=25"`
	Age   int    `validate:"number,min=18,max=60"`
	Email string `validate:"email"`
}

func main() {
	user := User{
		Id:    1,
		Name:  "Marson Shine",
		Age:   17,
		Email: "MarsonShine@example",
	}

	typeInfo := reflect.TypeOf(user)
	fmt.Println("Type:", typeInfo.Name())
	fmt.Println("Kind:", typeInfo.Kind())

	// 获取字段属性
	for i := 0; i < typeInfo.NumField(); i++ {
		field := typeInfo.Field(i)
		// 获取tag
		tag := field.Tag.Get(tagName)
		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
	}

	// 能获取tag，接下来就是解析tag
	user2 := User{
		Id:    0,
		Name:  "superlongstring",
		Age:   15,
		Email: "marsonshine",
	}
	for i, err := range ValidateStruct(user2) {
		fmt.Printf("\t%d. %s\n", i+1, err.Error())
	}
}

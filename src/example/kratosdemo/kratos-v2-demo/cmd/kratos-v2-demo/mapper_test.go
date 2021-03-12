package main

import (
	"kratos-v2-demo/internal/biz"
	"kratos-v2-demo/internal/data"
	"testing"

	"github.com/devfeel/mapper"
	"github.com/jinzhu/copier"
)

func TestMapper(t *testing.T) {
	createMapper()

	user := &data.UserEntity{}
	userDto := &biz.UserDto{
		UserName: "marsonshine",
		Age:      27,
		Address:  "创新大厦10楼",
	}
	mapper.Mapper(userDto, user)
	if user.Address == userDto.Address &&
		user.Age == userDto.Age {
		t.Logf("success!")
	}
}

func TestMapperTag(t *testing.T) {
	createMapper()
	user := &data.UserEntity{}
	userDto := &biz.UserDto{
		UserName: "marsonshine",
		Age:      27,
		Address:  "创新大厦10楼",
	}
	mapper.Mapper(userDto, user)
	if user.Name == userDto.UserName &&
		user.Age == userDto.Age {
		t.Logf("success!")
	} else {
		t.Errorf("error!")
	}
}

func TestCopier(t *testing.T) {
	user := data.UserEntity{}
	userDto := biz.UserDto{
		UserName: "marsonshine",
		Age:      27,
		Address:  "创新大厦10楼",
	}
	copier.Copy(user, userDto)
	if user.Name == userDto.UserName &&
		user.Age == userDto.Age {
		t.Logf("success!")
	} else {
		t.Errorf("error!")
	}
}

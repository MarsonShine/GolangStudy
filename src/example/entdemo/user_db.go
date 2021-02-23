package main

import (
	"context"
	"entdemo/ent"
	"entdemo/ent/user"
	"fmt"
	"log"
	"time"
)

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	sex := false
	u, err := client.User.
		Create().
		SetName("marsonshine").
		SetAge(27).
		SetAddress("深圳市南山区桃园街道创新大厦").
		SetNillableSex(&sex).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("添加用户失败：%v", err)
	}
	log.Printf("创建用户成功：%v", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.NameEQ("marsonshine")).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %v", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

func DeleteUser(ctx context.Context, client *ent.Client, id int) (bool, error) {
	d := client.User.DeleteOneID(id)
	if err := d.Exec(ctx); err != nil {
		return false, fmt.Errorf("failed deleting user: %v", err)
	}
	return true, nil
}

func DeleteUserByName(ctx context.Context, client *ent.Client, name string) (bool, error) {
	n, err := client.User.
		Delete().
		Where(user.NameEQ(name)).
		Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("failed deleting user: %v", err)
	}
	return n > 0, nil
}

func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %v", err)
	}
	log.Println("car was created: ", tesla)

	// 创建一个新的车品牌
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %v", err)
	}
	log.Println("car was created: ", ford)

	// 创建新用户，拥有多辆车
	marsonshine, err := client.User.
		Create().
		SetAge(27).
		SetName("marsonshine").
		AddCars(tesla, ford).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating user: %v", err)
	}
	log.Println("user was created: ", marsonshine)
	return marsonshine, nil
}

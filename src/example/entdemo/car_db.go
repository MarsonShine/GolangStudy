package main

import (
	"context"
	"entdemo/ent"
	"entdemo/ent/car"
	"fmt"
	"log"
)

func QueryCars(ctx context.Context, marsonshine *ent.User) error {
	cars, err := marsonshine.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %v", err)
	}

	log.Println("returned cars:", cars)

	// What about filtering specific cars.
	ford, err := marsonshine.QueryCars().
		Where(car.ModelEQ("Ford")).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %v", err)
	}
	log.Println(ford)
	return nil
}

func QueryCarUsers(ctx context.Context, marsonshine *ent.User) error {
	cars, err := marsonshine.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %v", err)
	}
	// 查询反向关系
	for _, ca := range cars {
		owner, err := ca.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying car %q owner: %v", ca.Model, err)
		}
		log.Printf("car %q owner: %q\n", ca.Model, owner.Name)
	}
	return nil
}

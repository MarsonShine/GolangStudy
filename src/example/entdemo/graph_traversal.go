package main

import (
	"context"
	"entdemo/ent"
	"log"
	"time"
)

func CreateGraph(ctx context.Context, client *ent.Client) error {
	summerzhu, err := client.User.
		Create().
		SetAge(30).
		SetName("summerzhu").
		Save(ctx)
	if err != nil {
		return err
	}
	marson, err := client.User.
		Create().
		SetAge(28).
		SetName("marson").
		Save(ctx)
	if err != nil {
		return err
	}
	// Then, create the cars, and attach them to the users in the creation.
	_, err = client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()). // ignore the time in the graph.
		SetOwner(summerzhu).         // attach this graph to Ariel.
		Save(ctx)
	if err != nil {
		return err
	}
	_, err = client.Car.
		Create().
		SetModel("Mazda").
		SetRegisteredAt(time.Now()). // ignore the time in the graph.
		SetOwner(summerzhu).         // attach this graph to Ariel.
		Save(ctx)
	if err != nil {
		return err
	}
	_, err = client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()). // ignore the time in the graph.
		SetOwner(marson).            // attach this graph to Neta.
		Save(ctx)
	if err != nil {
		return err
	}
	// Create the groups, and add their users in the creation.
	_, err = client.Group.
		Create().
		SetName("GitLab").
		AddUsers(marson, summerzhu).
		Save(ctx)
	if err != nil {
		return err
	}
	_, err = client.Group.
		Create().
		SetName("GitHub").
		AddUsers(summerzhu).
		Save(ctx)
	if err != nil {
		return err
	}
	log.Println("The graph was created successfully")
	return nil
}

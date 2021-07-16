package main

import "fmt"

func startg() {
	go func() {
		fmt.Println("start g")
	}()
}

package main

import (
	"fmt"
	"sync"
)

var once sync.Once

func main() {
	pool()

	for i := 0; i < 10; i++ {
		once.Do(func() {
			fmt.Println("execed...")
		})
	}

	duplicate(once)
}

func duplicate(once sync.Once) {
	for i := 0; i < 10; i++ {
		once.Do(func() {
			fmt.Println("execed2...")
		})
	}
}

func pool() {
	myPool := &sync.Pool{
		// New: func() interface{} {
		// 	fmt.Println("Creating new instance.")
		// 	return struct{}{}
		// },
	}
	o := myPool.Get() //1
	myPool.Put(1)
	fmt.Println(o)
	// instance := myPool.Get() //1
	// myPool.Put(instance)     //2
	// myPool.Get()             //3
}

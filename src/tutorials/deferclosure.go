package main

import "fmt"

func deferClosure() {
	var whatever [5]struct{}
	for i := range whatever {
		fmt.Println(i)
	}

	for i := range whatever {
		defer func() { fmt.Println(i) }()
	}

	for i := range whatever {
		defer func(n int) { fmt.Println(n) }(i)
	}
}

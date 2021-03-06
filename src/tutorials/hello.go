package main

// 关于迭代闭包取值的问题

import "fmt"

func hello(a int) int {
	c := a + 2
	return c
}

func hello2() int {
	a := 2
	b := 3
	return a + b
}

func main() {
	// hello2()
	deferClosure2()
}

func deferClosure2() {
	var whatever [5]struct{}
	for i := range whatever {
		fmt.Println(i)
	}

	for i := range whatever {
		// j := i 修复
		defer func() { fmt.Println(i) }()
	}

	for i := range whatever {
		defer func(n int) { fmt.Println(n) }(i)
	}
}

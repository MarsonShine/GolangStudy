package main

import "fmt"

// range function
func Backward[E any](s []E) func(func(int, E) bool) {
	return func(yield func(int, E) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(i, s[i]) {
				return
			}
		}
	}
}

func main() {
	// s := []string{"hello", "world"}
	// for i, x := range Backward(s) {
	// 	println(i, x)
	// }

	// 循环变量闭包问题
	// for i := range 10 {
	// 	println(10 - i)
	// }

	ch := make(chan bool)
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := range s {
		go func() {
			println(10 - i)
			ch <- true
		}()
	}

	for range s {
		<-ch
	}

	fmt.Println("go1.22 has lift-off!")
}

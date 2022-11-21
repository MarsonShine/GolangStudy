package main

import (
	"fmt"
)

func main() {
	// 无缓冲队列
	chOdd := make(chan int)
	chEven := make(chan int)
	source := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	go func(ch <-chan int) {
		for val := range ch {
			fmt.Println("奇数：", source[val])
		}
	}(chOdd)

	go func(ch <-chan int) {
		for val := range ch {
			fmt.Println("偶数：", source[val])
		}
	}(chEven)

	for i := range source {
		if i%2 == 0 {
			chOdd <- i
		} else {
			chEven <- i
		}
	}
}

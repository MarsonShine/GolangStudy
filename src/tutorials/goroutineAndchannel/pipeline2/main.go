package main

import "fmt"

// channel 分成不同的小功能函数，out, int

func main() {
	naturals := make(chan int)
	squarers := make(chan int)

	go counter(naturals)
	go squarer(squarers, naturals)

	printer(squarers)

}

// out
func counter(out chan<- int) {
	for x := 0; x < 100; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

// out
func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

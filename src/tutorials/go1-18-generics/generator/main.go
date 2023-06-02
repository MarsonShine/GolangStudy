package main

import "fmt"

type Addable interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | complex64 | complex128
}

func generator[T Addable](a T, b T) func() T {
	return func() T {
		t := a
		a += b
		return t
	}
}

func main() {
	f := generator(0, 1)
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
}

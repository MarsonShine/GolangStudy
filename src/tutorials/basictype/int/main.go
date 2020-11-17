package main

import "fmt"

func main() {
	var u uint8 = 255
	fmt.Println(u, u+1, u*u) // "255 0 1"

	var i int8 = 127
	fmt.Println(i, i+1, i*i) // "127 -128 1"

	var x uint8 = 1<<1 | 1<<5
	var y uint8 = 1<<1 | 1<<2

	fmt.Printf("%08b\n", x) // "00100010", the set {1, 5}
	fmt.Printf("%08b\n", y) // "00000110", the set {1, 2}

	fmt.Printf("%08b\n", x&y)  // "00000010", the intersection {1}
	fmt.Printf("%08b\n", x|y)  // "00100110", the union {1, 2, 5}
	fmt.Printf("%08b\n", x^y)  // "00100100", the symmetric difference {2, 5}
	fmt.Printf("%08b\n", x&^y) // "00100000", the difference {5}

	// 不同类型的相加会报错
	var apples int32 = 1
	var oranges int16 = 2
	var compute int
	// compute = apples + oranges // 报错
	compute = int(apples) + int(oranges)
	fmt.Printf("%d\n", compute)
}

package main

import (
	"fmt"
	"math"
)

//一个float32类型的浮点数可以提供大约6个十进制数的精度，而float64则可以提供约15个十进制数的精度；通常应该优先使用float64类型，因为float32类型的累计计算误差很容易扩散，并且float32能精确表示的正整数并不是很大（译注：因为float32的有效bit位只有23个，其它的bit位用于指数和符号；当整数大于23bit能表达的范围时，float32的表示将出现误差）：

func main() {
	var f float32 = 16777216 // 1 << 24
	fmt.Println(f == f+1)    // "true"!

	maxFloatValue := math.MaxFloat32
	fmt.Printf("%f\n", maxFloatValue)
}

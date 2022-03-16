//https://go.dev/doc/tutorial/generics
package main

import "fmt"

// 非范型版本
func SumInts(m map[string]int64) int64 {
	var s int64
	for _, v := range m {
		s += v
	}
	return s
}

func SumFloats(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}

// 有什么问题？其实就是入参和出参类型不一样之外，其他的全都一样
// 从调用端更能看出问题

func main() {
	ints := map[string]int64{
		"first":  34,
		"second": 16,
	}
	floats := map[string]float64{
		"first":  34.55,
		"second": 16.55,
	}
	fmt.Printf("非范型版本，总和：%v and %v\n", SumInts(ints), SumFloats(floats))
	// 类型推断，不用显式的申明string,int64,float64
	fmt.Printf("范型版本,总和：%v and %v\n", Sum(ints), Sum(floats))
	fmt.Printf("范型版本,总和：%v and %v\n", SumOfNumber(ints), SumOfNumber(floats))
}

// 这个时候可以定一个范型版本的函数，代替上面的两个仅入参和出参不同的函数
// 范型：每个类型参数都有一个参数约束，类似（act as）类型参数的一种元（metadata）类型
// 每个类型约束指定调用代码时每个输入参数的具体允许的类型
// 并且编译器还会根据用户调用范型函数传的参数，自动进行类型推断，这样就可以避免显式的申明参数类型
func Sum[K comparable, V int64 | float64](m map[K]V) V { // {V int64 | float64}表明参数类型约束
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// 显式定义范型约束
// 定义一个数字类型范型约束,这样就可以直接用 TNumber 代替 int64 | float64
type TNumber interface {
	int64 | float64
}

func SumOfNumber[K comparable, V TNumber](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

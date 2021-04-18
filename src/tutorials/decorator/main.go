package main

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

func decorator(f func(s string)) func(s string) {
	return func(s string) {
		fmt.Println("Started")
		f(s)
		fmt.Println("Done")
	}
}

func Hello(s string) {
	fmt.Println(s)
}
func main() {
	decorator(Hello)("Hello, World!")
	// 等价于
	hello := decorator(Hello)
	hello("Hello")

	sum1 := timedSumFunc(Sum1)
	sum2 := timedSumFunc(Sum2)
	fmt.Printf("%d, %d\n", sum1(-10000, 10000000), sum2(-10000, 10000000))

	// http.HandleFunc("/v1/hello", http1.WithServerHeader(http1.Hello))          // 这样写有点不好看，可以利用 Pipeline 优化
	// http.HandleFunc("/v2/hello", Handler(http1.Hello, http1.WithServerHeader)) // 可读性相对更好
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }

	// 范型 decorator
	type MyFoo func(int, int, int) int
	var myfoo MyFoo
	Decorator(&myfoo, foo)
	myfoo(1, 2, 3)

	mybar := bar
	Decorator(&mybar, bar)
	mybar("hello,", "world!")
}
func foo(a, b, c int) int {
	fmt.Printf("%d, %d, %d \n", a, b, c)
	return a + b + c
}
func bar(a, b string) string {
	fmt.Printf("%s, %s \n", a, b)
	return a + b
}

type SumFunc func(int64, int64) int64

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name() // 获取目标对象函数名
}
func timedSumFunc(f SumFunc) SumFunc {
	return func(start, end int64) int64 {
		defer func(t time.Time) {
			fmt.Printf("--- Time Elapsed (%s): %v ---\n", getFunctionName(f), time.Since(t))
		}(time.Now())
		return f(start, end)
	}
}

func Sum1(start, end int64) int64 {
	var sum int64
	sum = 0
	if start > end {
		start, end = end, start
	}
	for i := start; i < end; i++ {
		sum += i
	}
	return sum
}

func Sum2(start, end int64) int64 {
	if start > end {
		start, end = end, start
	}
	return (end - start + 1) * (end + start) / 2
}

// Pipeline
type HttpHandlerDecorator func(http.HandlerFunc) http.HandlerFunc

func Handler(h http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
	for i := range decors {
		d := decors[len(decors)-1-i] // 迭代反转
		h = d(h)
	}
	return h
}

// 利用反射构造范型decorator
func Decorator(decoPtr, fn interface{}) (err error) {
	var decoratedFunc, targetFunc reflect.Value
	decoratedFunc = reflect.ValueOf(decoPtr).Elem()
	targetFunc = reflect.ValueOf(fn)
	// reflect.MakeFunc 函数制出了一个新的函数
	v := reflect.MakeFunc(targetFunc.Type(), func(in []reflect.Value) (out []reflect.Value) {
		fmt.Println("before")
		out = targetFunc.Call(in)
		fmt.Println("after")
		return
	})
	decoratedFunc.Set(v) // v 赋值给 decoratedFunc
	return
}

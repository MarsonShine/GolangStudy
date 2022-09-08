package main

import (
	"fmt"
	"time"
)

func foobyval(n int) {
	fmt.Println(n)
}

func main() {
	for i := 0; i < 5; i++ {
		go func() {
			foobyval(i) // goroutine 闭包捕捉外部变量，编译器也会发出警告
		}()
	}
	time.Sleep(100 * time.Millisecond)
}

/*
为什么在循环闭包中捕捉外部值变量（非引用对象），还是输出相同的值？

*/

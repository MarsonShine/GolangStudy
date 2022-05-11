package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {
	// errorCode()
	correct()

	// for defer
	defer fmt.Println("first defer")
	for i := 0; i < 3; i++ {
		defer fmt.Printf("defer in for [%d]\n", i)
	}
	defer fmt.Println("last defer") // last defer  defer in for 2 1 0 first defer

	// goroutineError()
	// goroutineCorrect()
	goroutineWrap()
}

func errorCode() {
	fmt.Println("Enter function main.")
	// 引发 panic。
	panic(errors.New("something wrong"))
	p := recover() // 上一步引发panic，后续的代码无法继续执行
	fmt.Printf("panic: %s\n", p)
	fmt.Println("Exit function main.")
}

func correct() {
	fmt.Println("Enter function main.")
	defer func() {
		fmt.Println("Enter defer function.")
		if p := recover(); p != nil {
			fmt.Printf("panic: %s\n", p)
		}
		fmt.Println("Exit defer function.")
	}()
	// 引发 panic。
	panic(errors.New("something wrong"))
	fmt.Println("Exit function main.")
}

func goroutineError() {
	// go panic recover
	fmt.Println("Enter function main.")
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic: %s\n", p)
		}
	}()
	go func() {
		fmt.Println("Enter goroutine.")
		// 引发 panic。
		panic(errors.New("something wrong")) // 异常捕捉写到协程范围外是无法捕捉的
	}()
	time.Sleep(time.Second)
}

func goroutineCorrect() {
	// go panic recover
	fmt.Println("Enter function main.")
	go func() {
		defer func() {
			if p := recover(); p != nil {
				fmt.Printf("panic in goroutine: %s\n", p)
			}
		}()
		fmt.Println("Enter goroutine.")
		// 引发 panic。
		panic(errors.New("something wrong")) // 异常捕捉写到协程范围外是无法捕捉的
	}()
	time.Sleep(time.Second)
}

func goroutineWrap() {
	// go panic recover
	fmt.Println("Enter function main.")
	go func() {
		catch()
		fmt.Println("Enter goroutine.")
		// 引发 panic。
		panic(errors.New("something wrong")) // 异常捕捉封装到平级公共函数也是无法捕捉的
	}()
	time.Sleep(time.Second)
}

func catch() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic in goroutine: %s\n", p)
		}
	}()
}

package main

import (
	"fmt"
	"os"
)

func main() {
	// os.Args 第一个参数是程序的路径，后面的才是参数
	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	arg := os.Args[3]

	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(arg)
}

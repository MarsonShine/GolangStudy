package main

import (
	"flag"
	"fmt"
)

func main() {
	// go build -o 中的 -o 就是 flag
	wordPtr := flag.String("word", "foo", "a string")

	// 申明 numb 和 fork
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	// 所有标志申明完之后，调用 flag.Parse 来执行命令行解析
	flag.Parse()

	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("svar:", svar)
	// tail 指尾随的参数，尾随的位置参数可以出现在任何标志后面。
	fmt.Println("tail:", flag.Args())

	// 运行命令行 main.exe -h 程序就会自动生成命令行帮助文本提示
}

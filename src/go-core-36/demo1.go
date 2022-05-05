package main

import (
	"flag"
	"fmt"
	"os"
)

var name string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Useage of %s:\n", "question")
		flag.PrintDefaults()
	}
	flag.StringVar(&name, "name", "everyone", "请输入-name=value") // 第三个参数是没有输入对应的参数，则默认everyone，第四个参数是使用说明。
	flag.Parse()
}

func main() {
	fmt.Printf("Hello, %s!\n", name)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	linefilter()
}

// 行过滤器
// 对输入进行操作，返回返回给输出
// grep 就是这个效果
func linefilter() {
	scanner := bufio.NewScanner(os.Stdin)

	// 将输入的字符全部转化为大写并输出
	for scanner.Scan() {
		ucl := strings.ToUpper(scanner.Text())
		fmt.Println(ucl)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

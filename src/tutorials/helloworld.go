package main

import (
	"fmt"
	"os"
	"strconv"
)

func main1() {
	fmt.Println("Hello, 世界！")

	var s, seq string
	for i := 1; i < len(os.Args); i++ {
		s += seq + os.Args[i]
		seq = ""
	}
	fmt.Println(s)

	// for 表示 while
	for len(s) > 10 {
		s = s[0 : len(s)-1]
		fmt.Print(s + " ")
	}
	fmt.Println()

	// for 表示 while(true)
	// for {
	// 	// 无限循环
	// }

	var s1, sep = "", ""
	// for 表示范围
	for i, arg := range os.Args[1:] {
		is := strconv.Itoa(i)
		s1 += "索引=" + is + " " + sep + arg
		sep = " ~ "
	}
	fmt.Println(s1)

	// switch case
	// 在 go 语言中，switch case 中不必要标明 break，语言执行完 case 就会默认退出
	// 那么如果想要多个case执行同一个逻辑，则需要添加 fallthrough 语句覆盖默认的行为即可。

}

func SwitchCase(x int) int {
	switch {
	case x > 1:
		return x + 1
	default:
		return 0
	case x < 1:
		fallthrough // 表示继续执行后面的case
	case x == 1:
		return -1
	}
}

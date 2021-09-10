package main

import (
	"fmt"
	"time"
)

func main() {
	// 一个 case 可以判断多个表达式，用逗号分隔
	switch time.Now().Weekday() {
	case time.Friday:
		fmt.Println("今天是周五")
	case time.Monday:
		fmt.Println("今天是周一")
	case time.Tuesday:
		fmt.Println("今天是周二")
	case time.Wednesday:
		fmt.Println("今天是周三")
	case time.Thursday:
		fmt.Println("今天是周四")
	case time.Saturday, time.Sunday: // 多个表达式用逗号分隔
		fmt.Println("今天是周末")
	}

	// switch 可以不带表达式
	// 相当于 if else
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("12点之前")
	default:
		fmt.Println("12点了，中午了")
	}

	// 重点：还可以类型断言；类型开关(type switch)
	// 一般用作判断 interface 的具体类型。
	whatAmI := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			fmt.Println("布尔类型")
		case int:
			fmt.Println("int类型")
		case string:
			fmt.Println("字符串类型")
		default:
			fmt.Printf("未知类型 %T\n", t)
		}
	}

	whatAmI(true)
	whatAmI("marsonshine")
	whatAmI(100)
	whatAmI(123.99)
}

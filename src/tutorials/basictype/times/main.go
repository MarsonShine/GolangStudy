package main

import (
	"fmt"
	"time"
)

/*
	如果要做全球化跨时区的应用，你一定要把所有服务器和时间全部使用UTC时间。
	间一定要遵循：2006 为年，15 为小时，Monday 代表星期几等规则。
*/

func main() {
	p := fmt.Println
	t := time.Now()
	p(t.Format(time.RFC3339))

	t1, _ := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")
	p(t1)
	// 自定义布局
	p(t.Format("15:04PM"))
	p(t.Format("Mon Jan _2 15:04:05 2006"))
	p(t.Format("Mon Jan _2 15:04:05 2006"))
	p(t.Format("2006-01-02T15:04:05.999999-07:00"))

	p(t.Weekday().String())
}

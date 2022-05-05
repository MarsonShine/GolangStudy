// 自定义参数使用说明
package main

import (
	"flag"
)

func init() {
	flag.StringVar(&name, "name", "everyone", "请输入-name=value") // 第三个参数是没有输入对应的参数，则默认everyone，第四个参数是使用说明。
	flag.Parse()
}

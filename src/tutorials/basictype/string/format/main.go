package main

import "fmt"

type point struct {
	x, y int
}

func main() {
	p := point{1, 2}
	// v% 表示结构体对象
	fmt.Printf("%v\n", p)
	// v+% 将对象包括字段名内容也打印出来
	fmt.Printf("%+v\n", p)
	// v#% 打印该值的源码片段
	fmt.Printf("%#v", p)
	// %T 类型
	fmt.Printf("%T\n", p)
	// %t bool
	fmt.Printf("%t\n", true)
	// %d 格式化整形,十进制
	fmt.Printf("%d\n", 123)
	// %b 二进制
	fmt.Printf("%b\n", 14)
	// %c 格式化给定数字的ascii码值
	fmt.Printf("%c\n", 33)
	// %x 16进制
	fmt.Printf("%x\n", 456)
	// %f 浮点数
	fmt.Printf("%f\n", 78.9)
	// %e %E 科学计数法
	fmt.Printf("%e\n", 123400000.0)
	fmt.Printf("%E\n", 123400000.0)
	// %s 字符串输出
	fmt.Printf("%s\n", "\"string\"") // string
	// %q 像源码标记的一样，将带有双引号的字符串输出
	fmt.Printf("%q\n", "\"string\"") // "string"
	// %x 遇到字符串时，将使用 base-16 编码的字符串，每个字节使用2个字符表示
	fmt.Printf("%x\n", "hex this")
	// %p 指针
	fmt.Printf("%p\n", &p)
	// 显示数字要显示宽度（为了对其等格式要求），可以在"%"之后数字表示用空格对齐
	fmt.Printf("|%6d|%6d|\n", 12, 345)
	// %.2f 表示浮点数保留的精度
	fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45)
	// % 后 - 表示向左对齐
	fmt.Printf("|%-6.2f|%-6.2f|\n", 1.2, 3.45)
	// 字符串同上
	fmt.Printf("|%6s|%6s|\n", "foo", "b")
	fmt.Printf("|%-6s|%-6s|\n", "foo", "b")
}

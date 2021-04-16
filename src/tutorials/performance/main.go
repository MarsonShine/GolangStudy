package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	// 数字转字符串 strconv.Itoa() 要比 fmt.Sprintf() 快一倍
	_ = strconv.Itoa(20000000) // 直接转换为 ascii 码
	_ = fmt.Sprintf("%d", 2)
	// 尽可能避免 string 转成 []byte
	// 在 for-loop 中对 slice 做 append 操作，最佳实践是先初始化一个够用的 cap，原理与 C# 的 stringbuilder 初始化一个合适的 cap，来避免浪费内存和数据赋值。
	// stringbufer 或 stringbuild 拼接字符串会比使用 + 或 += 性能更好
	var b strings.Builder
	b.WriteString("ABC")
	b.WriteString("Hello World")
	fmt.Println(b.String())
	var sb bytes.Buffer
	sb.WriteString("DEF")
	sb.WriteString("Marson Shine")
	fmt.Println(sb.String())
	// 避免在热代码中进行内存分配，这样会导致gc很忙。尽可能的使用 sync.Pool 来重用对象。
	// 使用 I/O缓冲，I/O是个非常非常慢的操作，使用 bufio.NewWrite() 和 bufio.NewReader() 可以带来更高的性能。
	// 对于在for-loop里的固定的正则表达式，一定要使用 regexp.Compile() 编译正则表达式。性能会得升两个数量级。
	// 如果你需要更高性能的协议，你要考虑使用 protobuf 或 msgp 而不是JSON，因为JSON的序列化和反序列化里使用了反射。
	// 你在使用map的时候，使用整型的key会比字符串的要快，因为整型比较比字符串比较要快。
}

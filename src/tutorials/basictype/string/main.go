package main

// 关于字符串编码资料：https://blog.golang.org/strings
//第i个字节并不一定是字符串的第i个字符，因为对于非ASCII字符的UTF8编码会要两个或多个字节。我们先简单说下字符的工作方式。
import (
	"bytes"
	"fmt"
	"strconv"
	"unicode/utf8"
)

func main() {
	s := "hello, world"
	// s[i,j] 截取第i到第j（不包括j）
	fmt.Println(s[0:5]) // "hello"

	fmt.Println(s[:5]) // "hello"	相当于 s[0:5]
	fmt.Println(s[7:]) // "world"	相当于 s[7:len(s)]
	fmt.Println(s[:])  // "hello, world"  相当于 s[0:len(s)]
	//反引号里面的字符串内容是原生内容，没有转义操作。全部的内容都是字面的意思，包含退格和换行，因此一个程序中的原生字符串面值可能跨越多行
	const str = `Go is a tool for managing Go source code.

	Usage:
		go command [arguments]
	...`
	// 不同编码的字符串，len 得出的值不一样，字节码不一样
	ss := "hello, 世界"
	fmt.Println(len(ss))                    // "13"
	fmt.Println(utf8.RuneCountInString(ss)) // "9"
	// 要处理一直的信息，所以要解密
	for i := 0; i < len(ss); {
		r, size := utf8.DecodeRuneInString(ss[i:])
		fmt.Printf("%d\t%c\n", i, r)
		i += size
	}
	// 统计字符串的数量
	n := 0
	for range s {
		n++
	}

	n1 := utf8.RuneCountInString(ss)
	fmt.Println(n, n1)

	// 四个内置库 bytes、strings、strconv 和 unicode
	// 处理查询、替换、比较、截断、拆分和合并等功能

	// 字符串转换
	convert()

	const sample = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"
	fmt.Println(sample) // 这里会乱码
	// 为了输出真正的内容，可以下面几种方法
	// 方法一：分割开一个个字节输出
	for i := 0; i < len(sample); i++ {
		fmt.Printf("%x ", sample[i])
	}
	// 方法二：以正确的编码格式输出
	fmt.Printf("%x\n", sample)
	// 方式三：用转义符输出
	fmt.Printf("%q\n", sample)
}

func convert() {
	x := 123
	y := fmt.Sprintf("%d", x)
	fmt.Println(y, strconv.Itoa(x)) // "123 123"
	fmt.Println(strconv.FormatInt(int64(x), 2))
	// 字符串到数字
	xi, _ := strconv.Atoi("123")             // x is an int
	yi, _ := strconv.ParseInt("123", 10, 64) // base 10, up to 64 bits
	fmt.Printf("xi=%d,yi=%d", xi, yi)
}

// 注意，因为 string 是不变量，所以对字符串操作（分割，合并）会产生新的字符串，造成很多内存分配和数据复制
// 这个时候可以使用 bytes.Buffer 会更有效
// 一个Buffer开始是空的，但是随着string、byte或[]byte等类型数据的写入可以动态增长，一个bytes.Buffer变量并不需要初始化，因为零值也是有效的
func intsToString(values []int) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%d", v)
	}
	buf.WriteByte(']')
	return buf.String()
}

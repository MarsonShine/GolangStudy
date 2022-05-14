package main

import "fmt"

// Go 语言爱好者
// G o space 像这种英文字符都是英文字符，所对应的UTF-6编码值仅用一个字节即可。
// 而 语言爱好者 这写中文的十六进制范围都比较大，需要三个字节的编码值表示，8bed 8a00 7231 597d 8005 这些整数就是使用三个字节的编码值转换为整数的结果
// 每个字符的UTF-8的编码值还可以拆分为字节序列，如[% x]的输出就是如此：47 6f 20 e8 af ad e8 a8 80 e7 88 b1 e5 a5 bd e8 80 85

// 关于fmt.printf的输出格式可详见：https://pkg.go.dev/fmt
func main() {
	str := "Go 语言爱好者"
	fmt.Printf("The string: %q\n", str)                 // "Go 语言爱好者"
	fmt.Printf("  => runes(char): %q\n", []rune(str))   // ['G' 'o' ' ' '语' '言' '爱' '好' '者']
	fmt.Printf("  => runes(hex): %x\n", []rune(str))    // [47 6f 20 8bed 8a00 7231 597d 8005]	// 十六进制，8个十六进制分别代表8个字符
	fmt.Printf("  => runes(hex): [% x]\n", []byte(str)) // [47 6f 20 e8 af ad e8 a8 80 e7 88 b1 e5 a5 bd e8 80 85]

	// 遍历
	for i, c := range str {
		fmt.Printf("%d: %q [% x]\n", i, c, []byte(string(c))) // 中文字符一个rune打印出来了3个字节
	}
}

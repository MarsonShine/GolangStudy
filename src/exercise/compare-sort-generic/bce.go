//go:build ignore
// +build ignore

// 通过使用 -gcflags="-d=ssa/check_bce/debug=1" 编译器标志来显示是否需要边界检查信息
package main

//go run -gcflags="-d=ssa/check_bce/debug=1" bec.go
func f1(s []int) {
	_ = s[0] // 边界检查
	_ = s[1] // 边界检查
	_ = s[2] // 边界检查
}

func f2(s []int) {
	_ = s[2] // 边界检查
	_ = s[1] // 消除边界检查
	_ = s[0] // 消除边界检查
}

func f3(s []int, index int) {
	_ = s[index] // 边界检查
	_ = s[index] // 边界检查
}

func f4(a [5]int) {
	_ = a[4] // 消除边界检查
}

func f5(s []int) {
	for i := range s {
		_ = s[i]
		_ = s[i:len(s)]
		_ = s[:i+1]
	}
}

func f6(s []int) {
	for i := 0; i < len(s); i++ {
		_ = s[i]
		_ = s[i:len(s)]
		_ = s[:i+1]
	}
}

func f7(s []int) {
	for i := len(s) - 1; i >= 0; i-- {
		_ = s[i]
		_ = s[i:len(s)]
	}
}

func f8(s []int, index int) {
	if index >= 0 && index < len(s) {
		_ = s[index]
		_ = s[index:len(s)]
	}
}

func f9(s []int) {
	if len(s) > 2 {
		_, _, _ = s[0], s[1], s[2]
	}
}

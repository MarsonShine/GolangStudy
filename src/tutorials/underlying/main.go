package main

import (
	"fmt"
	"unsafe"
)

// 事实上，Go语言的调度器会自己决定是否需要将某个goroutine从一个操作系统线程转移到另一个操作系统线程。
// 一个指向变量的指针也并没有展示变量真实的地址。因为垃圾回收器可能会根据需要移动变量的内存位置，当然变量对应的地址也会被自动更新。
// unsafe.Sizeof 函数返回操作数在内存中的字节大小
// unsafe.Alignof 函数返回对应参数的类型需要对齐的倍数
// unsafe.Offsetof 函数的参数必须是一个字段 x.f，然后返回 f 字段相对于 x 起始地址的偏移量，
var a struct {
	bool
	float64
	int16
}
var b struct {
	float64
	int16
	bool
}
var c struct {
	bool
	int16
	float64
}

var x struct {
	a bool
	b int16
	c []int
}

func main() {
	fmt.Println(unsafe.Sizeof(a), unsafe.Sizeof(b), unsafe.Sizeof(c))                             // 24 16 16
	fmt.Println(unsafe.Sizeof(x), unsafe.Sizeof(x.a), unsafe.Sizeof(x.b), unsafe.Sizeof(x.c))     //
	fmt.Println(unsafe.Alignof(x), unsafe.Alignof(x.a), unsafe.Alignof(x.b), unsafe.Alignof(x.c)) //
	fmt.Println(unsafe.Offsetof(x.a), unsafe.Offsetof(x.b), unsafe.Offsetof(x.c))                 //
}

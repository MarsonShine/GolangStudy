package main

import (
	"fmt"
	"unsafe"
)

type (
	Dog struct {
		name string
		// category string
	}
)

func main() {
	dog := Dog{"斑点狗"}
	dogP := &dog
	dogPtr := uintptr(unsafe.Pointer(dogP))
	fmt.Printf("%p", &dogPtr)

	// dogPtr2 := uintptr(dogP) // 无法将指针值直接转换uintptr
	// 转换为uintptr之后，可以进行地址计算
	// unsafe.Offsetof 函数用于获取两个值在内存中的起始存储地址之间的偏移量，以字节为单位。
	namePtr := dogPtr + unsafe.Offsetof(dogP.name)
	nameP := (*string)(unsafe.Pointer(namePtr))
	fmt.Printf("%s", *nameP)
}

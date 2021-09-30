package main

import (
	"crypto/sha1"
	"fmt"
	"unsafe"
)

func main() {
	s := "sha1 this string"
	h := sha1.New()

	h.Write([]byte(s))

	bs := h.Sum(nil)
	fmt.Println(s)
	fmt.Printf("%x\n", bs)

	// 等价于
	bs2 := sha1.Sum([]byte(s))
	// 同样16进制显示
	fmt.Printf("%x\n", bs2)
}

func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

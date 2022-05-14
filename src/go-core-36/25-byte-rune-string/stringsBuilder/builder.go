package main

import (
	"fmt"
	"strings"
)

// strings.Builder相较于string有哪些优点？
// 1. 可以减少内存分配、数据复制，以降低GC的压力
// 2. 可以提高内存利用率：可以动态的添加/删减字符串

// strings.Builder 在写入字符串时会将字符串添加到内部一个连续地址空间的字节切片buf中
// 超过了内部预定的buf长度时，就会触发扩容：申请2倍之前的大小的空间，将数据拷贝到新申请的内存空间
func main() {
	sb := strings.Builder{}
	// 写入内容
	sb.WriteString("marsonshine")
	// 手动扩容
	sb.Grow(30)
	fmt.Printf("The length of contents in the builder is %d.\n", sb.Len())
}

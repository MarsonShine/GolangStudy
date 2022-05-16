package main

import (
	"bytes"
	"fmt"
)

// bytes.Buffer 高效读写字节，缓存区

func main() {
	var buffer1 bytes.Buffer
	contents := "Simple byte buffer for marshaling data."
	fmt.Printf("Writing contents %q ...\n", contents)
	buffer1.WriteString(contents)
	fmt.Printf("The length of buffer: %d\n", buffer1.Len()) // buffer.Len() 返回的是未读内容的长度
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap())

	p1 := make([]byte, 7)
	n, _ := buffer1.Read(p1)
	fmt.Printf("%d bytes were read. (call Read)\n", n)
	fmt.Printf("The length of buffer: %d\n", buffer1.Len())
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap())

	contents = "ab"
	buffer2 := bytes.NewBufferString(contents)
	fmt.Printf("The capacity of new buffer with contents %q: %d\n",
		contents, buffer1.Cap())
	fmt.Println()

	unreadBytes := buffer2.Bytes()
	fmt.Printf("The unread bytes of the buffer: %v\n", unreadBytes)
	fmt.Println()

	contents = "cdefg"
	fmt.Printf("Write contents %q ...\n", contents)
	buffer1.WriteString(contents)
	fmt.Printf("The capacity of buffer: %d\n", buffer1.Cap())
	fmt.Println()

	// 只要扩充一下之前拿到的未读字节切片unreadBytes，
	// 就可以用它来读取甚至修改buffer中的后续内容。
	unreadBytes = unreadBytes[:cap(unreadBytes)]
	fmt.Printf("The unread bytes of the buffer: %v\n", unreadBytes)
	fmt.Println()

	value := byte('X')
	fmt.Printf("Set a byte in the unread bytes to %v ...\n", value)
	unreadBytes[len(unreadBytes)-2] = value
	fmt.Printf("The unread bytes of the buffer: %v\n", buffer1.Bytes())
	fmt.Println()
}

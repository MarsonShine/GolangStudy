package main

import (
	"fmt"
	"io"
	"strings"
)

// 字符串高效读的最佳实践
func main() {
	var reader1 strings.Reader
	_ = reader1.Size() - int64(reader1.Len()) // 计算出的已读计数

	offset2 := int64(17)
	expectedIndex := reader1.Size() - int64(reader1.Len()) + offset2
	fmt.Printf("Seek with offset %d and whenc %d ...\n", offset2, io.SeekCurrent)
	readingIndex, _ := reader1.Seek(offset2, io.SeekCurrent)
	fmt.Printf("The reading index in reader: %d (returned by Seek)\n", readingIndex)
	fmt.Printf("The reading index in reader: %d (computed by me)\n", expectedIndex)
}

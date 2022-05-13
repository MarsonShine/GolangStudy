package main

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

var bufPool sync.Pool

// 数据块缓冲区
type Buffer interface {
	Delimiter() byte
	// 写一个数据块
	Write(content string) (err error)
	// 读一个数据块
	Read() (content string, err error)
	Free()
}

type myBuffer struct {
	buf       bytes.Buffer
	delimiter byte
}

func (b *myBuffer) Delimiter() byte {
	return b.delimiter
}

func (b *myBuffer) Write(content string) (err error) {
	if _, err = b.buf.WriteString(content); err != nil {
		return
	}
	return b.buf.WriteByte(b.delimiter)
}

func (b *myBuffer) Read() (content string, err error) {
	return b.buf.ReadString(b.delimiter)
}

func (b *myBuffer) Free() {
	bufPool.Put(b)
}

// delimiter 代表预定义的定界符。
var delimiter = byte('\n')

// sync.Pool不是开箱即用的，需要自己初始化
func init() {
	bufPool = sync.Pool{
		New: func() interface{} {
			return &myBuffer{delimiter: delimiter}
		},
	}
}

// GetBuffer 用于获取一个数据块缓冲区。
func GetBuffer() Buffer {
	return bufPool.Get().(Buffer)
}

func main() {
	// sync.Pool 可以避免多次new分配对象内存
	// 调试发现，10次调用，只有少量几次会调用New重新分配对象
	for i := 0; i < 10; i++ {
		buf := GetBuffer()
		defer buf.Free()
		buf.Write("A Pool is a set of temporary objects that" +
			"may be individually saved and retrieved.")
		buf.Write("A Pool is safe for use by multiple goroutines simultaneously.")
		buf.Write("A Pool must not be copied after first use.")

		fmt.Println("The data blocks in buffer:")
		for {
			block, err := buf.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(fmt.Errorf("unexpected error: %s", err))
			}
			fmt.Print(block)
		}
	}
}

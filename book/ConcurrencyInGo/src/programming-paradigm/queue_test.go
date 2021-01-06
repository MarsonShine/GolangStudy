package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// 队列不会减少程序的总运行时间，但是会改变程序的行为
// 队列提高了整体性能那种情况？
// 1：批处理，比如输入缓冲到内存而非硬盘中
func BenchmarkUnbufferedWrite(b *testing.B) {
	performWrite(b, tmpFileOrFatal())
}

func BenchmarkBufferedWirte(b *testing.B) {
	bufferedFile := bufio.NewWriter(tmpFileOrFatal())
	performWrite(b, bufferedFile)
}

func performWrite(b *testing.B, writer io.Writer) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for bt := range take(done, repeat(done, byte(0)), b.N) {
		writer.Write([]byte{bt.(byte)})
	}
}

func tmpFileOrFatal() *os.File {
	file, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return file
}

// 2:

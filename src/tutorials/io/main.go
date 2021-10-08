package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	readfile()
	writefile()
}

func readfile() {
	path, _ := filepath.Abs("./tmp/dat")
	data, err := ioutil.ReadFile(path)
	check(err)
	fmt.Println(string(data))

	f, err := os.Open(path)
	check(err)

	b1 := make([]byte, 5)
	n1, err := f.Read(b1)
	check(err)
	fmt.Printf("%d bytes: %s\n", n1, string(b1[:n1]))

	// 手动偏移位置
	o2, err := f.Seek(6, 0)
	check(err)
	b2 := make([]byte, 2)
	n2, err := f.Read(b2)
	check(err)
	fmt.Printf("%d bytes @ %d: ", n2, o2)
	fmt.Printf("%v\n", string(b2[:n2]))

	// ReadAtLeast 这种比上面更健壮
	o3, err := f.Seek(6, 0)
	check(err)
	b3 := make([]byte, 2)
	n3, err := io.ReadAtLeast(f, b3, 2)
	check(err)
	fmt.Printf("%d bytes @ %d: %s\n", n3, o3, string(b3))

	// 不适用 Seek 实现上面同样的功能
	// bufio 包实现了一个缓冲读取器，这可能有助于提高许多小读操作的效率，以及它提供了很多附加的读取函数。
	_, err = f.Seek(0, 0)
	check(err)

	r4 := bufio.NewReader(f)
	b4, err := r4.Peek(5)
	check(err)
	fmt.Printf("5 bytes: %s\n", string(b4))

	// 关闭
	f.Close()

}

func writefile() {
	d1 := []byte("hello\ngo\n")
	path, _ := filepath.Abs("./tmp/dat1")
	err := ioutil.WriteFile(path, d1, 0644)
	check(err)

	path, _ = filepath.Abs("./tmp/dat2")
	f, err := os.Create(path)
	check(err)
	defer f.Close()

	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(d2)
	check(err)
	fmt.Printf("wrote %d bytes \n", n2)

	n3, err := f.WriteString("writes\n")
	fmt.Printf("wrote %d bytes\n", n3)

	// 调用 Sync 将缓冲区的数据写入硬盘。
	f.Sync()

	w := bufio.NewWriter(f)
	n4, err := w.WriteString("buffered\n")
	fmt.Printf("wrote %d bytes\n", n4)
	// 使用 Flush 来确保，已将所有的缓冲操作应用于底层 writer。
	w.Flush()

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

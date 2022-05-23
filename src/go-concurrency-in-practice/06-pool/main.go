package main

import (
	"bytes"
	"sync"
)

var buffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func main() {

}

func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

// 有问题，会造成内存泄漏
// 因为传的参数是引用，如果这个buf尺寸越来越大，则底层的slice结构所占的空间就会越来越大。大到一定值的时候，gc就不会回收
// 这样就会导致了内存泄漏，解决方案就是在回收的时候判断buf的大小，内存太大就不要回收了
func PutBuffer(buf *bytes.Buffer) {
	// if cap(buf.Bytes()) > 1<<16 { // 判断是否大于64kb
	// 	return
	// }
	buf.Reset()
	buffers.Put(buf)
}

/*

 */

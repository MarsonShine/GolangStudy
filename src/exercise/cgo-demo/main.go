package main

/*
#include <stdlib.h>
#include "clibrary.h"

#cgo CFLAGS: -I .
#cgo LDFLAGS: -L . -lclibrary

extern void startCgo(void*, int);
extern void endCgo(void*, int, int);
*/
import "C"
import (
	"fmt"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

type Visitor interface {
	Start(int)
	End(int, int)
}

func GoTraverse(filename string, v Visitor) {
	cCallbacks := C.Callbacks{}

	cCallbacks.start = C.StartCallbackFn(C.startCgo)
	cCallbacks.end = C.EndCallbackFn(C.endCgo)

	// 分配一个C字符串来保存文件名的内容，并在指定时间释放它
	var cfilename *C.char = C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	//创建一个不透明的C指针，供访问者传递遍历。
	p := gopointer.Save(v)
	defer gopointer.Unref(p)

	// 调用c库的traverse函数
	C.traverse(cfilename, cCallbacks, p)
}

//export goStart
func goStart(user_data unsafe.Pointer, i C.int) {
	v := gopointer.Restore(user_data).(Visitor)
	v.Start(int(i))
}

//export goEnd
func goEnd(user_data unsafe.Pointer, a C.int, b C.int) {
	v := gopointer.Restore(user_data).(Visitor)
	v.End(int(a), int(b))
}

type MyVisitor struct {
	startState int
}

func (mv *MyVisitor) Start(i int) {
	mv.startState = i
	fmt.Println("End:", i)
}

func (mv *MyVisitor) End(a, b int) {
	fmt.Printf("Start: %v %v [state = %v]\n", a, b, mv.startState)
}

func main() {
	mv := &MyVisitor{startState: 0}
	GoTraverse("joe", mv)
}

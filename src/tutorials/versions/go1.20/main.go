package main

import (
	"fmt"
	"unsafe"
)

func main() {
	sliceToArray()

	n := 1
	var i interface{} = n // go 1.19 不支持这种写法，报interface{} do not implement comparable
	gcompareable(i)

	n1 := []byte{1}
	n2 := []byte{2}
	var x interface{} = n1
	var y interface{} = n2
	gcompareable2(x, y) // 切片类型不是可比较类型 error:runtime error: comparing uncomparable type []uint8
}

// go 1.20 新增语法糖：slice 和 array 互相转换支持
func sliceToArray() {
	var s1 = []int{1, 2, 3, 4, 5}
	var arr = [5]int(s1)
	arr[1] = 0
	fmt.Println(s1)
	fmt.Println(arr)

	s1[0] = 0
	fmt.Println(s1)
	fmt.Println(arr)

	var parr = (*[5]int)(s1) // 转换为指针，指针指向的是 s1切片的底层数组
	s1[0] = 0
	fmt.Println(s1)
	fmt.Println(arr)
	fmt.Println(parr)
}

// 泛型约束增强
func gcompareable[T comparable](t T) {

}

func gcompareable2[T comparable](t1 T, t2 T) {
	if t1 != t2 {
		return
	}
	println("equals")
}

func unsafeApi() {
	var arr = [6]byte{'h', 'e', 'l', 'l', 'o', '!'}
	s := unsafe.String(&arr[0], 6) // https://pkg.go.dev/unsafe#String
	fmt.Println(s)                 // hello!
	arr[0] = 'j'
	fmt.Println(s) // jello!

	s1 := "golang"
	fmt.Println(s1)            // golang
	b := unsafe.StringData(s1) // https://pkg.go.dev/unsafe#StringData 返回的字节是不允许被修改的
	*b = 'h'                   // fatal error: fault, unexpected fault address 0x10a67e5
	fmt.Println(s1)
}

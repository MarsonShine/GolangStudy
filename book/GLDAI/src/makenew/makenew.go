// make 和 new 都有初始化一个对象结构，但是两者有很大区别
// make: 只能初始化 go 内置的结构对象，如数组、切片、哈希表和 Channel
// new: 是根据传入结构类型分配一片内存空间并返回指向这片内存空间的指针

package main

import "fmt"

func main() {
	// slice 是一个包含 data、cap 和 len 的私有结构体 internal/reflectlite.sliceHeader；
	// hash 是一个指向 runtime.hmap 结构体的指针；
	// ch 是一个指向 runtime.hchan 结构体的指针；
	slice := make([]int, 0, 100)
	hash := make(map[int]string, 100)
	ch := make(chan int, 5)
	fmt.Println(slice, hash, ch)
	// new, 返回指针
	i := new(int)
	// 上面等价于下面
	var v int
	i = &v
	fmt.Printf("%p", i)
}

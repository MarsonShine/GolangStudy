package main

import (
	"bytes"
	"fmt"
	"reflect"
)

/*
slice 切片数据结构是
type slice struct {
	array unsafe.Pointer //指向存放数据的数组指针
	len int // 长度
	cap int	// 容量
}
里面存放的数据是指针，所以这是引用的数组的地址，这在操作数据的时候要注意，因为这是共享内存的。
*/

func main() {
	var foo = make([]int, 5)
	foo[3] = 42
	foo[4] = 100

	// 切片数据至另一个对象
	bar := foo[1:4] // index 1 ~ index 4
	bar[1] = 99

	// 由于 slice 里面存放数据的 array 对象是指针，是直接访问内存的。所以这里 bar 赋值对 foo 是有影响的
	fmt.Printf("bar[1] = 99, foo[1] = %d", foo[2])

	/*
		切片在用 append 追加数据时，如果此时 slice 的容量够用，那么就不会重新配分内存（容量翻倍），此时的 slice 的数据对象还是共享的
		如果不够用了，就会重新分配内存。
	*/
	foo = make([]int, 32) // cap = 32,len = 32
	// b := foo[1:16]
	foo = append(foo, 1) // len = 32+1 超出了 32，就会重新分配内容，此时 b，foo 的内存就不是共享的
	foo[2] = 42

	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/')
	dir1 := path[:sepIndex] // 因为 cap 够用，所以内存是共享的，此时如果要想不内存共享，则需要用 Full Slice Expression
	// 改成如下即是重新分配内存
	dir11 := path[:sepIndex:sepIndex]
	dir2 := path[sepIndex+1:]
	fmt.Println("dir1 =>", string(dir11)) //prints: dir1 => AAAA
	fmt.Println("dir1 =>", string(dir1))  //prints: dir1 => AAAA
	fmt.Println("dir2 =>", string(dir2))  //prints: dir2 => BBBBBBBBB
	dir1 = append(dir1, "suffix"...)
	fmt.Println("dir1 =>", string(dir1)) //prints: dir1 => AAAAsuffix
	fmt.Println("dir2 =>", string(dir2)) //prints: dir2 => uffixBBBB

	/*
		结构体的比较，是需要深度比较的，可以利用 reflect.DeepEqual(a,b)
	*/
	v1 := data{}
	v2 := data{}
	fmt.Println("v1 == v2:", reflect.DeepEqual(v1, v2))

	m1 := map[string]string{"one": "a", "two": "b"}
	m2 := map[string]string{"two": "b", "one": "a"}
	fmt.Println("v1 == v2:", reflect.DeepEqual(m1, m2))

	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	fmt.Println("s1 == s2:", reflect.DeepEqual(s1, s2))

	// 接口多态解耦
	d1 := Country{"USA"}
	d2 := City{"Los Angeles"}
	PrintStr(d1)
	PrintStr(d2)
	// 结构嵌套
	c1 := CountryName{WithName{"China"}}
	c2 := CityName{WithName{"ShenZhen"}}
	c1.PrintStr()
	c2.PrintStr()
}

type data struct {
}

/*
	接口实现的多种方式，面向接口编程，而不是具体实现
*/
type Country struct {
	Name string
}
type City struct {
	Name string
}
type Stringable interface {
	ToString() string
}

func (c Country) ToString() string {
	return "Country = " + c.Name
}
func (c City) ToString() string {
	return "City = " + c.Name
}
func PrintStr(s Stringable) {
	fmt.Println(s.ToString())
}

/*
也可以通过结构嵌套完成上面的目的，但是耦合度太强
*/
type WithName struct {
	Name string
}
type CountryName struct {
	WithName
}
type CityName struct {
	WithName
}

func (w WithName) PrintStr() {
	fmt.Println(w.Name)
}

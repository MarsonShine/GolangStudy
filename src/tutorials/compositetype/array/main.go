package main

import "fmt"

// 数组
// 长度是固定的，对内存布局友好，但是限制很多。因此在 Go 语言中使用比较少
// 与之对比的是 Slice，是可以动态变化的，可以增长和收缩

func main() {
	var a [3]int             // array of 3 integers
	fmt.Println(a[0])        // 打印第一个元素
	fmt.Println(a[len(a)-1]) // 打印第二个元素

	// 打印索引与元素
	for i, v := range a {
		fmt.Printf("%d %d\n", i, v)
	}

	// 只打印元素
	for _, v := range a {
		fmt.Printf("%d\n", v)
	}

	// 数组初始化
	var q [3]int = [3]int{1, 2, 3}
	var r [3]int = [3]int{1, 2}
	fmt.Println(r[2], q[2]) // "0"
	// ... 省略号，代表根据实际初始化元素个数来计算
	s := [...]int{1, 2, 3, 4, 5}
	fmt.Printf("%T\n", s) // "[5]int"

	// 可以指定索引和索引对应的值来初始化是数组
	symbol := [...]string{USD: "$", EUR: "€", GBP: "￡", RMB: "￥"}
	fmt.Println(RMB, symbol[RMB]) // "3 ￥"
	// 安装范围值初始化数组
	rg := [...]int{99: -1} // 定义含100个元素的数组 rg，最后一个元素初始化为 -1，其他元素都是用 0
	for i, v := range rg {
		fmt.Print(fmt.Sprintf("rg[%d]=%d", i, v), " ")
	}
	// 数组比较，要里面元素的值一样即为 true
}

type Currency int

const (
	USD Currency = iota // 美元
	EUR                 // 欧元
	GBP                 // 英镑
	RMB                 // 人民币
)

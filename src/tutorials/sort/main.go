package main

import (
	"fmt"
	"sort"
)

func main() {
	strs := []string{"c", "a", "b"}
	sort.Strings(strs)
	fmt.Println("内置排序，支持内置的数据类型 ", strs)

	ints := []int{1, 4, 2, 8}
	sort.Ints(ints)
	fmt.Println("整形排序 ", ints)

	// 表示是否已经是经过排序的顺序
	b := sort.IntsAreSorted(ints)
	fmt.Println("是否有序 ", b)

	customerSort()
}

// 自定义排序
// 主要是实现 sort.Interface 接口的 Len、Less、Swap 方法即可
// 从字符串长度由短到长排序
type byLength []string

// 实现接口方法
func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func customerSort() {
	fruits := []string{"peach", "banana", "kiwi"}
	sort.Sort(byLength(fruits))
	fmt.Println(fruits)
}

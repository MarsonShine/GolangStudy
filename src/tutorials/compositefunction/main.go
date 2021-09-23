package main

import (
	"fmt"
	"strings"
)

// 通过组装不同函数来实现满足给定条件的所有集合
// 这种模式对于集合操作非常有用并且常见

func Index(vs []string, s string) int {
	for i, v := range vs {
		if v == s {
			return i
		}
	}
	return -1
}

func Include(vs []string, s string) bool {
	// for _, v := range vs {
	// 	if v == s {
	// 		return true
	// 	}
	// }
	// return false
	// 组合
	return Index(vs, s) > -1
}

func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func Filter(vs []string, f func(string) bool) []string {
	arrs := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			arrs = append(arrs, v)
		}
	}
	return arrs
}

func Map(vss []string, f func(string) string) []string {
	arrs := make([]string, len(vss))
	for i, v := range vss {
		arrs[i] = f(v)
	}
	return arrs
}

func main() {
	var strs = []string{"peach", "apple", "pear", "plum"}
	fmt.Println(Index(strs, "pear"))
	fmt.Println(Include(strs, "grape"))
	fmt.Println(Any(strs, func(v string) bool {
		return strings.HasPrefix(v, "p")
	}))
	fmt.Println(All(strs, func(v string) bool {
		return strings.HasPrefix(v, "p")
	}))
	fmt.Println(Filter(strs, func(v string) bool {
		return strings.Contains(v, "e")
	}))

	fmt.Println(Map(strs, strings.ToUpper))
}

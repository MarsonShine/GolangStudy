package main

import "fmt"

var container = []string{"1", "2", "3"}

func main() {
	container := map[int]string{0: "1", 1: "2", 2: "3"}
	fmt.Printf("the element is %q.\n", container[0]) // 就近域范围

	// 类型推断
	value, ok := interface{}(container).([]string)
	if ok {
		fmt.Printf("container is []string:%v", value)
	}
	if value, ok := interface{}(container).(map[int]string); ok {
		fmt.Printf("container is map[int]string: %v", value)
	}
}

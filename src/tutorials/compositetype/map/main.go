package main

import (
	"bufio"
	"fmt"
	"os"
)

// map[key]value key 最好不要是float类型
func main() {
	ages := make(map[string]int)
	ages = map[string]int{
		"alice":   31,
		"charlie": 34,
	}

	// 可以删除
	delete(ages, "alice")         // 如果key不在，也不会报错
	ages["bob"] = ages["bob"] + 1 // 不存在 bob 也能正常运行
	ages["marsonshine"] = 27      //也可以直接对不存在的key复制
	printMap(ages)
	// map 中的元素不是变量，所以不能对其变量进行取址操作
	// _ = &ages["bob"] // 报错

	_ = make([]string, 0, len(ages)) // 0-len(ages) 的 string 数组

	// 如果key在map中是存在的，那么将得到与key对应的value；如果key不存在，那么将得到value对应类型的零值
	age, ok := ages["summerzhu"]
	if ok {
		fmt.Printf("存在该键：val=%d", age)
	}
	// 也可以直接合并
	if age1, ok1 := ages["summerzhu"]; !ok1 {
		fmt.Printf("不存在该键：val=%d", age1)
	}

	addEdge("a", "b")
	addEdge("c", "d")
	addEdge("a", "d")
	addEdge("d", "a")
	fmt.Println(hasEdge("a", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("a", "d"))
	fmt.Println(hasEdge("d", "a"))
	fmt.Println(hasEdge("x", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("x", "d"))
	fmt.Println(hasEdge("d", "x"))

	seen := make(map[string]bool) // a set of strings
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if !seen[line] {
			seen[line] = true
			fmt.Println(line)
		}
	}

	if err := input.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "dedup: %v\n", err)
		os.Exit(1)
	}
}

func printMap(kvs map[string]int) {
	for key, val := range kvs {
		fmt.Printf("%s\t%d\n", key, val)
	}
}

// map 比较是一定是通过判断ok是否为true，因为不存在的key其值默认是为0的
func equal(x, y map[string]int) bool {
	if len(x) != len(y) {
		return false
	}
	for k, xv := range x {
		if yv, ok := y[k]; !ok || yv != xv {
			return false
		}
	}
	return true
}

var graph = make(map[string]map[string]bool)

func addEdge(from, to string) {
	edge := graph[from]
	if edge == nil {
		edge = make(map[string]bool)
		graph[from] = edge
	}
	edge[to] = true
}

func hasEdge(from, to string) bool {
	return graph[from][to]
}

package main

import "fmt"

func main() {
	arr := []int{1, 2, 3}
	newArr := []*int{}

	for _, v := range arr {
		newArr = append(newArr, &v) // 引用地址覆盖，都会输出 3
	}

	for _, v := range newArr {
		fmt.Println(*v)
	}

	// 正确的形式
	for i := range arr {
		newArr = append(newArr, &arr[i])
	}
	for _, v := range newArr {
		fmt.Println(*v)
	}

	hash := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}
	// 运行多次，输出的内容都不同，这是因为哈希表引入了不稳定性
	for k, v := range hash {
		println(k, v)
	}
	// for 经典循环
	// for 循环初始化; 循环的条件; 循环结束时执行体; { 循环体 }
}

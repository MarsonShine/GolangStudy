package main

import "fmt"

func main() {
	// 数组值类型
	array1 := [3]string{"a", "b", "c"}
	fmt.Printf("The array: %v\n", array1)
	array2 := modifyArray(array1)
	fmt.Printf("The modified array: %v\n", array2)
	fmt.Printf("The original array: %v\n", array1)

	// 切片
	slices1 := []string{"a", "b", "c"}
	fmt.Printf("The slice: %v\n", slices1)
	slice2 := modifySlice(slices1)
	fmt.Printf("The modified array: %v\n", slice2)
	fmt.Printf("The original array: %v\n", slices1)
}

func modifyArray(a [3]string) [3]string {
	a[1] = "x"
	return a
}

func modifySlice(a []string) []string {
	a[1] = "x"
	return a
}

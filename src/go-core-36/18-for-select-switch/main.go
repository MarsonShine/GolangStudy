package main

import "fmt"

func main() {
	numbers1 := []int{1, 2, 3, 4, 5, 6}
	for i := range numbers1 {
		if i == 3 {
			numbers1[i] |= i
		}
	}
	fmt.Println(numbers1)

	// 动态设置数组的长度
	numbers2 := [...]int{1, 2, 3, 4, 5, 6}
	maxIndex2 := len(numbers2) - 1
	for i, e := range numbers2 {
		if i == maxIndex2 {
			numbers2[0] += e
		} else {
			numbers2[i+1] += e
		}
	}
	fmt.Println(numbers2)

	number3 := []int{1, 2, 3, 4, 5, 6}
	maxIndex3 := len(number3) - 1
	for i, e := range number3 {
		if i == maxIndex3 {
			number3[0] += e
		} else {
			number3[i+1] += e
		}
	}
	fmt.Println(number3)

	// value1 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	// switch 1 + 3 {
	// case value1[0], value1[1]: // 报错，类型无法匹配 int8 int
	// 	fmt.Println("0 or 1")
	// case value1[2], value1[3]:
	// 	fmt.Println("2 or 3")
	// case value1[4], value1[5], value1[6]:
	// 	fmt.Println("4 or 5 or 6")
	// }

	value2 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	switch value2[4] {
	case 0, 1: // 这里是可以通过编译的，会自动转换switch表达式的结果类型
		fmt.Println("0 or 1")
	case 2, 3:
		fmt.Println("2 or 3")
	case 4, 5, 6:
		fmt.Println("4 or 5 or 6")
	}

	// switch fallthrought
	name := "marsonshine"
	switch name {
	case "marsonshine":
		fmt.Println("marsonshine case")
		fallthrough // 表示继续执行下一个case
	case "summerzhu":
		fmt.Println("summerzhu case")
	default:
		fmt.Println("default case")
	}

	// switch type
	value6 := interface{}(byte(127))
	switch t := value6.(type) {
	case uint8, uint16:
		fmt.Println("uint8 or uint16")
	case byte: // 为什么无法通过编译？因为byte和uint8是一样的类型，byte就是uint8的别名
		fmt.Println("byte")
	default:
		fmt.Printf("unsupported type: %T", t)
	}
}

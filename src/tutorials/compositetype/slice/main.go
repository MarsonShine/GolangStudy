package main

import "fmt"

// slice 可以动态的增长或缩小,底层公用一个数组对象,给方法传参可以避免像数组一样会重新分配内存和数据拷贝到新的数组对象
// 对数组进行 slice 操作之后的结果就是 slice 对象。即：array[:]
// 也可以通过 make([]T, len),make([]T, len, cap), []{}
// slice 有三个属性，length: len(slice), cap: cap(slice), 指向数据对象的指针，只要索引不超过 cap 就能动态的扩容
func main() {
	a := [...]int{0, 1, 2, 3, 4, 5}
	reverse(a[:])
	fmt.Println(a) // "[5 4 3 2 1 0]"

	s := []int{0, 1, 2, 3, 4, 5}
	// Rotate s left by two positions.
	reverse(s[:2])
	reverse(s[2:])
	reverse(s)
	fmt.Println(s) // "[2 3 4 5 0 1]"

	array := [4]int{1, 2, 3, 4}
	printArray(array)
	// printSlice(array) // 报错，因为 array 不是 slice 数据类型
	printSlice(array[:])

	data := []string{"one", "", "three"}
	fmt.Printf("%q\n", nonempty(data)) // `["one" "three"]`
	// 由于共享底层数组对象，数组中的数据被覆盖了
	fmt.Printf("%q\n", data) // `["one" "three" "three"]`

	data2 := []string{"one", "", "three"}
	fmt.Printf("%q\n", nonempty2(data2)) // `["one" "three"]`
	fmt.Printf("%q\n", data2)            // `["one" "three" "three"]`

}

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// slice 基于每个项进行值比较
func equal(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func printArray(array [4]int) {
	fmt.Printf("形参array地址:%p\n", &array)
	for _, e := range array {
		fmt.Print(e)
	}
	fmt.Println()
}

func printSlice(slice []int) {
	fmt.Printf("形参slice地址:%p\n", &slice)
	for _, e := range slice {
		fmt.Print(e)
	}
	fmt.Println()
}

func appendInt(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1
	if zlen <= cap(x) {
		// 小于 cap 则可以扩容
		z = x[:zlen]
	} else {
		// 如果超过 cap，则以 2*len(x) 的容量扩容
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zcap)
		copy(z, x)
	}
	z[len(x)] = y
	return z
}

// 返回不为空的对象数组，但是覆盖了原来的项
func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}

func nonempty2(strings []string) []string {
	out := strings[:0] // zero-length slice of original
	for _, s := range strings {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

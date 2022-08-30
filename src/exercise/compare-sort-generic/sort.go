package main

// https://github.com/eliben/code-for-blog/tree/master/2022/genericsort
// https://eli.thegreenplace.net/2022/faster-sorting-with-go-generics/
// go build -o bubble.out
// .\bubble.out -cpuprofile cpui.out -kind strinterface  (must be linux)
//  go tool pprof -list bubbleSortInterface ./bubble.out cpui.out
// go tool objdump -S .\bubble.out >> dump.txt
// go tool compile -S .\sort.go >> sort_dump.txt
import (
	"math/rand"
	"sort"
	"strings"
)

type (
	/*
		~表达式的意思就是说，我们通常不关心特定的类型，比如string;我们对所有字符串类型都感兴趣。这就是 ~ 的作用。
		表达式 ~string 表示底层类型为 string 的所有类型的集合。这包括类型字符串本身以及所有类型声明的定义，如type MyString string
	*/
	Ordered interface {
		Integer | Float | ~string
	}
	Integer interface {
		Signed | Unsigned
	}
	Signed interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64
	}
	Unsigned interface {
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
	}
	Float interface {
		~float32 | ~float64
	}
)

/*
泛型版本的排序与泛型之前的排序用的算法是一样的，但是性能泛型版本要比非泛型要高一点
*/
func bubbleSortInterface(x sort.Interface) {
	n := x.Len()
	for {
		swapped := false
		for i := 1; i < n; i++ {
			if x.Less(i, i-1) {
				x.Swap(i, i-1)
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

func bubbleSortGeneric[T Ordered](x []T) {
	n := len(x)
	for {
		swapped := false
		for i := 1; i < n; i++ {
			if x[i] < x[i-1] {
				x[i-1], x[i] = x[i], x[i-1]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

func bubbleSortFunc[T any](x []T, less func(a, b T) bool) {
	n := len(x)
	for {
		swapped := false
		for i := 1; i < n; i++ {
			if less(x[i], x[i-1]) {
				x[i-1], x[i] = x[i], x[i-1]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

func makeRandomStrings(n int) []string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	ss := make([]string, n)
	for i := 0; i < n; i++ {
		// Each random string has length between 2-11
		var sb strings.Builder
		slen := 2 + rand.Intn(10)
		for j := 0; j < slen; j++ {
			sb.WriteRune(letters[rand.Intn(len(letters))])
		}
		ss[i] = sb.String()
	}
	return ss
}

type myStruct struct {
	a, b, c, d string
	n          int
}

type myStructs []*myStruct

func (s myStructs) Len() int           { return len(s) }
func (s myStructs) Less(i, j int) bool { return s[i].n < s[j].n }
func (s myStructs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func makeRandomStructs(n int) myStructs {
	structs := make([]*myStruct, n)
	for i := 0; i < n; i++ {
		structs[i] = &myStruct{n: rand.Intn(n)}
	}
	return structs
}

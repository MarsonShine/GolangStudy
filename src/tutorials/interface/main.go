package main

import (
	"bytes"
	"io"
	"reflect"
)

//因为在Go语言中只有当两个或更多的类型实现一个接口时才使用接口，它们必定会从任意特定的实现细节中抽象出来。结果就是有更少和更简单方法的更小的接口
//（经常和io.Writer或 fmt.Stringer一样只有一个）。当新的类型出现时，小的接口更容易满足。对于接口设计的一个好的标准就是 ask only for what you need（只考虑你需要的东西）

const debug = false

func main() {
	var buf *bytes.Buffer
	if debug {
		buf = new(bytes.Buffer) // enable collection of output
	}
	f(buf) // NOTE: subtly incorrect!
	if debug {
		// ...use buf...
	}
}

func f(out io.Writer) {
	// ...do something...
	if out != nil {
		out.Write([]byte("done!\n"))
	}
}

// 因为接口是隐式实现的，所以具体类型判断需要针对每个具体类进行转换，看是否转换成功
func formatOneValue(x interface{}) string {
	if err, ok := x.(error); ok {
		return err.Error()
	}
	if str, ok := x.(string); ok {
		return str
	}
	// ...其它类型
	return string("")
}

// 接口的完整性检查，这是最佳实践
type Shape interface {
	Sides() int
	Area() int
}
type Square struct {
	len int
}

func (s *Square) Sides() int {
	return 4
}
func (s *Square) Area() int {
	return 4
}

// 这里没有对 Shape 接口的所有方法进行实现，所以我们要通过如下方式对接口实现进行完整性检查
var _ Shape = (*Square)(nil)

type Munger interface {
	Munge(int)
}
type Foo struct{}

// 接口与类型之间的转换
func typeConvert() {
	var f Foo
	// _, ok := f.(Munger) // 转换错误，f必须要是接口类型

	// 可以先转换成空接口，然后在转成目标接口类型
	_, ok := interface{}(f).(Munger)

	// 也可以通过反射进行类型断言
	iMunger := reflect.TypeOf((*Munger)(nil)).Elem()
	ok = reflect.TypeOf(&f).Implements(iMunger)
}

package main

import "fmt"

type TestStruct struct{}

func NilOrNot(v interface{}) bool {
	return v == nil
}

// Go 语言根据interface的不同具体分为两种处理方式
// 一种是空接口（没有任何方法），在 golang 的数据结构之间用 eface 结构体表示，这个结构体只有两个字段，一个类型指针，一个是底层数据，这是最基本的底层结构，所以所有对象都能转换成 interface{} 类型
// 另一种是包含方法的接口，在 golang 中用数据结构 iface 表示，这个结构体对象有两个字段，一个是类型指针，一个是 *itab 结构指针，这个对象包含了很多数据，其中包含了常用的size、hash、equal 等。其中 hash 是用来快速判断类型是否相等的，在类型转换时就是靠这个来转换正确的具体类型的
type Duck interface {
	Quack()
}

type Cat struct {
	Name string
}

//go:noinline
func (c Cat) Quack() {
	println(c.Name + " meow")
}

func main() {
	var c Duck = Cat{Name: "grooming"} // 这里汇编会直接生成对应的实现类 Cat 调用方法 Quack ，编译器会做优化，将接口方法动态派发的过程省略掉直接转换成目标类型调用来减少额外的性能开销
	c.Quack()
	// 将接口转具体类型
	switch c.(type) {
	case *Cat:
		cat := c.(*Cat)
		cat.Quack()
	}

	var s *TestStruct
	fmt.Println(s == nil) // #=> true
	// 当把 s 变量传给 NilOrNot 方法时，发生了隐式类型转换，把 nil 的指针类型转换成了 interface{} 类型
	// 转换后包含了转换前的变量，还包含了变量的类型信息即 TestStruct，所以直接判断 nil 是不等于 nil 的
	// 应该通过反射来判断 interface{} 是否为一个空对象
	fmt.Println(NilOrNot(s)) // #=> false

	// 这里说一下具体转换过程
	// 因为 interface 类型变量包含了两个指针，一个指针指向值的类型，另一个指针指向实际的值
	// 因为 NilOrNot 函数接收一个 interface{} 类型的参数，所以当我们把一个 nil 的结构体指针传给这个方法就会发生隐式类型转换成 interface{}
	// 这个时候 interface{} 会依旧对自身两个指针进行赋值，一个是指向值的类型，一个是实际值 所以与 nil 这个是不一样的。
}

// 接口的动态派发
// 动态派发在结构体上的表现非常差，这也提醒我们应当尽量避免使用结构体类型实现接口。

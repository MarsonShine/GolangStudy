package main

// 反射主要围绕两点进行的：类型和数据值
// reflect.TypeOf 获取反射对象类型
// reflect.ValueOf 获取反射对象数据值
// 知道了类型和具体内容就能明确知道具体的类型信息了
// 获取类型 Type 之后就能知道这个类型有什么具体方法了
// 通过 Method 方法获取类型的方法
// 通过 Field 方法获取字段

import (
	"fmt"
	"reflect"
)

func main() {
	author := "marsonshine"
	fmt.Println("TypeOf author:", reflect.TypeOf(author))
	fmt.Println("ValueOf author:", reflect.ValueOf(author))

	// 从反射对象转换到接口
	// 过程为：反射对象 -> 接口对象 -> 显示转换到原始类型
	// 从接口值到反射对象是上面的逆过程
	// 基本类型 -> 接口类型 -> 转换为反射对象
	v := reflect.ValueOf(1)
	_ = v.Interface().(int)

	// 在利用反射修改值是要注意，要用指针修改对象值
	// 因为 golang 是按值传递的，当把反射对象传到方法中时会复制整个对象
	i := 1
	v = reflect.ValueOf(i)
	v.SetInt(10)
	fmt.Println(i) // 这样会报错，传递i时会重新创建一个新的对象

	v2 := reflect.ValueOf(&i)
	v2.SetInt(10)
	fmt.Println(i)

	// reflect.Value.Elem() 获取指针指向的变量
	v2.Elem().SetInt(20)

	// 判断一个类型是否实现了某个接口
	typeOfError := reflect.TypeOf((*error)(nil)).Elem()
	customErrorPtr := reflect.TypeOf(&CustomError{})
	customError := reflect.TypeOf(CustomError{})

	fmt.Println(customErrorPtr.Implements(typeOfError)) // true
	fmt.Println(customError.Implements(typeOfError))    // false

	// 利用反射调用某个函数、方法
	v3 := reflect.ValueOf(Add)
	if v3.Kind() != reflect.Func {
		return
	}
	t := v3.Type()
	argv := make([]reflect.Value, t.NumIn())
	for i := range argv {
		if t.In(i).Kind() != reflect.Int {
			return
		}
		argv[i] = reflect.ValueOf(i)
	}
	result := v3.Call(argv)
	if len(result) != 1 || result[0].Kind() != reflect.Int {
		return
	}
	fmt.Println(result[0].Int())
}

type CustomError struct{}

func (*CustomError) Error() string {
	return ""
}

func Add(a, b int) int {
	return a + b
}

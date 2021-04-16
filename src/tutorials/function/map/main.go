package main

import (
	"fmt"
	"reflect"
	"strings"
)

// https://coolshell.cn/articles/21164.html -- Map、Reduce、Filter 代码模式
// 泛型 Map
func Map(data interface{}, fn interface{}) interface{} {
	vfn := reflect.ValueOf(fn)
	vdata := reflect.ValueOf(data)
	result := make([]interface{}, vdata.Len())
	for i := 0; i < vdata.Len(); i++ {
		result[i] = vfn.Call([]reflect.Value{vdata.Index(i)})[0].Interface()
	}
	return result
}

// 上述简单的泛型版本是没有做类型检查的，在调用函数的时候参数不合规就会在运行时爆异常
// 类型检查
func SafeMap(slice, function interface{}) interface{} {
	return safeMap(slice, function, false)
}

func safeMap(slice, function interface{}, inPlace bool) interface{} {

	//检查slice类型是否为Slice
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("transform: not slice")
	}
	//检查函数签名
	fn := reflect.ValueOf(function)
	elemType := sliceInType.Type().Elem()
	if !verifyFuncSignature(fn, elemType, nil) {
		panic("trasform: function must be of type func(" + sliceInType.Type().Elem().String() + ") outputElemType")
	}
	sliceOutType := sliceInType
	if !inPlace {
		sliceOutType = reflect.MakeSlice(reflect.SliceOf(fn.Type().Out(0)), sliceInType.Len(), sliceInType.Len())
	}
	for i := 0; i < sliceInType.Len(); i++ {
		sliceOutType.Index(i).Set(fn.Call([]reflect.Value{sliceInType.Index(i)})[0])
	}
	return sliceOutType.Interface()
}

func verifyFuncSignature(fn reflect.Value, types ...reflect.Type) bool {
	// 检查是否为函数类型
	if fn.Kind() != reflect.Func {
		return false
	}
	// NumIn() - 返回函数输入参数的数量
	// NumOut() - 返回函数返回参数的数量
	if (fn.Type().NumIn() != len(types)-1) || (fn.Type().NumOut() != 1) {
		return false
	}
	// In() - 返回函数类型的输入参数的第i个参数类型
	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}
	// Out() - 返回函数类型的输出参数的第i个参数类型
	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}
	return true
}

func main() {
	square := func(x int) int {
		return x * x
	}
	nums := []int{1, 2, 3, 4}
	squared_arr := Map(nums, square)
	fmt.Println(squared_arr)

	upcase := func(s string) string {
		return strings.ToUpper(s)
	}
	strs := []string{"Marson", "Shine"}
	upstrs := Map(strs, upcase)
	fmt.Println(upstrs)
}

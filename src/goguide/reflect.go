package main

import (
	"fmt"
	"reflect"
)

func reflectExample() {
	var i int32 = 1
	itype := reflect.TypeOf(i)
	ivalue := reflect.ValueOf(i)
	fmt.Println(itype.Kind())

	ii := ivalue.Interface().(int32)
	fmt.Println(ii)

	type myInt32 int32
	var j myInt32 = 2
	jtype := reflect.TypeOf(j)
	fmt.Println(jtype.Kind())
}

func editReflectValue() {
	var i int32 = 1
	ivalue := reflect.ValueOf(i)
	if ivalue.CanSet() {
		ivalue.SetInt(2)
	}

	fmt.Println(ivalue)
}

func editAddressReflectValue() {
	var i int32 = 1
	ivalue := reflect.ValueOf(&i)
	fmt.Println(ivalue)

	fmt.Println(ivalue.CanSet())

	value := ivalue.Elem()
	fmt.Println(value.CanSet())
	value.SetInt(2)
	fmt.Println(value)
	fmt.Println(i)
}

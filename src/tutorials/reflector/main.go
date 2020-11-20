package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"
	"tutorials/reflector/display"
	"tutorials/reflector/method"
)

// Go 类型一般类型结构由两个结构组成：Type 和 Value，Type 表示 Go 的类型，Value 表示对象值
// reflect.TypeOf 返回的是具体类型
// %T 格式化符号输出对象类型
// reflect.Value 有很多方法
// NumFiled 获取结构体的成员数量
// MapKeys 获取每个map中的key
// Elem 获取指针
// 接口
// 通过反射获取的值即 reflector.ValueOf(inter) 有些是不允许更改的，因为是不可取地址的
// 例如：
// x := 2                   // value   type    variable?
// a := reflect.ValueOf(2)  // 2       int     no
// b := reflect.ValueOf(x)  // 2       int     no
// c := reflect.ValueOf(&x) // &x      *int    no
// d := c.Elem()            // 2       int     yes (x)

func main() {
	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w)) //"*os.File"

	fmt.Println(reflect.TypeOf(3)) // int

	v := reflect.ValueOf(3)
	fmt.Println(v)
	fmt.Printf("%v\n", v)   // "3"
	fmt.Println(v.String()) // NOTE: "<int Value>"

	// 可以直接对 Value 进行取 Type：v.Type()
	t := v.Type() //reflect.type
	fmt.Println(t.String())

	// 逆操作
	fmt.Println("反射逆操作")
	v = reflect.ValueOf(3)
	x := v.Interface()
	i := x.(int)
	fmt.Printf("%d\n", i)

	display_test()
	update_property_reflection()
	show_method()
}

func display_test() {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Color           bool
		Actor           map[string]string
		Oscars          []string
		Sequel          *string
	}
	//!-movie
	//!+strangelove
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},

		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	display.Display("strangelove", strangelove)
}

// 其中a对应的变量不可取地址。因为a中的值仅仅是整数2的拷贝副本。b中的值也同样不可取地址。c中的值还是不可取地址，它只是一个指针&x的拷贝。实际上，所有通过reflect.ValueOf(x)返回的reflect.Value都是不可取地址的。
func update_property_reflection() {
	x := 2
	a := reflect.ValueOf(x)
	b := reflect.ValueOf(2)
	c := reflect.ValueOf(&x)
	d := c.Elem() // 这个是可以修改的
	// 通过 CanAddr 判断能否获取地址
	fmt.Println(a.CanAddr()) // "false"
	fmt.Println(b.CanAddr()) // "false"
	fmt.Println(c.CanAddr()) // "false"
	fmt.Println(d.CanAddr()) // "true"

	px := d.Addr().Interface().(*int) // px:=&x
	*px = 3
	fmt.Println(x) // x 的值改变了

	// 如果不适用指针，则可以用 reflector.Value.Set
	d.Set(reflect.ValueOf(4))
	fmt.Println(x) // x 的值改变了

	// 如果类型不对，则会发生转换类型失败的panic错误
	defer func() {
		switch p := recover(); p {
		case nil: // no panic
		default:
			fmt.Println("数据类型转换错误") // unexpected panic; carry on panicking
		}
	}()
	d.Set(reflect.ValueOf(int64(5)))

	var y interface{}
	ry := reflect.ValueOf(&y).Elem()
	ry.SetInt(2)                     // 错误，无法在 interface Value 上调用 SetInt
	ry.Set(reflect.ValueOf(3))       // ok
	ry.SetString("hello")            // 错误，无法在 interface Value 上调用 SetString
	ry.Set(reflect.ValueOf("hello")) // ok
}

func show_method() {
	method.Print(time.Hour)
	method.Print(new(strings.Replacer))
}

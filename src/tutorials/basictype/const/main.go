package main

import "fmt"

// 批量申明，可以省略
const (
	a = 1
	b
	c = 2
	d
)

// iota 常量生成器
// 在一个const声明语句中，在第一个声明的常量所在的行，iota将会被置为0，然后在每一个有常量声明的行加 1。
type Weekday int

const (
	Sunday Weekday = iota //  相当于枚举
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Flags uint

const (
	FlagUp           Flags = 1 << iota // is up
	FlagBroadcast                      // supports broadcast access capability
	FlagLoopback                       // is a loopback interface
	FlagPointToPoint                   // belongs to a point-to-point link
	FlagMulticast                      // supports multicast access capability
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
	GiB // 1073741824
	TiB // 1099511627776             (exceeds 1 << 32)
	PiB // 1125899906842624
	EiB // 1152921504606846976
	ZiB // 1180591620717411303424    (exceeds 1 << 64)
	YiB // 1208925819614629174706176
)

func main() {
	fmt.Println(a, b, c, d) // "1 1 2 2"

	fmt.Println(Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday) // "0 1 2 3 4 5 6"
}

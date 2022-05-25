package main

// var a, b int
var a string
var done bool

func main() {
	go f()
	g()
}

func f() {
	a = 1 // 写之前的操作
	b = 2 // 写操作w
}

func g() {
	print(b) // 读操作
	print(a) // ？？？
}

func setup() {
	a = "hello, world"
	done = true
}

/*
指令重排，以及可见性的问题，多核CPU并发执行导致程序的运行和代码的书写顺序不一致
*/

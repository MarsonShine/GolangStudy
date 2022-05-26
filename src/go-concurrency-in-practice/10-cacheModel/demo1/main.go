package main

var a, b int

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

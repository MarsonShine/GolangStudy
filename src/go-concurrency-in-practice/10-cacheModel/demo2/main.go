package main

var a string
var done bool

// 下面的赋值指令可能会随着CPU架构的不同，执行的顺序的也不同
// 所以就会导致main中最后打印出来的a可能是空
func setup() {
	a = "hello, world"
	done = true
}

func main() {
	go setup()
	for !done {
	}
	print(a)
}

package main

type T struct {
	msg string
}

var g *T

// 下面的赋值指令可能会随着CPU架构的不同，执行的顺序的也不同
// 所以就会导致main中最后打印出来的msg可能是空
func setup() {
	t := new(T)
	t.msg = "hello, world"
	g = t
}

func main() {
	go setup()
	for g == nil {
	}
	print(g.msg)
}

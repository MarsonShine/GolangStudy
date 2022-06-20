package main

import (
	"fmt"
	"tutorials/delegate/undo"
)

type Widget struct {
	X, Y int
}

type Label struct {
	Widget // 嵌入
	Text   string
	X      int
}

func (l Label) Paint() {
	fmt.Printf("[%p] - Label.Paint(%q)\n", &l, l.Text)
}

func main() {
	lable := Label{Widget{10, 10}, "State", 100} // label.X = 100; lable.Widget.X = 10;
	fmt.Printf("X=%d, Y=%d, Text=%s Widget.X=%d\n", lable.X, lable.Y, lable.Text, lable.Widget.X)

	ints := undo.NewIntSet()
	for _, i := range []int{1, 3, 5, 7, 9} {
		ints.Add(i)
		fmt.Println(ints)
	}

	for _, i := range []int{1, 2, 3, 4, 5, 6, 7} {
		fmt.Print(i, ints.Contains(i), " ")
		ints.Delete(i)
		fmt.Println(ints)
	}
}

package main

import (
	"fmt"
	"time"
)

func foobyref(n *int) {
	fmt.Println(*n)
}

func main() {
	for i := 0; i < 5; i++ {
		go foobyref(&i)
	}
	time.Sleep(100 * time.Millisecond)
}

/*
byvalue执行的时候是正常结果；
byreference循环执行的都是最后一个值；
为什么会这样？

go语言规范中强调：由init语句申明的变量在每次迭代中都会被重复使用
这意思就是说，在程序运行时，只有一个代表 i 的对象，而不是每次迭代都会分配一个新对象。这个对象在每次迭代时都会分配一个新值。
*/

package main

import "sync"

/*
指令重排，以及可见性的问题，多核CPU并发执行导致程序的运行和代码的书写顺序不一致

重点：
在一个 goroutine 内部，程序的执行顺序和它们的代码指定的顺序是一样的
即使编译器或者 CPU 重排了读写顺序，从行为上来看，也和代码指定的顺序一样

解析：
在 goroutine 内部对一个局部变量 v 的读，一定能观察到最近一次对这个局部变量 v 的
写。如果要保证多个 goroutine 之间对一个共享变量的读写顺序，在 Go 语言中，可以使
用并发原语为读写操作建立 happens-before 关系，这样就可以保证顺序了。

channel 的机制可以保证执行顺序的一致性，规则如下：
1.send-revc机制，send一定发生在recv完成之前
2.close一个channel，一定发生在从关闭的channel中读取一个零值之前
3.对于unbuffer的channel，读取操作的调用一定发生在发送数据的调用之前
4.channel容量是m（m>0），那么第n个receive一定发生在第n+m个send的完成之前

tips：
在讨论并发编程时，当我们说x事件在y事件之前发生（happens before），我们并不是说x事件在时间上比y时间更早；我们要表达的意思是要保证在此之前的事件都已经完成了，例如在此之前的更新某些变量的操作已经完成，你可以放心依赖这些已完成的事件了。
当我们说x事件既不是在y事件之前发生也不是在y事件之后发生，我们就说x事件和y事件是并发的。这并不是意味着x事件和y事件就一定是同时发生的，我们只是不能确定这两个事件发生的先后顺序。
--《go语言圣经》
*/

// rule 1
var s string
var ch = make(chan struct{}, 10)

func f1() {
	s = "hello, world"
	ch <- struct{}{}
}

func main() {
	go f1()
	<-ch
	print(s)
}

// rule2
func f2() {
	s = "hello, world"
	close(ch)
}

func main2() {
	go f2()
	<-ch
	print(s)
}

// 上述s的初始化，发生在close之前，而根据规则2，close一定发生在读取一个零值之前

// rule3
var ch3 = make(chan struct{})

func f3() {
	s = "hello, world"
	<-ch
}
func main3() {
	go f3()
	ch <- struct{}{}
	print(s)
}

// 如果第 60 行发送语句执行成功（完毕）
// 那么根据这个规则，第 56 行（接收）的调用肯定发生了（执行完成不完成不重要，重要的是这一句“肯定执行了”），
// 那么 s 也肯定初始化了，所以一定会打印出“hello world”

// mutex/rwmutex 也有顺序保证
var mu sync.Mutex

func foo() {
	s = "hello, world"
	mu.Unlock()
}
func main5() {
	mu.Lock()
	go foo()
	mu.Lock()
	print(s)
}

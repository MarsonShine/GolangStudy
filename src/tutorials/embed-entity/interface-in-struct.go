package main

import "fmt"

type Fooer interface {
	Foo() string
}

type Container3 struct {
	Fooer
}

func (co Container3) Foo() string {
	return co.Fooer.Foo()
}

/*
这种可以实现策略模式，状态模式等来达到OCP
*/

func sink(f Fooer) {
	fmt.Println("sink:", f.Foo())
}

type FooerImp struct {
}

func (f FooerImp) Foo() string {
	return "FooerImp Foo"
}

func show2() {
	co := Container3{Fooer: FooerImp{}}
	sink(co)
}

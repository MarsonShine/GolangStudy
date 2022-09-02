package main

import (
	"fmt"
)

// 嵌入

type Base struct {
	b int
}

type Container struct { // Container 是嵌入结构体
	Base // Base是被嵌入的结构体
	c    string
}

func (base Base) Describe() string {
	return fmt.Sprintf("base %d belongs to us", base.b)
}

func (container Container) Describe() string {
	return container.Base.Describe()
}

func show() {
	container := &Container{
		Base: Base{
			b: 1,
		},
		c: "marsonshine",
	}
	container.c = "marsonshine"
	// 属性
	fmt.Printf("能直接访问被嵌入结构体的属性: container.b = %d; container.Base.b = %d\n", container.b, container.Base.b)
	fmt.Printf("container.c = '%s'", container.c)
	// 方法
	fmt.Printf(container.Describe())
}

// 当嵌入结构体的属性与被嵌入结构体的属性相同时，让访问对应的属性时，得到的就是嵌入结构体的属性
// 而此时的被嵌入结构体的属性就成为了阴影属性

type Base2 struct {
	b   int
	tag string
}

type Container2 struct {
	Base2
	c   string
	tag string
}

func (c Container2) DescribeTag() string {
	return fmt.Sprintf("Container tag is %s", c.tag)
}

func (b2 Base2) DescribeTag() string {
	return fmt.Sprintf("Base2 tag is %s", b2.tag)
}

func showShadows() {
	co := Container2{
		Base2: Base2{
			b:   1,
			tag: "base",
		},
		c:   "marsonshine",
		tag: "container2",
	}
	fmt.Println(co.Base2.DescribeTag())
	fmt.Println(co.DescribeTag())
}

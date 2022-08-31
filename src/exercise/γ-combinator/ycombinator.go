package main

import "fmt"

type Func[T, U any] func(T) U
type TagFunc[T, U any] func(Func[T, U]) Func[T, U]
type CombinatorFunc[T, U any] func(CombinatorFunc[T, U]) Func[T, U]

/*
Y组合子的定义使用了几种匿名函数。在Go里面我们必须要有显式类型。Func是底层计算函数，如果我们使用普通递归，它就是函数类型。它使用两种泛型类型进行参数化：T代表参数类型，U代表返回类型
TagFunc是一个函数类型，CombinatorFunc只在Y组合符本身的定义中使用
*/

func Y[T, U any](f TagFunc[T, U]) Func[T, U] {
	return func(self CombinatorFunc[T, U]) Func[T, U] {
		return f(func(n T) U {
			return self(self)(n)
		})
	}(func(self CombinatorFunc[T, U]) Func[T, U] {
		return f(func(n T) U {
			return self(self)(n)
		})
	})
}

var factorial_tag = func(recurse Func[int, int]) Func[int, int] {
	return func(n int) int {
		if n == 0 {
			return 1
		}
		return n * recurse(n-1)
	}
}

var fib_tag = func(recurse Func[int, int]) Func[int, int] {
	return func(n int) int {
		if n <= 1 {
			return n
		}
		return recurse(n-1) + recurse(n-2)
	}
}

type Node struct {
	val   int
	left  *Node
	right *Node
}

var treesum_tag = func(recurse Func[*Node, int]) Func[*Node, int] {
	return func(n *Node) int {
		if n == nil {
			return 0
		} else {
			return n.val + recurse(n.left) + recurse(n.right)
		}
	}
}

func main() {
	fac := Y(factorial_tag)
	fmt.Printf("fac(%d)=%d", 10, fac(10))

	treesum := Y(treesum_tag)
	tree := &Node{
		val: 1,
		left: &Node{
			val: 2,
			left: &Node{
				val: 3,
				left: &Node{
					val: 4,
				},
				right: &Node{
					val: 5,
				},
			},
			right: &Node{
				val: 6,
			},
		},
		right: &Node{
			val: 7,
		},
	}
	fmt.Printf("treesum(node)=%d", treesum(tree))

	fibsum := Y(fib_tag)
	fmt.Printf("fib(%d)=%d", 100, fibsum(10))
}

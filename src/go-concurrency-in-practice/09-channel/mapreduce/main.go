package main

import "fmt"

// map-reduce 分为两个步骤，第一步是映射（map），处理队列中的数据；
// 第二步是规约（reduce），把列表中的每一个元素按照一定的处理方式处理成结果，放入到结果队列中
// 与JavaScript中的array.reduce，array.map类似
func main() {
	in := asStream(nil)
	// map
	mapFn := func(v interface{}) interface{} {
		return v.(int) * 10
	}
	// reduce
	reduceFn := func(c, v interface{}) interface{} {
		return v.(int) + c.(int)
	}

	sum := reduce(mapChan(in, mapFn), reduceFn) // 返回累加结果
	fmt.Println(sum)
}

func mapChan(in <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{}) // 创建一个输出chan
	if in == nil {
		close(out)
		return out
	}

	go func() {
		defer close(out)
		for v := range in {
			out <- fn(v)
		}
	}()

	return out
}

func reduce(in <-chan interface{}, fn func(r, v interface{}) interface{}) interface{} {
	if in == nil {
		return nil
	}

	out := <-in // 先读取第一个元素
	for v := range in {
		out = fn(out, v)
	}
	return out
}

// 生成一个数据流
func asStream(done <-chan struct{}) <-chan interface{} {
	s := make(chan interface{})
	values := []int{1, 2, 3, 4, 5}
	go func() {
		defer close(s)
		for _, v := range values { // 从数组生成
			select {
			case <-done:
				return
			case s <- v:
			}
		}
	}()
	return s
}

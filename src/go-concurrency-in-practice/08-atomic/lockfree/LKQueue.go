package lockfree

import (
	"sync/atomic"
	"unsafe"
)

type LKQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

type node struct {
	value interface{}
	next  unsafe.Pointer
}

func NewLKQueue() *LKQueue {
	n := unsafe.Pointer(&node{})
	return &LKQueue{
		head: n,
		tail: n,
	}
}

func (q *LKQueue) Enqueue(v interface{}) {
	n := &node{value: v}
	for {
		tail := load(&q.tail)
		next := load(&tail.next)
		if tail == load(&q.tail) {
			if next == nil { // 没有新元素入队列
				if cas(&tail.next, next, n) { // 增加到队尾
					cas(&q.tail, tail, n) // 入队成功，移动尾指针
				}
			} else {
				// 已有新数据加入到队列后面，需要移动尾指针
				cas(&q.tail, tail, next)
			}
		}
	}
}

// 将unsafe.Pointer原子加载转换成node
func load(p *unsafe.Pointer) (n *node) {
	return (*node)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *node) (ok bool) {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}

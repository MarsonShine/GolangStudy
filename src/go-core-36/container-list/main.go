package main

import (
	"container/list"
	"container/ring"
	"fmt"
)

// List 是双向链表
// 所以插入/删除第一个，最后一个都是O1级别的操作
func main() {
	list := list.New()
	list.Len() // O1
	// 遍历，慢查询操作On
	for {
		head := list.Front()
		if head.Next() != nil {
			fmt.Printf("cur = %v", head.Value)
			head = head.Next()
		}
	}

	{
		ring := ring.New(10)
		ring.Len() // On,内部采用遍历计算
	}

}

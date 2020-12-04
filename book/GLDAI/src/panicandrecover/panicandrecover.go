package main

import (
	"fmt"
)

func main() {
	defer println("in main")
	//////////////////////////////////////////////////////////////////////
	// go func() {
	// 	defer println("in goroutine")
	// 	panic("") // 抛错，当前协程域 defer 是有效的，所以只会输出 in goroutine
	// }()

	// time.Sleep(1 * time.Second)
	// // 这里也没有正常结束
	// // 因为recover只有在panic 发生之后才能出发，而这个是在发生异常之前就调用了
	// if err := recover(); err != nil {
	// 	fmt.Println(err)
	// }
	// 所以应该放在 defer 中延迟执行才能在 panic 之后触发恢复
	/////////////////////////////////////////////////////////////////////
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()
	// panic("unkonwn err")
	////////////////////////////////////////////////////////////////////
	multicall()
}

// panic 是可以嵌套多次调用
func multicall() {
	defer fmt.Println("in multicall")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			panic("panic again and again")
		}()
		panic("panic again")
	}()
	panic("panic once")
}

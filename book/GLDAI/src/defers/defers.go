package main

// defer 使用传值方式传递参数会进行预计算；#2
// defer 的执行顺序是在当前域方法执行完毕之后执行的 #1
// 编译器生成中间代码
// 使用三种不同的机制处理该关键字：开放编码、堆、栈
// 堆，性能最低
// 栈，优化内存布局，可以提高30%性能
// 开放编码，使用代码内联来优化，有使用条件：函数的 defer 数量少于 8 个；defer 不能循环中；函数的 return 语句与 defer 语句的乘积不小于 15 个
import (
	"fmt"
	"time"
)

func main() {
	// #1
	// for i := 0; i < 5; i++ {
	// 	defer fmt.Println(i) // 倒序输出
	// }
	// {
	// 	defer fmt.Println("defer runs")
	// 	fmt.Println("block ends")
	// }

	// fmt.Println("main ends") // 输出 block ends - main ends = defer runs

	// 计算程序运行的时间 #2
	function()
}

func function() {
	startedAt := time.Now()
	defer fmt.Println(time.Since(startedAt))
	time.Sleep(time.Second) // 输出 0s 因为预计算原因，当 defer 后的函数虽然会在函数结束后运行，但是引用的参数会拷贝，所以 time.Since(startedAt) 的结果不是在 main 函数结束时调用的，而是在 defer 关键调用的时候就”预计算“了，导致输出的结果为 0
	// 修改如下，使用匿名函数
	defer func() {
		fmt.Println(time.Since(startedAt))
	}()
	time.Sleep(time.Second)
}

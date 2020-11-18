package main

import (
	"flag"
	"fmt"
	"time"
)

// 实现 flag.Value 接口，就能自定义命令行标记，可以定义新的符号
var period = flag.Duration("period", 1*time.Second, "sleep period") // 运行  go run main.go -period 50ms/2m30s

func main() {
	flag.Parse()
	fmt.Printf("Sleeping for %v...", *period)
	time.Sleep(*period)
	fmt.Println()

	// 其实在 flag 库中已经构建好了，只需要实现 flag.Value 就可以达到上面的效果，详见 tempconv/main.go
}

package main

// go test -timeout 300s -run ^TestFail$ gocore36/test/demo2 -count=1
// 如果要选择打印日志 go test -v -timeout 300s -run ^TestFail$ gocore36/test/demo2 -count=1

// 基准测试，又称性能测试
// go test -benchmem -run=^$ -bench ^BenchmarkGetPrimes$ gocore36/test/demo3 -count=1
// go test -benchmem -run=^$ -bench ^BenchmarkGetPrimesWith100$ gocore36/test/demo4
// 可选参数：
// -benchmem 输出基准测试的内存分配统计信息。
// -benchtime 用于指定基准测试的探索式测试执行时间上限
func main() {

}

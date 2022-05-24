package main

// 流/管道模式
// 提供跳过几个元素，或者只取其中几个元素的功能
func main() {

}

// 创建流
func asStream(done <-chan struct{}, values ...interface{}) <-chan interface{} {
	s := make(chan interface{}) // 场景一个未缓冲的chan
	go func() {
		defer close(s)
		for _, v := range values {
			select {
			case <-done:
				return
			case s <- v: // 将数组元素塞入到chan中
			}
		}
	}()
	return s
}

// 获取前n个数据
// 过滤指定条件的数据
// 只取满足条件的数据，一旦不满足则停止
// 跳过前n个数据
// 跳过满足条件的数据
// 跳过满足条件的数据，一旦不满足，当前这个元素及之后的元素都会输出给 Channel 的 receiver

func takeN(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{}) // 输出流
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream: // 从输入流中读取元素发送给输出流（如同linux的管道（如 grep)将流的输出结果做另一个流的输入）
			}
		}
	}()
	return takeStream
}

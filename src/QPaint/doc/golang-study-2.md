# defer 关键字

首先来看官网的定义：

> A "defer" statement invokes a function whose execution is deferred to the moment the surrounding function returns, either because the surrounding function executed a [return statement](https://golang.org/ref/spec#Return_statements), reached the end of its [function body](https://golang.org/ref/spec#Function_declarations), or because the corresponding goroutine is [panicking](https://golang.org/ref/spec#Handling_panics).

就是被标记了 defer 的片段会调用一个函数，这个函数会推迟在周围函数执行完返回后执行。要注意，在最后的说明中还带有 panicking，这是什么呢？

看了官网文档对 panicking 的解释，我认为是就是一个**运行期的未知的**异常处理程序，比如数组越界就会触发一个 run-time panic。就相当于调用内置函数 panic，并用实现了接口类型 runtime-Error 的值做为参数。触发这个错误就代表着它是未确定的错误。

defer 调用的格式为

```
DeferStmt = "defer" Expression .	// Expression 必须是方法或函数
```

这里面有个很重要点：

```
Each time a "defer" statement executes, the function value and parameters to the call are evaluated as usual and saved anew but the actual function is not invoked. Instead, deferred functions are invoked immediately before the surrounding function returns, in the reverse order they were deferred.
```

每次一个 defer 语句调用，这个普通函数值和参数会被(重点)**重新保存**，但是实际上并没有调用。而是在环绕函数返回之前立即调用，（重点来了）并且会将标记的 defer 的执行的**顺序反转再调用**。

举个例子：

```go
// out func
func f() (result int) {
	defer func() {
		result *= 7
	}()
	return 6
}
for i := 0; i <= 3; i++ {
  defer fmt.Println(i)
}
fmt.Println("==================")
f()
```

通过之前我们讲的概念和重点，我们大概也知道输出的是什么了。结果我就不放了，大家看到了这篇文章就自己动手去试试吧。
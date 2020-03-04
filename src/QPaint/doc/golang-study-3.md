# golang 自学系列（三）—— if，for，channel

一般情况下，if 语句跟大多数语言的 if 判断语句一样，根据一个 boolean 表达式结果来执行两个分支逻辑。

但凡总是有例外，go 语言还有这种写法：

```
// 写法1
if i:= getID(); i < currentID {
	execute some logic
} else {
	execute some other logic
}

// 写法2
var obj = map[string]interface{}
if val,ok := obj["id"]; ok {
	execute some logic
} else {
	execute some other logic
}
```

写法 1 的意思是在判断逻辑前，可以加一个表达式，比如获取 ID 赋值给 i，然后参与后续的判断是否小于当前 ID。

写法 2 的意思同样是在判断逻辑前，可以加一个表达式，获取对象 ID（obj["id"]）给 val，但是与 1 不同的是，这里 val、ok 的值是有直接关联的。val 值取得成功与否，就是 ok 的结果值。

**即 ok 一定是 boolean 类型值，表示 val = obj["id"] 是否赋值成功**。我认为这种特性很好，完全不用取得值是否不存在，会报错等。

# for 语句

for 语句一般表示为重复执行块。由三个部分组成，一个是单条件控制的迭代，一个是 “for” 语句，最后一个是 “range” 语句。

```
ForStmt = "for" [ Condition | ForClause | RangeClause ] Block .
Condition = Expression .
```

for 查了资料才发现用法特别多

一. 使用单条件的 for 语句

```go
for a < b {
	a *= 2
}
```

这个是最简单的，意思就是只要你的条件计算得出来的是 true 就会重复执行代码段。就如上面所示，只要 a < b 就会一直会执行 a *= 2。相当于 while 死循环。

二. 使用 for 从句

```
// 句法格式
ForClause = [ InitStmt ] ";" [ Condition ] ";" [ PostStmt ] .
InitStmt = SimpleStmt .
PostStmt = SimpleStmt .

for i := 0; i < 10; i++ {
	f(i)
}
```

使用 for 从句的 for 语句也是通过条件控制的，但是会额外指定一个 init 以及 post 语句，就好比分配一个数，这个数会自增或递减。满足条件判断，就会在重复运行执行体。

只要初始化语句变量不为空，就会在第一次迭代运行之前计算。有几点要注意：

for 从句中的任何元素都可以为空，除非它之后一个条件，否则这种情况下分号是不能丢的。如果条件是缺省的，它就等价于这个条件是 true。例如

```
for condtion { exec() }  等同于 for ; condition ; { exec() }
for 				 { exec() }  等同于 for   true			 ; { exec() }
```

三. 使用 range 从句

使用了 range 从句的 for 语句代表从执行的这些对象，这些对象会是数组、分片、字符串或映射以及通道（channel）上接收的值。如果迭代的条目存在，就把它赋值给迭代变量。

```
RangeClause = [ ExpressionList "=" | IdentifierList ":=" ] "range" Expression .
```

这个表达式 “range” 的后边的表达式被称为 range 表达式，它可能是数组、数组指针、分片、字符串、映射（map）或者是通道接收操作（channel permitting [receive operations](https://golang.org/ref/spec#Receive_operator).）。就像赋值一样，如果左边有操作数，那么则必须是可寻址的或是映射索引表达式。

**如果范围表达式是一个通道（channel），那么最多只有一个迭代变量，否则最多有两个变量。**如果最后一个迭代变量是空标识符，那么就相当于没有这个标识符的 range 表达式。

range 表达式 x 要在循环体开始之前计算一次，有一个例外：如果存在最多一个迭代变量以及 len(x) 是常熟，那么 range 表达式就不会计算。下面是官网给出的例子

```go
var testdata *struct {
	a *[7]int
}
for i, _ := range testdata.a {
	// testdata.a is never evaluated; len(testdata.a) is constant
	// i ranges from 0 to 6
	f(i)
}

var a [10]string
for i, s := range a {
	// type of i is int
	// type of s is string
	// s == a[i]
	g(i, s)
}

var key string
var val interface{}  // element type of m is assignable to val
m := map[string]int{"mon":0, "tue":1, "wed":2, "thu":3, "fri":4, "sat":5, "sun":6}
for key, val = range m {
	h(key, val)
}

// key == last map key encountered in iteration
// val == map[key]

var ch chan Work = producer()
for w := range ch {
	doWork(w)
}

// empty a channel
for range ch {}
```

# Channel 类型

`chan` 关键字：起初这个官网没有查到具体对 chan 的解释，就只知道是一个关键字，后来经过一番资料查询，发现这个是用来方便创建 `channel` 类型的快捷方式。其默认值是 nil。

既然知道 chan 是用来创建 `channel` 的，那么我们就来看 `channel` 类型的定义：

```
A channel provides a mechanism for concurrently executing functions to communicate by sending and receiving values of a specified element type. The value of an uninitialized channel is nil.
```

就是说 channel 类型为并发运行函数提供一个机制，通过发送和接收指定元素类型的值通信。未分配的 channel 的值是 nil。

在来了解下 `chan` 的操作符 <-，它表示指定 channel 方向、发送或接收。如果没有给定方向，就是双向。channel 被限制只能通过赋值或显示转换来发送或接收。

```go
chan T 	// 可以用来发送和接收 T 类型的值
chan <- float64	// 只能被用于发送 64 位浮点数
<- chan int	// 只能用于接收 int 数
```

官网还描述下面这种用法，说实话我没怎么看懂

```go
chan<- chan int    // same as chan<- (chan int)
chan<- <-chan int  // same as chan<- (<-chan int)
<-chan <-chan int  // same as <-chan (<-chan int)
chan (<-chan int)
```

go 语言提供内置函数 `make` 来构建新的 channel 值，该函数传递 channel 类型参数和一个可选的容量（capacity）参数值

```go
make(chan int, 100)
```

其中的 100 是容量值（capacity），这个容量值是指 channel 内的缓冲大小。如果值为 0 说明没有缓冲，这种情况下只有当接受者和发送者准备好才能成功通讯。否则，只要这个缓冲块没满（推送）或不为空（接收），那已经缓冲的 channel 就能成功通信。

channel 能通过调用方法 `close` 关闭。多值赋值通过接收操作符的形式报告这个接收的值是否在 channel 在关闭之前。

单个 channel 能用在发送语句、接收操作符以及调用内置 `cap` 函数和 `len` 函数，也能在任意数量的 goroutline 使用，而不需要同步。通道是一个先进先出（FIFO）的队列。

上面提到了两点，发送语句和接收操作符

# send 语句

简单来讲就是发送一个值给 channel。

`ch <- 3` 意思是发送一个值 3 给 channel 变量 ch。如果 channel 关闭了，会报 run-time panic 错误。

# 接收操作符

对于 channel 类型的操作数 ch，接收操作符的值 <- ch 意思是从 channel 类型值 ch 接收的值。<- 右边是 channel 类型元素。表达式块只有在值可用才不会阻塞。所以空 channel 无法接受值，因为永远阻塞。

现在来看一下 chan 的一些例子

```go
var c chan int	// nil
c = make(chan int)	// 初始化
fmt.Printf("c 的类型是%T \n", cc)	// chan int
fmt.Printf("c 的值是%v \n", cc)	//	0xc0000820c0
```

这能看出 chan int 的值是一个地址，像是指针一样。不过目前我还是不知道这个具体的用法是什么，用在什么地方？

而当我尝试读取值的时候，却发现好像一直在阻塞：

```go
c <- 3
fmt.Printf("c 的值是%v \n", c)
<-c
fmt.Printf("c 的值是%v \n", c)
```

我想的是 chan 在发送时，没有接收前是堵塞的，所以一直没有执行下面的输出。所以我又在后面加了 `<- c` 让 c 接收。结果发现还是不执行上面的输出。

后来又查了相关的资料，得知 channel 类型很像是一个通道，消费者-生产者之间的关系。

于是我又写了下面代码

```go
cc := make(chan int)
defer close(cc)
cc <- 3 + 4
i := <-cc
fmt.Println(i)
```

断电调试发现（如何在 vscode 断电调试我稍微会另起文章说）程序运行到 cc <- 3 + 4 就不往下执行，进程也没结束，说明是阻塞的。从之前的概念上将，cc 初始化出来的类型是没有设置初始容量，即没有缓存，难道在不能发送数据了？为了验证想法，我在初始化 channel 的时候家了初始缓存 buffer：`cc := make(chan int, 100)`；能正常输出结果值。但是概念上并没有说没有缓冲区的就不能正常发送数据啊。

在描述 send 语句的时候说过，在接收器准备好了，发送器才会被处理。按照这个结果来看，接收器是没有准备好的。那要怎么才能使接收器准备好呢？

我又查了资料，发现基本上都是这么种写法，以上面的代码为前提

```go
cc := make(chan int)
defer close(cc)
go func() {
  cc <- 3 + 4
}()
i := <-cc
fmt.Println(i)
```

这样就是正常的。难道要把发送器以一种函数调用的形式存在？

我又尝试把立即执行函数改为普通函数

```go
cc := make(chan int)
defer close(cc)
fchan(cc)
i := <-cc
fmt.Println(i)
```

结果发现还是不行，这个问题先放一边吧。等以后有时间在仔细学习一下。

// TODO Range 也是可以处理 chan 

# 初始化数组

`[]type{}` 当我看到这个语句的时候，我内心是奔溃的，因为突然看到这种写法的我不知道这属于哪个特征，都不好查关键字资料。刚开始我以为是 type 类型的数组，然后后面再接一个空对象 `{}`。但是翻遍了特性，都没看到这种概念。后来直接写了一个例子，查看具体的输出

```go
var sa = []string{}
fmt.Printf("sa 的值是%v \n", sa)	//sa 的值是[] 
```

发现它就是一个数组，并没有什么对象。后来我尝试去掉后面 `{}` 发现根本通不过编译，这才知道，这就是初始化一个 string 数组。
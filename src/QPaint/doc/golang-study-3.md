# golang 自学系列（三）—— if，for 语句

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

```
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

# 待弄清的语句

`chan` 关键字:	// TODO

`[]ShapeID{}` 表达式:  // TODO
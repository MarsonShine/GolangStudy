# Golang自学系列

我在使用vscode进行go编程时，总会显示一下警告

```go
type Service struct {
		a *ClassName
}
exported type Service should have comment or be unexported
```

这是因为你安装插件 `gopls` 对代码会有一些规则。即如果你定一个 变量/类 时，如果开头是大写字母，那么编辑器就会检测出你没有针对这个 变量/类 进行注释。这个时候有两种选择

1. 将大写字母改成小写字母 service
2. 进行特定格式注释，如上面注释就是这样`//Service 服务类`

其实除了这个还有一种做法，直接设置vscod的相关检测的属性`"go.lintFlags":["--disable=all"]`，这样就不用写那些烦人的注释啦。

golang 的编码习惯有个很有意思，就是它的所有变量、方法、类 等等在代码的末句全都不需要打分号，就算你打了分号，编辑器一样也会给你自动省略

# 怎么导入第三方库？

直接在官方仓库地址选择自己要导入的模块，地址见：https://pkg.go.dev

```go
import ("moduleName")
```

但是我当初这么做的时候，发现虽然编译器没有检测到错误，但是在引用 module 的 api 时，没有智能提示。后来发现跟 nuget 是一样的，有快捷键导入 module: `shift + command + p` 然后选择 `Go: Add Import` 选中你本地 clone 下来的第三方类库。但是这里有一个疑问，难道要每次下载源代码吗？而是不能够直接下载一个类似 dll 的可执行的 “微文件” 么。讲道理是肯定有的，不可能一个项目发布出去，还把人家的源代码发不出去的。这个之后弄到发布的时候在回过头来查资料吧

错误1: `expected ';', found f` 这个是不同编辑器的编码问题，拿我现在用的 vscode 为例，我新建的 go 项目的默认编码是 `LF/CRLF`，切换成 `CRLF/LF` 。注意，如果切换的时候发现还是报同样的错误，那极有可能是 vscode 没有反应过来，只要在当前页面随便输个空格在保存即可。



接下来就是各种变量函数的基本用法介绍了

null 值：golong 用 nil 代表 null，这个很特别啊，大多数语言都是 null

定义类：用 type 关键字 `type ClassName struct{ someField int}`

申明变量：`var scopeVar = "string"` 我试了一下，这种显示写法也行 `var scopeVar string = ""`。

但是有这么一种写法 `localVar := ""`，也很有趣，我尝试了一下，这个好像只能在方法里面写（就相当于 `var localVar string = ""`），这种写法提升到 “全局” 则不行。

# golong 具有指针概念

指针保存的是变量的内存地址。比如 `var p *int` p 表示的是整形的地址，其零值是 nil

这个跟 c++ 的指针是一样的，比较复杂，当时上大学的我上课上到这个地方的时候很懵，什么“指针”，“指针的引用”，“指针的指针”，“ * ”，“ & ” 等总是搞不清楚。在用 vc++ 6.0 时代下，编写代码没有任何提示，简直是难如登天。

但是现在时代不同啦，ide/编辑器 可以自动帮你做正确的选择。这次偶然的机会学习 golang，不过我还是有必要把这地方的知识弄清楚。

首先看下面代码，我把注释写在边上

```go
var p *int	// 变量 p 代表是整形的内存地址
i := 11	// 就一般的变量赋值
p = &i	// 给内存地址指针变量 p 赋值 11 的指针变量
*p = 1	// p 地址的值赋值为 1
```

第二行我就不解释了。

第一行代码就是定义一个指向整形的内存地址的变量 p。

第三行代码表示你要给一个整形的地址赋值，那么肯定不是直接赋值一个整数 i，而是这个变量 i 指向的整数的地址 &i。其实可以理解为 i 的一个引用。

第四行我要直接给 p 指针指向的地址具体的值，那就是我们之前说的 “指针的指针：*p = 1”。

# 函数申明

**函数在 golang 里面同 js 是一样 —— 一等公民。也就是你无论写在哪里，它都是可以在当前域是有效，可以引用的。**

函数申明分两种

1. 无返回值：`func SomeMethod() {}`
   1. 带参数：`func SomeMethod(a int) {}`
2. 有返回值：`func SomeMethodAndReturn() ReturnValue {}`

函数这里面有个好玩的约定：

- **函数名首字母是小写就是 private 私有方法**
- **函数名首字母是大写则是 public 共有方法**

我们还可以定一个函数类（函数类就相当于 C# 的委托，委托对于 CLR 而言就是一个含有这么一个函数的类，也可以当做 Java 中的内部类处理）。实例代码如下所示

```go
type delegateFunc func(string)	// 定一个委托
func serve(msg string){
  fmt.Printf(msg)
}
func main(){
  d := delegateFunc(serve)	// 把函数当作参数传递
  d("marson shine")
}
```

# 方法定义

方法的定义跟函数很像: `func (type类型参数) MethodName(parameters) ReturnValue {}`

先来看[官网对方法的定义](http://docscn.studygolang.com/ref/spec#Method_sets)：

```
一个 type 指定的类型可以关联方法集。一个接口类型的方法集是其接口。任何类型 T 的方法集，由它作为接收器接收所有方法。对应的指针类型 *T 的方法集是由接收器 *T 或者是 T 申明的方法集（那也就说，它包含了 T 的所有方法集）。更多的规则运用在包含那些匿名字段的结构（struct）上。任何类型都有空的方法集。在一个方法集中，每个方法必须有一个唯一的不为空的名称。
```

我们举个例子来说明：

```go
func (typeName ClassName) MethodName(parameter string) string {

}
```

这里我们定义了一个方法，指定的接收器就是 ClassName 类型，即我们得先有个接收器，才能有这个方法集 MethodName

```go
type ClassName struct {
		userName string
}
```

那么这个时候我就可以出实话一个 ClassName，然后就可以调用方法 MethodName 了。

对于上面的方法定义，其实还有一种写法是这样的：

```go
func (typeName *ClassName) MethodName(parameter string) string {

}
```

这个我翻阅了下资料，发现这个还是很有趣的，这个跟编译器以及 golang 本身的 “函数式编程” 的特性有关。函数式有个很重要的特征，就是 “无状态” 的。举个例子，我新建一个函数，这个函数本身是无状态的，只要你传入的参数不变，那么这个函数得到的值就是恒定值。我们拿之前的方法为例子 `func (typeName ClassName) MethodName(parameter string) string {}` 这个就是说你无法更改 typeName 这个值，它是不变量的。你只能在这个函数领域下更改，一旦这个方法返回（指定的栈地址）则传入的 typeName 就还是之前的状态。如果你要像 C# 一样传递一个引用，在局部更改返回后，这个引用对象同样也会更改的话，改怎么实现呢？也很简单，只要在原来的基础之上加个 “ * ”，也就是上面的写法。
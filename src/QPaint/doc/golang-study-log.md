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

**怎么导入第三方库？**

直接在官方仓库地址选择自己要导入的模块，地址见：https://pkg.go.dev

```go
import ("moduleName")
```

但是我当初这么做的时候，发现虽然编译器没有检测到错误，但是在引用 module 的 api 时，没有智能提示。后来发现跟 nuget 是一样的，有快捷键导入 module: `shift + command + p` 然后选择 `Go: Add Import` 选中你本地 clone 下来的第三方类库。但是这里有一个疑问，难道要每次下载源代码吗？而是不能够直接下载一个类似 dll 的可执行的 “微文件” 么。讲道理是肯定有的，不可能一个项目发布出去，还把人家的源代码发不出去的。这个之后弄到发布的时候在回过头来查资料吧

错误1: `expected ';', found f` 这个是不同编辑器的编码问题，拿我现在用的 vscode 为例，我新建的 go 项目的默认编码是 `LF/CRLF`，切换成 `CRLF/LF` 。注意，如果切换的时候发现还是报同样的错误，那极有可能是 vscode 没有反应过来，只要在当前页面随便输个空格在保存即可。



接下来就是各种变量函数的基本用法介绍了

定义类：用 type 关键字 `type ClassName struct{ someField int}`

申明变量：`var scopeVar = "string"` 我试了一下，这种显示写法也行 `var scopeVar string = ""`。

但是有这么一种写法 `localVar := ""`，也很有趣，我尝试了一下，这个好像只能在方法里面写（就相当于 `var localVar string = ""`），这种写法提升到 “全局” 则不行。

函数申明：

函数申明分两种

1. 无返回值：`func SomeMethod(){}`
   1. 带参数：`func SomeMethod(a int)`




# 深入理解泛型

设计泛型的一个准要原则

> do you want slow programmers, slow compilers and bloated binaries, or slow execution times?
>
> 你是想要缓慢的程序，慢编译器和臃肿的二进制，又或是慢执行时间？

“慢速编译器和臃肿的二进制文件”指的是通过完全单态化实现模板的 C++ 方法——每个模板调用都被视为有点像宏扩展，其完整的代码复制了其正确的类型。

“慢速执行时间”指的是Java的装箱方法，或者指的是由于每次调用都透明地动态分发，代码非常通用。

这些特性都是tradeoff。

从go的两篇泛型设计文档来看：[泛型实现——Stenciling设计文档](https://go.googlesource.com/proposal/+/refs/heads/master/design/generics-implementation-stenciling.md)；[泛型实现——字典设计文档](https://go.googlesource.com/proposal/+/refs/heads/master/design/generics-implementation-dictionaries.md)；go都考虑这两种泛型设计方向

> 在go语言中，stenciling 就是 c++ 的单态化（monomorphization ）泛型模板；字典（dictionaries）就是动态分发

由于上述原因，这两种方法本身都不完美。因此，提出了另一种设计：

[泛型实现 - GC Shape Stenciling](https://github.com/golang/proposal/blob/master/design/generics-implementation-gcshape.md)

这种“GC Shape”方法是模板和字典两种极端方法的折衷。根据实例化的类型，我们可以单态化或使用动态调度。有一个最新的[文档](https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md)详细描述了Go 1.18是如何做到这一点的。

具体来说，不同的底层类型（如整数和字符串）将获得自己的 GC Shape，**这意味着将为每种类型生成不同的函数，并且类型是硬编码的（因此这是单态化）。另一方面，所有指针类型都将分组在同一个 GC Shape 中，并将使用动态调度。**

> **All pointers to objects belong to the same GCShape, regardless of the object being pointed at**.
>
> 这意味着 time.Time 指针与 uint64、bytes.Buffer 和 trings.Builder 具有相同的 GCShape。但是GC Shape并不知道在具体调用方法时发生的事情。

要注意，这是在目前 Go1.18 中的状态，后续可能会发生变化，因为 Go 团队正在与社区合作，以了解什么最适合现实生活中的工作负载。

## 为什么泛型排序版本要比非泛型性能要好？

如前一节所讨论的，字符串类型将获得自己的GC Shape，因此将为该字符串类型硬编码自己的函数。让我们看看它在程序集中是什么样子的。

首先，翻找二进制文件的调试信息，我们会找到这个符号：`bubbleSortGeneric[go.shape.string_0]`，它表示该字符串当前唯一成员的 GC Shape 的 bubbleSortGeneric 的模板版本。但是，我们不会发现它是一个独立的函数来调用，**因为它被内联到它的调用站点中**。这种内联不会影响性能，因此我们将只关注内部循环的指令，提醒您这样做：

```go
for i := 1; i < n; i++ {
  if x[i] < x[i-1] {
    x[i-1], x[i] = x[i], x[i-1]
    swapped = true
  }
}
```

它生成的汇编代码如下：

```assembly
MOVQ  0x80(SP), R8
INCQ  R8
MOVQ  0x70(SP), CX
MOVQ  0x78(SP), BX
MOVQ  R8, DX
MOVL  AX, SI
MOVQ  0xb0(SP), AX
CMPQ  DX, BX
JLE   0x4aef20
MOVQ  DX, 0x80(SP)
MOVB  SI, 0x3d(SP)
MOVQ  DX, SI
SHLQ  $0x4, DX
MOVQ  DX, 0x90(SP)
MOVQ  0(DX)(AX*1), R8
MOVQ  0x8(DX)(AX*1), BX
LEAQ  -0x1(SI), R9
SHLQ  $0x4, R9
MOVQ  R9, 0x88(SP)
MOVQ  0(R9)(AX*1), CX
MOVQ  0x8(R9)(AX*1), DI
MOVQ  R8, AX
CALL  runtime.cmpstring(SB)
MOVQ  0xb0(SP), DX
MOVQ  0x90(SP), SI
LEAQ  0(SI)(DX*1), DI
MOVQ  0x88(SP), R8
LEAQ  0(R8)(DX*1), R9
TESTQ AX, AX
JGE   0x4af01a
```

首先要注意的是，Less方法没有动态分派。每次循环迭代都直接调用cmpstring。其次，程序集的后一部分类似于前面显示的 Less 代码，有一个关键的区别—**没有边界检查**！Go 包含一个[边界检查消除(BCE)通道](https://go101.org/article/bounds-check-elimination.html)，它可以消除比较时的边界检查:

```go
// ... earlier we had n := len(x)
for i := 1; i < n; i++ {
  if x[i] < x[i-1] {
```

编译器知道 i 在任何时候都在 1 和 `len(x)` 之间（通过查看循环描述和 i 未被修改的事实），因此 x[i] 和 x[i-1] 都在安全地访问切片边界。

在接口版本中，编译器不会消除 Less 中的边界检查；函数是这样定义的:

```go
func (x StringSlice) Less(i, j int) bool { return x[i] < x[j] }
```

谁知道传入的索引值是多少呢！此外，由于动态调度，这个函数没有内联它的调用者，编译器可能对正在发生的事情有更深入的了解。Go 编译器具有一些去虚拟化(devurtualization capabilities)的功能，但它们并没有在这里发挥作用。这是编译器改进的另一个有趣领域。

> BCE（Bounds Check Elimination）消除边界检查：
>
> 在数组/切片元素操作中，程序会检查所有涉及的索引是否超出范围。如果超出范围就会引起程序崩溃，于是通过边界检查来消除这种破坏。这就是边界检查。
>
> 这种检查能使程序运行得更安全，但是性能也会有所损失。
>
> 从Go Toolchain 1.7开始，标准的Go编译器使用了一个新的编译器后端，它基于SSA(静态单赋值形式)。SSA帮助Go编译器有效地使用像[BCE(边界检查消除)](https://en.wikipedia.org/wiki/Bounds-checking_elimination)和[CSE(公共子表达式消除)](https://en.wikipedia.org/wiki/Common_subexpression_elimination)这样的优化。BCE可以避免一些不必要的边界检查，CSE可以避免一些重复的计算，使标准的Go编译器可以生成更高效的程序。有时这些优化的改进效果是明显的。

## 泛型排序与自定义排序委托

为了验证前面描述的一些观察结果，这一次，不依赖`constraints.Ordered`，但使用比较函数代替：

```go
func bubbleSortFunc[T any](x []T, less func(a, b T) bool) {
  n := len(x)
  for {
    swapped := false
    for i := 1; i < n; i++ {
      if less(x[i], x[i-1]) {
        x[i-1], x[i] = x[i], x[i-1]
        swapped = true
      }
    }
    if !swapped {
      return
    }
  }
}
```

通过如下调用排序：

```go
bubbleSortFunc(ss, func(a, b string) bool { return a < b })
```

其性能如下：

```powershell
go test -bench . compare-sort-generic 
goos: windows
goarch: amd64
pkg: compare-sort-generic
cpu: Intel(R) Core(TM) i5-9400 CPU @ 2.90GHz
BenchmarkSortStringInterface-6               130           9115803 ns/op
BenchmarkSortStringGeneric-6                 159           7445174 ns/op
BenchmarkSortStringFunc-6                    141           8446509 ns/op
```

从结果看比较函数排序的性能居中。

这种比较具有有趣的现实意义，因为 SortFunc 也是添加到 `golang.org/exp/slices` 的变体，以提供更通用的排序功能（对于不受约束的类型）。此版本还提供了针对 `sort.Sort` 的加速。

另一个含义是对指针类型进行排序；如前所述，1.18 中的 Go 编译器会将所有指针类型分组到一个 GC Shape 中，这意味着它需要传递一个字典来进行动态分发。这可能会使代码变慢，尽管 BCE 仍然应该启动 - 所以不会慢很多。

## 注意事项

并不是所有都上泛型都是有利的，有些场景使用泛型会使系统变慢，具体详见：https://planetscale.com/blog/generics-can-make-your-go-code-slower；

Go后续的版本更新可能会修复已有的问题，但是一定要记住：<font color="red">另一方面，由于[《泛型困境》](*Generic Dilemma*)中所讨论的原因，go 的泛型不太可能在所有可能的情况下都是“零成本”的。Go优先考虑快速的编译时间和紧凑的二进制大小，因此它必须在任何设计中做出一定的权衡。</font>

## 原文地址

https://eli.thegreenplace.net/2022/faster-sorting-with-go-generics/



# 泛型会使你的代码变慢

在 1.18 中的当前泛型实现中，泛型函数的每次运行时调用都将透明地接收一个静态字典作为其第一个参数，其中包含有关传递给函数的参数的元数据。该字典将放置在 AMD64 的寄存器 AX 中，以及 Go 编译器尚不支持基于寄存器的调用约定的平台中的堆栈中。这些字典的完整实现细节在上述设计文档中进行了深入解释，但作为总结，它们包括所有必需的类型元数据，以将参数传递给进一步的泛型函数，将它们从接口转换或转换为接口，对于我们来说，最重要的是调用它们的方法。没错，在单态化（monomorphization）步骤之后，生成的函数 Shape 需要将其所有泛型参数的虚拟方法表（virtual method table）作为运行时输入。直观地说，虽然这大大减少了生成的唯一代码的数量，但这种广泛的单态化并不适合去虚拟化（de-virtualization）、内联或任何类型的性能优化。

事实上，对于绝大多数的 Go 代码来说，让它泛型就意味着让它变慢。

> 去虚拟化（de-virtualization）：// TODO

接口是一种涉及装箱的多态形式，即确保我们操作的所有对象具有相同的Shape。对于 Go 接口，这个 Shape 是一个 16 字节的胖指针（iface），其中前半部分指向有关装箱值的元数据（我们称之为 itab），后半部分指向值本身。

```go
type iface struct {
	tab *itab
	data unsafe.Pointer
}

type itab struct {
	inter *interfacetype // offset 0
	_type *_type // offset 8
	hash  uint32 // offset 16
	_     [4]byte
	fun   [1]uintptr // offset 24...
}
```

itab 包含大量关于接口内部类型的信息。inter、_type 和 hash 字段包含所有必需的元数据，以允许接口之间的转换、反射和切换接口的类型。但是这里我们关心的是 itab 末尾的 fun 数组：虽然在类型描述中显示为 `[1]uintptr`，但这实际上是一个变长分配。 itab 结构的大小在特定接口之间变化，结构末尾有足够的空间来存储接口中每个方法的函数指针。这些函数指针是我们每次要调用接口上的方法时需要访问的；它们在 Go 中等价于 C++ 虚拟表。

举个例子，非泛型版本`buf.WriteByte('\\')`生成了如下代码：

```assembly
0089  MOVQ "".buf+48(SP), CX
008e  MOVQ 24(CX), DX
0092  MOVQ "".buf+56(SP), AX
0097  MOVL $92, BX
009c  CALL DX
```

要在 buf 上调用 WriteByte 方法，我们首先需要一个指向 buf 的 itab 的指针。尽管 buf 最初是通过一对寄存器传递到我们的函数中，但编译器在函数体的开头将其溢出到堆栈中，以便它可以将寄存器用于其他事情。要调用 buf 上的方法，我们首先必须将 *itab 从堆栈中加载回寄存器 (CX)。现在，我们可以取消引用 CX 中的 itab 指针来访问它的字段：我们将偏移 24 处的双字移动到 DX 中，快速浏览一下上面 itab 的原始定义表明，事实上，itab 中的第一个函数指针位于偏移量 24 处——到目前为止，这一切都说得通。

由于 DX 包含我们要调用的函数的地址，我们只是缺少它的参数。Go 所谓的“结构附加方法(struct-attached method)”是对一个独立函数的[糖](https://en.wikipedia.org/wiki/Syntactic_sugar)，该函数将其接收者作为其第一个参数，例如 `func (b *Builder) WriteByte(x byte)` 脱糖（desugars）到 `func "".(*Builder).WriteByte(b *Builder, x byte)`。因此，函数调用的第一个参数必须是 `buf.(*iface).data`，它是指向位于我们接口内的 `strings.Builder` 的实际指针。该指针在堆栈中可用，在我们刚刚加载的制表符指针(tab pointer)之后的 8 个字节。最后，我们函数的第二个参数是字面量 `\\`, (ASCII 92)，我们可以调用 DX 来执行我们的方法。

再来看看泛型版本：

```assembly
MOVQ ""..dict+48(SP), CX
0094  MOVQ 64(CX), CX
0098  MOVQ 24(CX), CX
009c  MOVQ "".buf+56(SP), AX
00a1  MOVL $92, BX
00a6  CALL CX
```

它看起来很相似，但有一个明显的区别。偏移量 `0x0094` 包含我们不希望函数调用站点包含的内容：另一个指针解引用。这里发生了什么事情：**由于我们将所有指针 shape 单态化为 *uint8 的单个形状实例化，因此该形状不包含有关可以在这些指针上调用的方法的任何信息。**这些信息将保存在哪里？理想情况下，它将存在于与我们的指针关联的 itab 中，但没有与我们的指针直接关联的 itab，因为我们函数的形状采用单个 8 字节指针作为其 buf 参数，而不是 16 字节胖指针 `*itab` 和 `data` 字段，就像接口一样。如果您还记得，这就是 stenciling 实现将字典传递给每个泛型函数调用的全部原因：该字典包含指向函数的所有泛型参数的 `itab` 的指针。

好了，这个程序集，加上额外的负载，现在完全讲得通了。方法调用的开始，不是加载我们的 buf 的 itab，而是加载已传递给我们的泛型函数（并且也已溢出到堆栈中）的字典。使用 CX 中的字典，我们可以解引用它，并且在偏移量 64 处我们找到了我们正在寻找的 `*itab`。遗憾的是，我们现在需要另一个解引用 (`24(CX)`) 来从 `itab` 内部加载函数指针。方法调用的其余部分与前面的代码生成相同。

这种额外的解引用在实践中有多糟糕？直观上，我们可以假设在泛型函数中调用对象的方法总是比在只接受接口作为参数的非泛型函数中慢，<font color="red">因为泛型将把以前的指针调用转变为两次间接的接口调用，表面上比普通接口调用慢</font>。

```
name                      time/op      alloc/op     allocs/op
Monomorphized-16          5.06µs ± 1%  2.56kB ± 0%  2.00 ± 0%
Iface-16                  6.85µs ± 1%  2.59kB ± 0%  3.00 ± 0%
GenericWithPtr-16         7.18µs ± 2%  2.59kB ± 0%  3.00 ± 0%
```

​														（上结果引自作者）

这个简单的基准测试使用3个略有不同的实现来测试相同的函数体。`GenericWithPointer` 将 `strings.Builder` 传递给 `func Escape[W io.ByteWriter](W, []byte)` 泛型函数。`Iface` 的基准测试是直接采用接口的 `func Escape(io.ByteWriter, []byte)`。`Monomorphized` 用于手动单态化的 `func Escape(*strings.Builder, []byte)`。

结果不令人意外。专门用于直接获取 `strings.Builder` 的函数是最快的，因为它允许编译器在其中**内联** `WriteByte` 调用。泛型函数比以 `io.ByteWriter` 接口作为参数的最简单实现慢得多。我们可以看到，来自泛型字典的额外负载的影响并不显着，因为在这个微基准测试中，`itab` 和泛型字典在缓存中都会非常温暖（swarm）（但是，请继续阅读以分析缓存争用如何影响泛型代码）。

这是我们可以从该分析中收集到的第一个见解：**在 1.18 中没有动力将带有接口的纯函数转换为使用泛型。它只会让它变慢，因为 Go 编译器目前无法生成通过指针调用方法的函数shape。相反，它将引入具有两层间接的接口调用。这与我们想要的方向完全相反，即去虚拟化，并在可能的情况下进行内联。**

在结束本节之前，让我们指出 Go 编译器逃逸分析中的一个细节：我们可以看到我们的单态化函数在我们的基准测试中有 2 个 allocs/op。这是因为我们传递了一个指向位于堆栈中的 `strings.Builder` 的指针，编译器可以证明它没有逃逸，因此不需要堆分配。即使我们还从堆栈中传递了一个指针，Iface 基准测试也显示了 3 个 allocs/op。这是因为我们将指针移动到接口，并且总是分配。令人惊讶的是，GenericWithPointer 实现还显示了 3 个 allocs/op。即使为函数生成的实例化直接采用指针，转义分析也不能再证明它是非转义的，因此我们获得了额外的堆分配。

## 泛型接口调用

如果我们将 `strings.Builder` 隐藏在接口后面会发生什么？

```go
var buf strings.Builder
var i io.ByteWriter = &buf
BufEncodeStringSQL(i, []byte(nil))
```

泛型函数的参数现在是一个接口，而不是指针。这个调用显然是有效的，因为我们传递的接口与方法上的约束是相同的。但是我们生成的实例化shape是什么样的呢？我们没有嵌入完整的反汇编，因为它会产生混乱，但就像我们之前做的一样，让我们分析函数中`WriteByte`方法的调用点:

```assembly
00b6  LEAQ type.io.ByteWriter(SB), AX
00bd  MOVQ ""..autotmp_8+40(SP), BX
00c2  CALL runtime.assertI2I(SB)
00c7  MOVQ 24(AX), CX
00cb  MOVQ "".buf+80(SP), AX
00d0  MOVL $92, BX
00d5  CALL CX
```

变化很大，这与我们之前生成的代码相比，差距很大。

这发生了什么？我们可以在 Go 运行时中找到 [runtime.assertI2I](https://github.com/golang/go/blob/c6d9b38dd82fea8775f1dff9a4a70a017463035d/src/runtime/iface.go#L421-L430) 方法：这是一个断言接口之间转换的帮助方法。它接受一个 `*interfacetype` 和一个 `*itab` 作为它的两个参数，并且仅当给定 `itab` 中的接口也实现了我们的目标接口时才返回给定 `interfacetype` 的 `itab`。

假设有这样一个接口：

```go
type IBuffer interface {
	Write([]byte) (int, error)
	WriteByte(c byte) error
	Len() int
	Cap() int
}
```

这个接口没有提到`io.ByteWriter`和`io.Writer`，然而任何实现`IBuffer`的类型也隐士的实现了这两个接口。这对我们泛型函数生成的代码会有影响：因为我们对函数的泛型约束是`io.ByteWriter`，我们可以任何一个实现了`io.ByteWriter`的接口参数——当然这也包括了`IBuffer`。但是我们需要在参数上调用`WriteByte`方法时，该方法在我们收到的接口`itab.fun`数组中的哪个位置存在？我们并不知道！如果我们传递的是`*strings.Builder`作为`io.ByteWriter`接口参数，那么这个方法就会在该接口（指strings.Builder）的 `itab`将在`func[0]`的位置。如果我们传递的是`IBuffer`，它就会在`func[1]`中。我们需要一个帮助类，它可以将 `itab` 用作 `IBuffer`，并返回一个 `io.ByteWriter` 的 `itab`，其中我们的 `WriteByte` 函数指针始终稳定在 fun[0]。

这是`assertI2I`的工作，也是函数中的每个调用站点所做的。让我们一步一步来分析。

```assembly
00b6  LEAQ type.io.ByteWriter(SB), AX
00bd  MOVQ ""..autotmp_8+40(SP), BX
00c2  CALL runtime.assertI2I(SB)
00c7  MOVQ 24(AX), CX
00cb  MOVQ "".buf+80(SP), AX
00d0  MOVL $92, BX
00d5  CALL CX
```

首先，**它将 io.ByteWriter 的接口类型（这是一个硬编码的全局，因为这是我们的约束中定义的接口类型）加载到寄存器 `AX`。然后，它将我们传递给函数的接口的实际`itab`加载到`BX`中**。这就是函数`assertI2I`需要的两个参数，在调用它之后剩下的就是在`AX`中的`io.ByteWriter`的`itab`，我们可以继续我们的接口函数调用，就像我们在之前的代码生成中所做的那样，知道我们的函数指针现在总是在`itab`中偏移24。本质上，这个shape实例化所做的就是将每个方法调用从`buf.WriteByte(ch)`转换为`buf.(io.ByteWriter). writebyte(ch)`。

## 字节序列

翻看go源码的可以直到，以`[]byte`切片作为参数有大量重复的代码（如`(*Buffer).Write`和`(*Buffer).WriteString`）,其中以`encoding/utf8`这个包中的代码有超过50%的API代码都是重复的

| Bytes            | String                   |
| ---------------- | ------------------------ |
| `DecodeLastRune` | `DecodeLastRuneInString` |
| `DecodeRune`     | `DecodeRuneInString`     |
| `FullRune`       | `FullRuneInString`       |
| `RuneCount`      | `RuneCountInString`      |
| `Valid`          | `ValidString`            |

需要指出的是，这种代码重复实际上是一种性能优化：API可以很好的只提供它想要的`[]byte`切片数据，这样就会迫使用户在调用包之前必须转换为`[]byte`，并且它们的互相转换会额外强制分配内存。

所以对于重复的代码来看，使用泛型是一个很好的目标，但是由于代码重复首先是为了防止额外的分配，所以在实现泛型统一之前，必须要明确生成的shape实例的行为是符合预期的。

> 当利用泛型优化时，我们可以通过查看汇编代码，看生成的指令是否有很大差别

具体的例子给查看末尾给的原文

## 结论（1.18泛型最佳实践建议）

- 建议减少重复相同的那些使用`ByteSeq`约束对采用`string`和`byte[]`的方法。这生成的 shape 实例化非常接近于手动编写两个几乎相同的函数。
- **建议使用泛型数据结构**。这是迄今为止他们最好的用例：以前使用 `interface{}` 实现的通用数据结构复杂且不符合人体工程学。**删除类型断言并以类型安全的方式存储未装箱的类型，使这些数据结构更易于使用且性能更高。**
- 请尝试通过回调类型参数化功能助手。在某些情况下，它可能允许 Go 编译器将它们展平（flatten）。
- 不要尝试使用泛型去虚拟化或内联方法调用。它不起作用，因为所有指针类型都有一个可以传递给泛型函数的相同 shape；相关的方法信息存在于运行时字典中。
- 在任何情况下都不要将接口传递给泛型函数。由于 shape 实例化适用于接口的方式，而不是去虚拟化，您要添加另一个虚拟化层，该层涉及每个方法调用的全局哈希表查找。在性能敏感的上下文中处理泛型时，只使用指针而不是接口。
- 不要重写基于接口的 API 来使用泛型。鉴于当前实现的限制，如果继续使用接口，当前使用非空接口的任何代码都会表现得更可预测，并且会更简单。在方法调用方面，泛型将指针转化为两次间接接口，并将接口进而转化为......好吧，如果我说实话，这是非常可怕的事情。
- 不要觉得失望或是高兴，因为 Go 泛型的语言设计没有技术限制，可以防止（最终）实现更积极地使用单态化来内联或去虚拟化方法调用。

## 原文链接

https://planetscale.com/blog/generics-can-make-your-go-code-slower


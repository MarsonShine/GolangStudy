# Go1.8——泛型支持

## 优劣

对于语言的任何特性更改，我们都需要讨论利益最大化和成本最小化。

在 Go 中，我们的目标是通过可以自由组合的独立的、正交的（orthogonal）语言特性来降低复杂性。我们通过简化单个功能来降低复杂性，并通过允许它们自由组合来最大化功能的好处。我们想对泛型做同样的事情。

为了更具体地说明这一点，我将列出一些我们应该遵循的指导方针。

### 减少新概念

我们应该尽可能少地向语言中添加新概念。这意味着最少的新语法，最少的新关键字和其他名称。

### 复杂性落于泛型作者本身，而不是使用泛型的用户身上

编写泛型包的程序员应该尽可能地承担复杂性。我们不希望包的用户担心泛型。这意味着应该能够以一种自然的方式调用泛型函数，**这意味着使用泛型包时的任何错误都应该以一种易于理解和修复的方式报告。它还应该易于调试对泛型代码的调用。**

### 泛型作者与用户都能独立工作

同样，我们应该很容易地将泛型代码的作者和用户的关注点分开，这样他们就可以独立地开发自己的代码。它们不应该担心对方在做什么，就像不同包中的普通函数的编写者和调用者不需要担心一样。这听起来很明显，但在其他所有编程语言中，泛型并不是这样的。

### 构建时间短，执行时间快

当然，我们希望尽可能地保持 Go 现在提供给我们的短构建时间和快速执行时间。泛型倾向于在快速构建和快速执行之间进行权衡。尽可能地，我们两者都想要。

### 保持 Go 的清晰和简单

最重要的是，Go 至今都是一门简单的语言。Go 程序通常清晰易懂。在我们探索泛型领域的漫长过程中，一个主要部分是试图理解如何在**保持清晰和简单的同时添加泛型**。我们需要找到能够很好地适应现有语言的机制，而不是把它变成完全不同的东西。

这些准则应该适用于 Go 中的任何泛型实现。这是我今天想给你们传达的最重要的信息：**泛型可以为语言带来显著的好处，但是只有能保持 Go 的核心理念（Go still feels like Go），它们才值得去做。**

## 设计草案

幸运的是，我认为这是可以做到的。在这篇文章的最后，我将从讨论为什么我们需要泛型以及转到泛型的要求，简单地讨论如何将它们添加到语言中的设计。

这是这个设计中通用的 Reverse 功能实现：

```go
func Reverse (type Element) (s []Element) {
    first := 0
    last := len(s) - 1
    for first < last {
        s[first], s[last] = s[last], s[first]
        first++
        last--
    }
}
```

你会注意到函数的主体是完全相同的。只有签名变了。

切片的元素类型 `type Element` 已经提出来了。称为 `Element` 并且我们称它为*类型参数（type parameter）*。它不再是 slice 参数类型的一部分，而是一个单独的、附加的类型参数。

要用类型形参调用函数，一般情况下你要传递一个类型参数，除了它是一个类型外，它和其他参数一样。

```go
func ReverseAndPrint(s []int) {
    Reverse(int)(s)
    fmt.Println(s)
}
```

这就是 `int` 的 `Reverse` 的实现例子。

幸运的是，在大多数情况下，包括这个例子，**编译器可以从常规参数的类型推断出类型参数，而且您根本不需要提到类型参数**。

调用一个泛型函数就像调用任何其他函数一样。

```go
func ReverseAndPrint(s []int) {
    Reverse(s)
    fmt.Println(s)
}
```

换句话说，尽管通用的 `Reverse` 函数比 `ReverseInt` 和 `ReverseStrings` 稍微复杂一些，但复杂性落在函数的编写者身上，而不是调用者。

### 契约（Contracts）

因为 Go 是一种静态类型语言，所以我们必须讨论类型参数的类型。这个元类型（meta-type）告诉编译器在调用泛型函数时允许什么类型的实参，以及泛型函数可以对类型形参的值执行什么类型的操作。

`Reverse` 函数可以处理任何类型的切片。它对 `Element` 类型的值所做的唯一的事情是赋值，这适用于 Go 中的任何类型。对于这种泛型函数，这是一种非常常见的情况，我们不需要对类型参数说任何特殊的话。

让我们看一下另一个函数。

```go
func IndexByte (type T Sequence) (s T, b byte) int {
    for i := 0; i < len(s); i++ {
        if s[i] == b {
            return i
        }
    }
    return -1
}
```

目前标准库中的 `bytes` 包和 `strings` 包都有一个 `IndexByte` 函数。这个函数返回 `s` 序列中的索引 `b`，其中 `s` 是一个字符串或一个 `[]byte`。我们可以使用单个泛型函数来替换 `bytes` 和 `strings` 包中的两个函数。在实践中，我们可能不会费心这么做，但这是一个有用的简单例子。

这里，我们需要知道类型参数 `T` 的作用类似于字符串或字节数组。我们可以在它上面调用 `len`，我们可以对它进行索引，我们可以将索引操作的结果与一个字节值进行比较。

要让这个编译，类型参数 `T` 本身需要一个类型。它是一个元类型（meta-type），但是因为我们有时需要描述多个相关的类型，而且因为它描述了泛型函数的实现和它的调用者之间的关系，所以我们实际上把类型 `T` 称为契约。在这个例子，这个契约被命名为 `Sequence`。它出现在类型参数列表之后。

这就是本例中 `Sequence` 契约的定义方式。

```go
contract Sequence(T) {
    T string, []byte
}
```

这非常简单，因为这是一个简单的示例：类型参数 `T` 既可以是字符串，也可以是字节数组。这里的 `contract` 可能是新加的关键字，或者在包范围内识别的特殊标识符；具体详见设计稿。

原来的 contract 的设计方案，使用者反馈太过复杂（具体详见[2018年 Gophercon 大会上展示的设计](https://github.com/golang/proposal/blob/4a530dae40977758e47b78fae349d8e5f86a6c0a/design/go2draft-contracts.md)），现在的 contract 就非常简单了。

它们允许您指定类型参数的基础类型，以及/或列出类型参数的方法。它们还允许您描述不同类型参数之间的关系。

### 方法契约

下面是另一个简单的例子，这个函数使用 `String` 方法返回一个字符串数组，它是 `s` 中所有元素的。

```go
func ToStrings (type E Stringer) (s []E) []string {
    r := make([]string, len(s))
    for i, v := range s {
        r[i] = v.String()
    }
    return r
}
```

它非常简单：遍历切片，在每个元素上调用 `String` 方法，并返回结果字符串的一个切片。

你可能注意到这个合同看起来像 `fmt.Stringer` 接口，所以值得指出的是 `ToStrings` 函数的参数不是 `fmt.Stringer` 的一个切片。**它是某个元素类型的切片，其中元素类型实现了 `fmt.Stringer`。**元素类型切片以及 `fmt.Stringer` 的切片再内存表示上通常是不同的。Go 不支持它们之间的直接转换。所以这是值得写的，即使是 `fmt.Stringer` 已经存在。

### 多类型的契约

下面是一个带有多个类型参数的契约示例。

```go
type Graph (type Node, Edge G) struct { ... }

contract G(Node, Edge) {
    Node Edges() []Edge
    Edge Nodes() (from Node, to Node)
}

func New (type Node, Edge G) (nodes []Node) *Graph(Node, Edge) {
    ...
}

func (g *Graph(Node, Edge)) ShortestPath(from, to Node) []Edge {
    ...
}
```

这里我们描述的是一个由节点和边组成的图。我们不需要图的特定数据结构。相反，我们说的是 `Node` 类型必须有一个 `Edges` 方法来返回连接到 `Nodes` 的边列表。`Edge` 类型必须有一个 `Nodes` 方法，它返回 `Edge` 连接的两个 node。

我跳过了这个实现，但它显示了返回 `Graph` 的 `New` 函数的签名，以及 `Graph` 上的 `ShortestPath` 方法的签名。

这里重要的一点是，契约不只是关于单一类型。它可以描述两个或多个类型之间的关系。

### 有序类型（Ordered types）

一个令人惊讶的普遍抱怨是 Go 没有 `Min` 或者是 `Max` 函数。这是因为一个有用的 `Min` 函数应该适用于任何有序类型，这意味着它必须是泛型。

尽管自己编写 `Min` 非常简单，但任何有用的泛型实现都应该允许我们将其添加到标准库中。这就是我们的设计。

```go
func Min (type T Ordered) (a, b T) T {
    if a < b {
        return a
    }
    return b
}
```

`Ordered` 契约说类型 `T` 必须是有序类型，这意味着它支持小于、大于等操作符。

```go
contract Ordered(T) {
    T int, int8, int16, int32, int64,
        uint, uint8, uint16, uint32, uint64, uintptr,
        float32, float64,
        string
}
```

`Ordered` 契约只是该语言定义的所有有序类型的列表。此契约接受任何列出的类型，或任何基础类型为上述类型之一的指定类型。基本上，您可以使用 `<` 操作符的任何类型。

事实证明，简单地枚举支持小于操作符的类型要比发明一种适用于所有操作符的新表示法容易得多。毕竟，在 `Go` 中，只有内置类型支持操作符。

同样的方法可以用于任何操作符，或者更普遍地用于为任何打算使用内置类型的泛型函数编写契约。它让泛型函数的编写人员清楚地指定函数要使用的类型集。它让泛型函数的调用者清楚地看到该函数是否适用于所使用的类型。

在实践中，这个契约可能会进入标准库，因此实际上 `Min` 函数(它可能也会在标准库的某个地方)看起来是这样的。这里我们只是谈及了契约包中中定义的 Ordered 契约。

```go
func Min (type T contracts.Ordered) (a, b T) T {
    if a < b {
        return a
    }
    return b
}
```

### 泛型数据结构

最后，让我们来看一个简单的泛型数据结构，二叉树。在这个例子中，树有一个比较函数，所以对元素类型没有要求。

```go
type Tree (type E) struct {
    root    *node(E)
    compare func(E, E) int
}

type node (type E) struct {
    val         E
    left, right *node(E)
}
```

下面是如何创建一个新的二叉树。比较函数被传递给 `New` 函数

```go
func New (type E) (cmp func(E, E) int) *Tree(E) {
    return &Tree(E){compare: cmp}
}
```

其中一个私有方法返回的指针要么指向存放 `v` 的槽，要么指向树中它应该在的位置。

```go
func (t *Tree(E)) find(v E) **node(E) {
    pn := &t.root
    for *pn != nil {
        switch cmp := t.compare(v, (*pn).val); {
        case cmp < 0:
            pn = &(*pn).left
        case cmp > 0:
            pn = &(*pn).right
        default:
            return pn
        }
    }
    return pn
}
```

这里的细节并不重要，特别是因为我还没有测试这段代码。我只是想展示一下写一个简单的泛型数据结构是什么样子的。

下面用于测试树是否包含值的代码

```go
func (t *Tree(E)) Contains(v E) bool {
    return *t.find(e) != nil
}
```

下面代码展示插入新值

```go
func (t *Tree(E)) Insert(v E) bool {
    pn := t.find(v)
    if *pn != nil {
        return false
    }
    *pn = &node(E){val: v}
    return true
}
```

注意类型节点的类型参数 `E`。这就是编写泛型数据结构的样子。正如您所看到的，它看起来像编写普通的 Go 代码，除了一些类型参数散落地出现在这里和那里。

使用起来非常简单

```go
var intTree = tree.New(func(a, b int) int { return a - b })

func InsertAndCheck(v int) {
    intTree.Insert(v)
    if !intTree.Contains(v) {
        log.Fatalf("%d not found after insertion", v)
    }
}
```

就是应该这样的。编写泛型数据结构有点困难，因为您常常必须显式地写出支持类型的类型参数，但是尽可能地使用一个泛型数据结构与使用普通的非泛型数据结构没有什么不同。

### 下一步

我们正在进行实际的实现，以允许我们对这个设计进行实验。能够在实践中尝试这些设计是很重要的，以确保我们能够编写我们想要编写的程序。它并没有像我们希望的那样快，但当这些实现可用时，我们会发布更多的细节。

Robert Griesemer 编写了一个[初步的CL](https://go.dev/cl/187317)，修改了 `go/types` 包。这允许测试使用泛型和契约的代码是否可以进行类型检查。它目前还不完整，但它主要适用于单个包，我们将继续努力。

我们希望人们对这个和未来的实现做的是尝试编写和使用泛型代码，看看会发生什么。我们希望确保人们能够写出他们需要的代码，并且能够按照预期使用这些代码。当然，并不是所有的事情在一开始都能顺利进行，随着我们探索这个领域，我们可能不得不改变一些事情。而且，要清楚的是，我们对语义的反馈比对语法的细节更感兴趣。

# 原文链接

https://go.dev/blog/why-generics


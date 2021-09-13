# Go Slices 切片

Go slices 切片类型为处理类型化数据序列提供了一种方便而有效的方法。Slices 非常类似于其它语言的数组类型，但又有一些非常用的属性。这篇文章主要就是介绍 slices 是如何使用的。

## 数组

slices 是一个构建在 Go 数组类型之上的一种抽象类型，所以为了更好的理解 slice 类型，我们有必要先理解数组类型。

数字类型指定了数组长度与元素类型。例子，`[4]int` 就表示含有四个 int 元素的数组。数组的尺寸是固定的；它的长度是它类型的一部分（`[4]int` 和 `[5]int` 是不同的，无法兼容。）数组通常能通过索引方式查询，即表达式 `s[n]` 表示访问是从 0 开始的第 n 个元素。

```go
var a [4]int
a[0] = 1
i := a[0]
// i == 1
```

数组不需要显式的初始化；数组元素的零值是默认初始值：

```
// a[2] == 0, int 类型的默认值是 0
```

`[4]int` 在内存的表示就是 4 个整数值顺序排列：

![](https://go.dev/blog/slices-intro/slice-array.png)

Go 数组是值。数组变量表示是整个数组；它不是指向第一个数组元素的指针（就像 C 语言中的一样）。就是说当你要分配或传递一个数组值时，会对数组的内容进行拷贝。（为了避免拷贝，你可以传递一个数组指针，但是这个是指向的数组的指针而不是数组）。考虑数组的一种方法是，它是一种结构体，但带有索引字段而不是命名字段：一个固定大小的复合值。

数组字面量可以如下指定：

```
b := [2]string{"Penn", "Teller"}
```

或者你可以靠编译器统计你数组的元素：

```
b := [...]string{"Penn", "Teller"}
```

上面两种都代表的是 `[2]string` 数组类型。

## Slices

数组有它们自己的应用场景，但是还是有点不灵活，所以你在 Go 代码中经常看不到这些数组代码。反而 Slices 到处都是。它构建在数组之上，更加有力量和方便的。

Slices 类型说明是 `[]T`，其中 `T` 是 Slice 中元素的类型。不像数组，**slice 不需要为其指定长度**。

slice 字面量可以像数组那样申明，只是你要省去元素计数：

```
letters := []string{"a", "b", "c", "d"}
```

slice 可以通过 Go 内置的 make 创建：

```
func make([]T, len, cap) []T
```

其中 T 表示创建的 slice 的元素类型。`make` 函数需要传递类型，长度以及可选项容量参数。当调用时，`make` 会**分配数组内存以及返回一个指向这个数组的 slice**。

```
var s []byte
s = make([]byte, 5, 5)
// s == []byte{0, 0, 0, 0}
```

当 cap 参数省略时，默认就会只当一个长度。下面是同一版本更简单的版本代码：

```
s := make([]byte, 5)
```

len 和 cap 参数可以通过内置的 `len()` 和 `cap()` 函数检索校验。

```
len(s) == 5
cap(s) == 5
```

下面两节来讨论 len 和 cap 的关系。

slice 的零值是 `nil`。其函数 `len()` 和 `cap()` 返回的是 0。

slice 也能将已有的 slice 对象以及数组对象通过”切片“构成 slice。切片是通过指定一个半开范围(half-open range)，两个下标用冒号隔开来完成的。例如，表达式`b[1:4]` 创建一个包含 b 中元素 1 到 3 的切片(结果切片的索引将是 0 到 2)。

```
b := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
// b[1:4] == []byte{'o', 'l', 'a'}, 与 b 共享相同的存储空间
```

切片表达式的开始与结束参数是可选的；默认情况是从 0 到 slice 的长度：

```
// b[:2] == []byte{'g', 'o'}
// b[2:] == []byte{'l', 'a', 'n', 'g'}
// b[:] == b
```

这里也有语法支持从给定的数组创建一个 slice：

```
x := [3]string{"Лайка", "Белка", "Стрелка"}
s := x[:] // slice s 与数组 x 共享一个引用（内存地址空间）
```

## Slice 内部细节

**slice 是数组段的描述符，它由指向数组的指针、段的长度和它的容量(段的最大长度)组成**。

![](https://go.dev/blog/slices-intro/slice-struct.png)

我们变量 s，通过 `make([]byte, 5)` 创建，其数据结构就像：

![](https://go.dev/blog/slices-intro/slice-1.png)

长度通过 slice 指向元素的数量。capacity 是指在底层数组的元素数量（从切片指针所指向的元素开始）。为了更加清晰的区分 len 与 cap 两者的区别，我们通过下面的例子来说明。

我们有一个 slice 变量 s，观察 slice 数据结构的变化以及在其底层数组的关系：

```
s := s[2:4]
```

![](https://go.dev/blog/slices-intro/slice-2.png)

切片操作并没有发生数据的拷贝。它创建了新的 slice 值，它指向的原始数组对象。这使得切片操作与数组索引操作一样高效。因此修改 slice 中的元素，就修改数组的原始元素：

```
d := []byte{'r', 'o', 'a', 'd'}
e := d[2:]
// e == []byte{'a', 'd'}
e[1] = 'm'
// e == []byte{'a', 'm'}
// d == []byte{'r', 'o', 'a', 'm'}
```

早先我们把 s 切成的长度比它的容量短的长度。我们可以通过再次切割来达到它的容量：

```
s = s[:cap(s)]
```

![](https://go.dev/blog/slices-intro/slice-3.png)

一个切片不能超过它的容量。如果要这么做就会发生一个运行时的 panic，只需要索引设置成 slice 和数组的边界之外就能发生错误。同样的，slice 不能重切片至 0 以下以访问数组元素。

## 增长切片（拷贝与追加函数）

**为了增加切片的容量，必须创建一个新的、更大的切片，并将原始切片的内容复制到其中**。这个技术在其它语言也是动态数组的实现。下面的例子是通过构建新的 slice 将原来的容量翻倍，拷贝 s 的元素到 t，然后分配切片值 t 到 s：

```
t := make([]byte, len(s), (cap(s)+1)*2) // +1 in case cap(s) == 0
for i := range s {
        t[i] = s[i]
}
s = t
```

内置的复制函数使这个常见操作的循环部分变得更容易。顾名思义，copy 将数据从源片复制到目标片。它返回复制的元素个数。

```
func copy(dst, src []T) int
```

copy 函数支持两个长度不一的切片（它只会拷贝元素数量小的）。另外，copy 能处理源和目标 slice 共享相同的底层数组，即能够正确处理重叠的 slice。

通过 copy，我们可以高效的使用下面代码代替上面的代码：

```
t := make([]byte, len(s), (cap(s)+1)*2)
copy(t, s)
s = t
```

一个通用的操作就是追加数据到 slice 的尾部。这个函数追加字节元素到字节类型的 slice 对象，如必要就会增长 slice 并返回已经更新的 slice 值：

```go
func AppendByte(slice []byte, data ...byte) []byte {
    m := len(slice)
    n := m + len(data)
    if n > cap(slice) { // if necessary, reallocate
        // allocate double what's needed, for future growth.
        newSlice := make([]byte, (n+1)*2)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0:n]
    copy(slice[m:n], data)
    return slice
}
```

使用方式：

```
p := []byte{2, 3, 5}
p = AppendByte(p, 7, 11, 13)
// p == []byte{2, 3, 5, 7, 11, 13}
```

`AppendByte` 函数非常有用，因为它们提供了对切片增长方式的完全控制。根据程序的特性，可能希望以更小或更大的块进行分配，或者对重新分配的大小设置一个上限。

但绝大多数程序都不需要完全控制，所以 Go 提供了内置了 `append` 函数来应付大多数目的；它的签名如下：

```
func append(s []T, x ...T) []T
```

append 函数将元素 x 附加到切片 s 的末尾，如果需要更大的容量，则增加切片。

```
a := make([]int, 1)
// a == []int{0}
a = append(a, 1, 2, 3)
// a == []int{0, 1, 2, 3}
```

将一个 slice 追加到另一个 slice，可以使用 `...` 操作符拓展参数集合

```
a := []string{"John", "Paul"}
b := []string{"George", "Ringo", "Pete"}
a = append(a, b...) // equivalent to "append(a, b[0], b[1], b[2])"
// a == []string{"John", "Paul", "George", "Ringo", "Pete"}
```

由于 slice 的零值的行为很象长度为 0 的 slice，你可以申明一个 slice 变量并在循环中追加元素：

```
// Filter returns a new slice holding only
// the elements of s that satisfy fn()
func Filter(s []int, fn func(int) bool) []int {
    var p []int // == nil
    for _, v := range s {
        if fn(v) {
            p = append(p, v)
        }
    }
    return p
}
```

## 潜在的"陷阱"

如前所述，对切片进行重新切片并不会生成底层数组的副本。整个数组将一直保存在内存中，直到它不再被引用。偶尔，这可能会导致程序在只需要一小部分数据时将所有数据保存在内存中。

例如，`FindDigits` 函数将一个文件加载到内存中，并在其中搜索第一组连续的数字，将它们作为一个新片返回。

```
var digitRegexp = regexp.MustCompile("[0-9]+")

func FindDigits(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    return digitRegexp.Find(b)
}
```

这段代码的行为与发布的一样，但是返回的 `[]byte` 指向包含整个文件的数组。因为切片引用原始数组，**所以只要切片被保存在垃圾回收器周围，就不能释放数组**；文件的很少有用字节将整个内容保存在内存中。

为了修复这种问题，我们可以在它返回之前只对有兴趣的内容进行拷贝：

```
func CopyDigits(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    b = digitRegexp.Find(b)
    c := make([]byte, len(b))
    copy(c, b)
    return c
}
```
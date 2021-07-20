# 不安全包—unsafe

unsafe 包主要是围绕 Go 类型进行一些类型安全的操作，那为什么又是不安全(unsafe)的呢？因为使用这个功能的时候，可能会存在不兼容的问题，存在不可移植的问题，并且不受 Go 1 兼容指南的保护。

不安全包里面总共就一个类型和三个方法：

```go
type ArbitraryType int
type Pointer *ArbitraryType

func Sizeof(x ArbitraryType) uintptr

func Offsetof(x ArbitraryType) uintptr

func Alignof(x ArbitraryType) uintptr
```

## ArbitraryType

这个只是一个标记类型，没有实际意义；它表示的是 Go 中的任意类型。

## Pointer

这个类型是 unsafe 中的核心对象， 它表示的一个到任意类型的指针。Pointer 这里有四种操作可用，而其它类型是不可用的：

- 任意类型的指针值可以转换为一个 Pointer。
- Pointer 可以转换为任意类型的指针值。
- 一个 unpointer 能转型成一个 Pointer。
- 一个 Pointer 能转换成一个 unpointer。

因此 Pointer 允许一个程序击穿整个类型系统，读写任意的内存。所以在使用它的时候要小心。

关于下面 Pointer 的使用模式是有效的。没有使用这些模式的代码现在很可能是无效的，将来也可能是无效的。即使是下面的有效模式也有一些重要的警告。

运行 “go vet” 能帮助我们在使用 Pointer 时提示我们是否符合上面这些模式，但是它并不能保证这些代码一定就是有效的。

### 模式一：*T1 转型成 指向 *T2 的 Pointer 指针

前提是 **T2 不大于 T1 以及两者共享相同的内存布局**，这个约定允许重新解析一个类型的数据到另一类型的数据。math.Float64bits 就是其中一个实现：

```go
func Float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}
```

上面就是由 float64 转换成 uint64 的例子，是符合模式1的：uint64 <= float64，并且能共享相同的内存布局。

### 模式二：Pointer 转换成 uintptr（但是不能转回 Pointer）

将一个指针转换为一个 uintptr 会产生一个整数，这个整数指向值的内存地址。uintptr 通常的用法就是打印它。

uintptr 转回 Pointer 通常是无效的。uintptr 它是一个数字而非一个引用。

Pointer 转换成 uintptr 会创建一个没有指针语义的整数值。一旦 uintptr 持有了一些对象的地址，如果这些对象移除的时候，垃圾回收器将不会更新 uintptr 的值。

### 模式三：使用算法将 Pointer 与 uintptr 互转

如果 p 指向一个已分配的对象，它可以通过转换为 uintptr，增加一个偏移量，然后转换回 Pointer 的方式进入该对象。

```go
p = unsafe.Pointer(uintptr(p) + offset)
```

这个模式大多都用在一个结构中访问字段或者在数组中访问元素：

```go
// 等价于 f := unsafe.Pointer(&s.f)
f := unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Offsetof(s.f))

// 等价于 e := unsafe.Pointer(&x[i])
e := unsafe.Pointer(uintptr(unsafe.Pointer(&x[0])) + i*unsafe.Sizeof(x[0]))
```

在这种模式下通过 uintptr 加减偏移量都是有效的。通常对于对齐，在指针上使用 &^ 同样也是有效的。

在所有的例子里，结果必须还要继续将指针指向原始已分配的对象。不像 C 将指针推进到其初始分配的末端是无效的：

```c
// 无效：指针最后分配在已分配空间之外
var s thing
end = unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Sizeof(s))

// 无效：指针最后分配在已分配空间之外
b := make([]byte, n)
end = unsafe.Pointer(uintptr(unsafe.Pointer(&b[0])) + uintptr(n))
```

注意，两个转换都必须要求出现在相同的表达式，它们之间只有中间的算术：

```go
// 无效: uintptr 在转换回 Pointer 之前无法存储在变量中
u := uintptr(p)
p = unsafe.Pointer(u + offset)
```

注意，point 必须要指向已分配的对象，所以它一定不能为 nil：

```go
// 无效: nil 指针转换
u := unsafe.Pointer(nil)
p := unsafe.Pointer(uintptr(u) + offset)
```

### 模式四：当调用 syscall.Syscall 时，Pointer 转换为 uintptr

当触发系统调用时会直接传递它们的 uintptr 参数给操作系统，然后根据调用的细节，操作系统可能会将其中一些重新解释为指针。

那也就是说，系统调用隐式实现了将这些参数（uintptr）从 uintptr 转换回了 pointer。

如果一个 pointer 参数必须转换为 uintptr 作为参数使用，那么这个转换必须出现在调用表达式本身中：

```go
syscall.Syscall(SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(p)), uintptr(n))
```

在程序集中的函数调用的参数列表中，通过安排引用已分配的对象，编译器就会处理这个 Pointer 转换为 uintptr。如果有这个对象，它会保留这个对象直到调用结束后才移走，即使单独从这个类型来看，这个对象在调用期间都不再需要了。

对于这个模式下的编译器重新组织，这个转换必须要出现在参数列表中：

```go
// 无效：在通过系统调用隐式转换回 Pointer 之前，uintptr 无法存储在变量中
u := uintptr(unsafe.Pointer(p))
syscall.Syscall(SYS_READ, uintptr(fd), u, uintptr(n))
```

### 模式五：reflect.Value.Pointer 转换结果或 reflect.Value.UnsafeAddr 从 uintptr 转换到 Pointer

reflect 包中的 Value 方法命名为 Pointer 和 UnsafeAddr，它们都返回 uintptr 而不是 unsafe.Pointer。Pointer 防止调用方在不首先导入“unsafe” 的情况下将结果更改为任意类型。但是这就意味着结果是易碎的(fragile)并且必须要在相同的表达式内调用之后立即转换为 Pointer：

```go
p := (*int)(unsafe.Pointer(reflect.ValueOf(new(int)).Pointer()))
```

在上面的例子，在转换之前将结果存到一个变量是无效的：

```go
// 无效: uintptr 在转换回 Pointer 之前无法存储在临时变量
u := reflect.ValueOf(new(int)).Pointer()
p := (*int)(unsafe.Pointer(u))
```

### 模式六：reflect.SliceHeader 或 reflect.StringHeader 数据字段与 Pointer 转换

在上节模式案例中，反射数据结构 SliceHeader 和 StringHeader 申明了 uintptr 字段。防止调用方在不首先导入 "unsafe" 的情况下将结果更改为任意类型。但是这就意味着 SliceHeader 和 StringHeader 仅在解释实际切片或字符串值的内容时有效。

```go
var s string
hdr := (*reflect.StringHeader)(unsafe.Pointer(&s)) // 模式 1
hdr.Data = uintptr(unsafe.Pointer(p))              // 模式 6 (this case)
hdr.Len = n
```

在这个用法中 hdr.Data 实际上是在字符串头文件中引用底层指针的另一种方法，而不是 uintptr 变量本身。

一般来说，reflect.SliceHeader 和 reflect.StringHeader 应该只作为 *reflect.SliceHeader 和 *reflect.StringHeader 指向实际的切片或字符串使用，而不是普通的结构体。

程序不应该为这些结构类型申明或分配变量。

```go
// 无效: 一个直接申明的 header 将不会持有数据的引用
var hdr reflect.StringHeader
hdr.Data = uintptr(unsafe.Pointer(p))
hdr.Len = n
s := *(*string)(unsafe.Pointer(&hdr)) // hdr 转换的 Point 可能已经丢失了
```

## unsafe.Sizeof(x ArbitraryType) uintptr

Sizeof 有一个任意类型参数，并返回假设变量 v 的大小(以字节为单位)，就像 v 是通过var v = x声明的一样。返回的 size 值不包括任何引用 x 的地址。

例如，如果 x 是一个切片，Sizeof 返回的就是 slice 描述的大小，而不是 slice 引用的内存大小。Sizeof 返回的 Go 常数。

## unsafe.OffsetOf(x ArbitraryType) uintptr

Offsetof 返回 x 表示的字段结构内的偏移量，该结构必须是 structValue.field 的形式。换句话说，它返回**结构开始和字段开始之间的字节数**。Offsetof 的返回值是一个 Go 常数。

## unsafe.Alignof(x ArbitraryType) uintptr

Alignof 传递一个任意类型参数 x，并返回假设变量 v 所需的对齐方式，就像 v 是通过 var v = x 申明的一样。

它是最大的值 m，使 v 的地址总是对 m 整除的。它与 reflect.TypeOf(x).Align() 返回的结果是一样的。特殊情况下，如果一个变量 s 是结构类型以及 f 是在这个结构下的一个字段，那么调用 Alignof(s.f) 将会返回这个结构类型下的字段所需的对齐。这个例子与调用 reflect.TypeOf(s.f).FieldAlign() 是一致的。Alignof 的返回值是一个 Go 常熟。


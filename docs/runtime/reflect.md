# Go 语言中的反射—Reflect

反射提供了程序在运行时动态获取任意类型的对象，并对其进行操作的能力。它是一种元编程。它的使用绝大多数都是针对 TypeOf、ValueOf 这两个方法返回的对象进行操作的。

- TypeOf：返回任意对象在运行期的反射对象类型（reflect.Type）
- ValueOf：返回任意对象在运行期的反射对象值（reflect.Value）

## 反射三定律

在操作反射对象的时候，Go 团队要求我们必须符合三个定律：

1. 从接口对象 `interface{}` 获取反射对象；
2. 从反射对象也可以获取接口对象 `interface{}`；
3. 要想修改反射对象值，这个值就必须是可修改的(settable)

因为 Go 是静态类型语言，所以在 Go 中所有的静态类型我们都可以通过反射获取两个对象；一个代表这个对象的具体值(type value)；另一个代表这个对象的类型描述(type description)。而这一对变量就可以从 `reflect.TypeOf` 和 `reflect.ValueOf` 这两个方法获取。前者返回 reflect.Type 对象，后者返回 reflect.Value 对象。

## reflect.Type

[reflect.Type](https://github.com/golang/go/blob/master/src/reflect/type.go#L38) 代表了 Go 语言中的所有类型，这个对象中没有暴露出任何属性，有的只是方法：

```go
type Type interface {
		...
    Method(int) Method
    MethodByName(string) (Method, bool)
    Name() string
    Size() uintptr
    String() string
    Kind() Kind
    Implements(u Type) bool
    Comparable() bool
    Elem() Type
    Field(i int) StructField
    FieldByNameFunc(match func(string) bool) (StructField, bool)
    ...
    common() *rtype
    uncommon() *uncommonType
}
```

Type 对象里面的方法有很多，我还省略了一些方法。其中 [rtype]([reflect.rtype](https://github.com/golang/go/blob/master/src/reflect/type.go#L311)) 是 Type 接口对象的默认实现，它实现了 Type 接口中描述的所有方法。请注意，方法中有 `Comparable()` 就说明此类型是可以进行比较的。

我们通过 [reflect.TypeOf(x)](https://github.com/golang/go/blob/master/src/reflect/type.go#L1413) 就可以得到反射类型对象了。而方法中的参数是 interface{} 类型的，这就是反射第一个定律：要想获取对象的反射对象，就必须先将其转换为 interface{} 类型。然后再将**接口对象地址**转换为头对象 emptyInterface 并最终返回 Type 对象。而这些方法操作的对象的元数据信息是存在 [reflect.rtype](https://github.com/golang/go/blob/master/src/reflect/type.go#L311) 这个私有类型下的：

```go
func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}

type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}
```

通过上面的方法我们可以调用 Method 方法开获取该对象的方法信息，通过调用 Field 方法获取字段信息。这两个消息都提供了按名称、索引等需求检索方法和字段信息。

这里的 Kind() 方法要注意一下，这个方法返回的是反射对象的**底层类型**。什么意思呢？且看下面代码：

```go
var i int32 = 1
itype := reflect.TypeOf(i)
fmt.Println(itype.Kind())

type myInt32 int32
var j myInt32 = 2
jtype := reflect.TypeOf(j)
fmt.Println(jtype.Kind())
```

int32 类型的变量 i 经过反射打印 Kind 出来是 int32；但是 myInt32 类型的变量 j 经过反射打印出来的 Kind 出来的类型同样是 int32。因为 myInt32 类型还是由底层类型封装的。

## reflect.Value

reflect.Value 是 Go 值的反射接口。还是用上面的例子表述：

```go
var i int32 = 1
itype := reflect.TypeOf(i)
ivalue := reflect.ValueOf(i)
fmt.Println("itype = ", itype)		// itype = int32
fmt.Println("ivalue = ", ivalue)	// ivalue = 1
```

[reflect.Value](https://github.com/golang/go/blob/master/src/reflect/value.go#L39) 是一个 struct 对象而不是接口，所以内部对 Value 对象实现了很多方法。有兴趣的可以翻源码去看看，总之获取了这个反射对象，就能对结构和数据进行操作了。

需要注意的是，reflect.Value 记录的不仅仅是值，同样还有值的类型信息，我们从源代码中就可以知道 Value 对象含有 rtype 对象，所以 reflect.Type 中的方法在 reflect.Value 同样也有。

犹如前面所说，我们可以通过将 Go 中任意的对象通过调用 reflect.TypeOf/ValueOf 转换为 interface{} 对象(第一定律)获取反射对象。那么当我们操作完反射对象之后又该如何转换回去呢？

这就是第二定律了，我们通过**反射对象先转回为 interface{} 对象，然后在返回具体的 Go 中的类型对象**：

```go
func (v Value) Interface() interface{}
```

具体例子如下：

```go
var i int32 = 1
ivalue := reflect.ValueOf(i)
ii := ivalue.Interface().(int32)
fmt.Println(ii)
```

第三定律就涉及到当我们获取反射对象，要进行修改时的问题了。它要求我们要修改的值必须要是可赋值的(settable)。

我们从一个例子出发：

```go
var i int32 = 1
ivalue := reflect.ValueOf(i)
ivalue.SetInt(2)
```

当我们执行上面的代码时，就会报 `panic: reflect: reflect.Value.SetInt using unaddressable value` 错误。提示的错误也明显，因为 ivalue 这个对象不是可寻址的(addressable)。其实我们也可以通过系统包提供的 `Value.CanSet` 方法来判断是否可以调用 Set 方法修改值：

```go
fmt.Println("settability of v:", ivalue.CanSet())	// settability of v: false
```

因为 Go 语言的类型都是值传递的，所以当调用 `ivalue := reflect.ValueOf(i)` 这个方法传递的参数其实是 i 这个变量的值拷贝。所以即使修改了 ivalue 的值，变量 i 也不会受到影响。所以我们在获取一个变量的反射对象并且要对其某些值进行修改时，我们要传这个变量的地址：

```go
var i int32 = 1
ivalue := reflect.ValueOf(&i)
fmt.Println(ivalue.CanSet())
```

我们经过上面的修改发现，输出的还是 false。这是为什么呢？

由于我们在获取反射对象的时候传递的是目标变量的地址，所以 ivalue 代表的是变量 i 的地址对象。我们当然不能修改栈上的地址，而是通过 Value.Elem 来获取地址对应的值，然后再进行修改操作：

```go
var i int32 = 1
ivalue := reflect.ValueOf(&i)

value := ivalue.Elem()
fmt.Println(value.CanSet())	// true
value.SetInt(2)	
fmt.Println(value)	// 2
fmt.Println(i)	// 2
```

我们从 [Value.Elem](https://github.com/golang/go/blob/master/src/reflect/value.go#L1146) 的实现就能看发现，其实就是通过判断 `ivalue.kind` 是否是指针进行转换获取的。

PS：有些东西还没想好怎么梳理，如反射函数调用，unsafe.Pointer 是如何使用等。// TODO

# 参考链接

- https://blog.golang.org/laws-of-reflection#TOC_9.
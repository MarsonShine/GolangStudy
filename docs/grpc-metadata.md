# RPC 客户端与服务端传值的利器——Metadata

原文地址：https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md

grpc 支持在客户端和服务端之间发送元数据（metadata）。这个文档展示了如何在 grpc-go 发送和接受数据。

## 背景

有四种 RPC 服务方法：

- [一元 RPC（Unary RPC）](https://grpc.io/docs/guides/concepts.html#unary-rpc)
- [服务流 RPC（Server Stream RPC）](https://grpc.io/docs/guides/concepts.html#server-streaming-rpc)
- [客户端流 RPC（Client Stream RPC）](https://grpc.io/docs/guides/concepts.html#client-streaming-rpc)
- [双向绑定流 RPC（Bidirectional Stream RPC）](https://grpc.io/docs/guides/concepts.html#bidirectional-streaming-rpc)

这里首先得知道什么是 Metadata

## 什么是 Metadata

Metadata 是 RPC 服务调用时的附带信息（如[授权验证](https://grpc.io/docs/guides/auth/)），这些信息是以键值对集合的形式存在的，其中键和值通常都是字符串类型，但是也可以是二进制数据。**Metadata 本身对 grpc 是不透明的——它必须要让客户端在调用服务端的时候提供相关的信息，反之亦然**。

访问 Metadata 是依赖各自语言的实现的

> 以上概念出自 https://grpc.io/docs/what-is-grpc/core-concepts/#metadata

## 构造 metadata

通过包 [metadata](https://godoc.org/google.golang.org/grpc/metadata) 来创建 metadata。里面的 MD 类型实际上是一个 string - string 集合的 map。

```go
type MD map[string][]string
```

Metadata 可以像 map 一样读取。要注意这个 map 的值是一个 string 数组，所以用户可以用单个键存储多个值

### 创建新的 Metadata

我们可以使用 New 方法创建 metadata：

```go
md := metadata.New(map[string][]string{"key1": "val1","key2": "val2"})
```

还有另一种方式，使用 `Pairs`。用相同的 key 可以将多个值追加到集合中去：

```go
md := metadata.Pairs(
	"key1", "val1",
	"key1", "val1-2",	// 这个 key1 相同，但是可以创建一个 []string{"val1", "val2"}
	"key2", "val2",
)
```

注意：这里面所有的键都会自动转换为小写形式，所以 "key1" 和 "Key1" 是同一个键，并且他们会将值追加到集合中去。`New` 和 `Pairs` 都是如此。

### Metadata 存储二进制数

在 metadata，键通常是字符串。但是有时候也可以是二进制数。为了可以在 metadata 中存储二进制数据，只需要在键值前面简单的加一个 "-bin" 后缀。在创建 metadata 时，带有 "-bin" 后缀的键都会将值转码。

```go
md := metadata.Pairs(
	"key", "string value",
	"key-bin", string([]byte{96, 102}),	// 这个二进制数在发送之前将被转码，并且转移后就会解码
)
```

## 从 context 检索 metadata

我们可以用方法 `FromIncomingContext` 在 context 检索 metadata：

```go
func (s *server) SomeRPC(ctx context.Context, in *pb.SomeRequest) (*pb.SomeResponse, err) {
	md, ok := metadata.FromInComingContext(ctx)
}
```

## 发送和检索 metadata —— 客户端

### 发送 metadata

这里有两种方式发送 metadata 给服务端。其中推荐的方式是使用方法 `AppendToOutgoingContext` 追加 kv 对到上下文中。这个可以与当前上下文的已经存在的 metadata 一起使用，也可以不使用。当这里没有之前的 metadata 时，就会创建；如果当前上下文已经存在 metadata，就会将 kv 对合并进去。

```go
// 用一些元数据创建新的 context
ctx := metadata.AppendToOutgoingContext(ctx, "k1","v1", "k2", "v2", "v3")
// 接着添加更多的元数据到上下文中（例如在拦截器中）
ctx := metadata.AppendToOutgoingContext(ctx, "k3", "v4")
// 使用 unary rpc
response, err := client.SomeRPC(ctx, someRequest)
// 或者使用流 rpc
response, err := client.SomeStreamingRPC(ctx)
```

另外，metadata 也可以使用 `NewOutgoingContext` 追加到上下文中。但是这个是把现有的 context 全部替换，所以必要小心保持已经存在的那些 metadata。并且这个要比使用 `AppendToOutgoingContext` 要慢。下面是使用例子：

```go
// 创建一些带有元数据的 metadata
md := metadata.Pairs("k1", "v1", "k1", "v2", "k2", "v3")
ctx := metadata.NewOutgoingContext(context.Background(), md)
// 接着添加更多的元数据到上下文中（例如在拦截器中）
send, _ := metadata.FromOutgoingContext(ctx)	// 这里是重新创建了一个 context 返回
newMD := metadata.Pairs("k3", "v3")
ctx = metadata.NewOutgoingContext(ctx, metadata.Join(send, newMD))
// 使用 unary rpc
response, err := client.SomeRPC(ctx, someRequest)
// 或者使用流 rpc
stream, err := client.SomeStreamingRPC(ctx)
```

### 检索 metadata

在客户端也能检索元数据，包括头和尾部。

### Unary Call

与 unary call 调用一起发送的 Header 和 Trailer 可以使用 CallOption 中的 Header 和 Trailer 函数检索：

```go
var header, trailer metadata.MD // variable to store header and trailer
r, err := client.SomeRPC(
    ctx,
    someRequest,
    grpc.Header(&header),    // 检索 header
    grpc.Trailer(&trailer),  // 检索 trailer
)

// do something with header and trailer
```

### Stream Call

对于流调用包括：

- 服务端流 RPC
- 客户端流 RPC
- 双向流RPC

都能使用 [ClientStream](https://godoc.org/google.golang.org/grpc#ClientStream) 接口的的方法 `Header` 和 `Trailer` 将返回的流中检索出来 Header 和 Trailer：

```go
stream, err := client.SomeStreamingRPC(ctx)

// retrieve header
header, err := stream.Header()

// retrieve trailer
trailer := stream.Trailer()
```

## 发送和接收 Metadata —— 服务端

### 检索 metadata


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

Metadata 是 RPC 服务调用时的附带信息（如[授权验证](https://grpc.io/docs/guides/auth/)），这些信息是以键值对集合的形式存在的，其中键和值通常都是字符串类型，但是也可以是二进制数据。Meta

> 以上概念出自 https://grpc.io/docs/what-is-grpc/core-concepts/#metadata


# 调度器——GMP 调度模型

Goroutine 调度器，它是负责在工作线程上分发准备运行的 goroutines。

首先在讲 GMP 调度模型之前，我们先了解为什么会有这个模型，之前的调度模型是什么样子的？为什么要改成现在的模式？

我们从当初的[Goroutine 调度设计文档](https://docs.google.com/document/d/1TTj4T2JO42uD5ID9e89oa0sLKhJYD0Y_kqxDv3I3XMw/edit#)得知之前采用了 GM 的调度模型，并且在高并发测试下性能不高。文中提到测试显示 Vtocc 服务器在 8 核机器上的CPU最高为70%，而文件显示 `rutime.futex()` 就消耗了14%。通常，在性能至关重要的情况下，调度器可能会禁止用户使用惯用的细粒度并发。

那么是什么原因导致这些问题呢？Dmitry Vyukov 总结四个原因：

- 使用了一个全局互斥锁 mutex 处理整个与 goroutine 相关的操作（创建，完成，再调度等）。
- 频繁的 Goroutine 切换。工作线程会在那些可运行的 goroutine 之间频繁切换，这就导致了增加延迟以及额外的开销。
- 每个线程M都需要处理内存缓存（每个M的缓存与运行 G 所需要的缓存比例差距太大，100:1），这就导致了大量的内存占用影响了数据局部性。
- 系统调用(syscall)会导致工作线程频繁阻塞以及解除阻塞，这会导致大量的开销。

为了解决这个问题，于是就引入了 Processor 这个概念。









现在的调度模型主要分为三个概念：

- Goroutine(G)，表示待执行的任务
- 工作线程(M)，表示操作系统线程
- 处理器(P)，执行 Go 代码所需要的一种资源

P 必须要绑定到 M 上来执行具体的 Go 代码。





## 参考链接

- https://docs.google.com/document/d/1TTj4T2JO42uD5ID9e89oa0sLKhJYD0Y_kqxDv3I3XMw/edit#
- https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-goroutine/
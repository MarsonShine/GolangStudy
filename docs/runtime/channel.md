# Channel —— CSP

Go channel 并发原理设计的原型就是 CSP，即通信顺序进程（Communicating sequential processes，CSP）。

其核心思想就是：**不要通过共享内存的方式进行通信，而是应该通过通信的方式共享内存**。

![](https://img.draveness.me/2020-01-28-15802171487080-channel-and-goroutines.png)

​															                           （图出自：[https://draveness.me/golang/](https://img.draveness.me/2020-01-28-15802171487080-channel-and-goroutines.png)）

Goroutine 和 Channel 分别对应 CSP 中的实体和传递信息的媒介，Goroutine 之间会通过 Channel 传递数据。

上图两个 goroutine，一个会向 channel 发送数据，另一个会从 channel 接收数据，它们两个能够独立运行且不存在直接关系，而是通过 channel 间接关联起来的。

## 内部结构

Go 内部初始化 channel 的时候运行时会创建一个 [runtime.hchan](https://github.com/golang/go/blob/master/src/runtime/chan.go#L33) 的结构体。

```
type hchan struct {
	qcount   uint           // 在队列中的所有数据
	dataqsiz uint           // 环形队列的大小
	buf      unsafe.Pointer // 指向由 dataqsiz 个元素组成的数组的指针
	elemsize uint16
	closed   uint32	// 表示是否关闭通道的标识
	elemtype *_type // 元素类型
	sendx    uint   // 发送信息的当前位置索引
	recvx    uint   // 接收信息索引
	recvq    waitq  // 等待接收队列的集合
	sendq    waitq  // 等待发送的队列集合

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}
```

在实现的过程中有一些不变性需要了解一下：

c.sendq 和 c.recvq 之间至少有一个为空，除了使用 select 语句发送和接收时阻塞单个 goroutine 的非缓冲通道，在这种情况下，c.sendq 和 c.recvq 的长度只受 select 语句的大小限制。 

对于缓冲通道，我们有：

- c.qcount > 0 隐式代表 c.recvq 为空
- c.qcount < c.datasiz 隐式表示 c.sendq 为空

虽然我们通过 CSP 的方式来提高并发的性能，但并不意味着 chan 实现的是无锁队列。从结构体中的 `lock` 字段可以看出内部还是使用锁。

在 chan 中，**锁的目的是为了保护 hchan 所有的字段信息，以及在该通道上阻塞的 sudogs 中的几个字段**。

当占据锁的时候不要更改其它 g 的状态（尤其是不要准备另一个 G），因为这会导致堆栈收缩而死锁。

在 hchan 结构里面还有一个 waitq 队列结构：

```
type waitq struct {
	first *sudog
	last  *sudog
}
```

sendq 与 recvq 这两个循环队列都是因为缓冲区满了而进入队列等待的集合。

## 创建 channel

go 内置了 make 函数创建 channel，内部经由 runtime 转换调用 [runtime.makechan](https://github.com/golang/go/blob/master/src/runtime/chan.go#L72) 函数。makechan 函数会检查元素类型，并根据元素创建的元素类型和缓冲区的大小来具体创建不同的 channel：

```go
func makechan(t *chantype, size int) *hchan {
	elem := t.elem
	//...
	checkElemType...
	//...
	var c *hchan
	switch {
	case mem == 0:
		// Queue or element size is zero.
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		// Race detector uses this location for synchronization.
		c.buf = c.raceaddr()
	case elem.ptrdata == 0:
		// Elements do not contain pointers.
		// Allocate hchan and buf in one call.
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		// Elements contain pointers.
		c = new(hchan)
		c.buf = mallocgc(mem, elem, true)
	}
	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)
	lockInit(&c.lock, lockRankHchan)
	...
	return c
}
```

- 如果当前创建的 channel 不存在缓冲区的话，那么就只会为 runtime.hchan 分配一段内存空间；
- 如果当前创建的 channel 不包含指针，那么就会为 runtime.hchan 以及缓冲区 buf 分配一段连续的内存；
- 如果包含指针，就会单独为 runtime.hchan 和缓冲区分配内存；

## 发送信息

发送消息是通过调用 [runtime.chansend](https://github.com/golang/go/blob/master/src/runtime/chan.go#L159) 实现的。通常单通道的 send 和 recv，如果 block 不为 nil，那么这些协议不会休眠，如果没有立即完成则会直接返回。

发送信息的总体思路：

当一个通道中还有休眠线程待处理时，如果这时通道关闭，休眠的可以通过 `g.param == nil` 判断来唤醒。循环和重新运行操作是最容易的；我么后续可以看到它是如何关闭的。

在通道发送消息时，会经过一系列的状态判断：

快速路径：在没有获取锁的情况下检查非阻塞的失败操作。

之后在观察到通道没有关闭，我们观察到通道还没有准备好发送。每一个观察值都是一个单词大小的读数(第一个c.closed，第二个full())。因为一个关闭的通道不能从“准备发送”转变到“未准备发送”，即使这个通道在两个观察者之间关闭了，它们也意味着两个观察者之间存在一个时刻，既通道既没有关闭，也没有准备发送。我们如果观察到这个通道存在这个时刻，我们就会报告这个发送操作不能继续。

如果读取在这里被重新排序是可以的：如果我们观察到通道未准备发送并没有关闭，这意味着在第一次观察期间是无法关闭的。但是没有任何东西能保证一直继续往前。我们依靠 chanrecv() 和 closechan() 中释放锁的副作用来更新 c.closed 和 full() 的线程视图。

接下来看具体是如何发送信息：

```go
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	...
	if !block && c.closed == 0 && full(c) {
		return false
	}
	...
	lock(&c.lock)

	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}
	... 
    // 直接发送消息
	if sg := c.recvq.dequeue(); sg != nil {
		// Found a waiting receiver. We pass the value we want to send
		// directly to the receiver, bypassing the channel buffer (if any).
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}
	...
	if !block {
		unlock(&c.lock)
		return false
	}
	...
	return true
}
```

在发送之前先检查 channel 是否关闭，注意检查状态之前必须要先上锁，防止其它程序更改状态字段。如果通道关闭则释放锁并 panic。

如果没有关闭，则首先会去接收者等待队列中获取 g，如果存在则传递数据直接发送给向接收者。

如果缓冲区中还有可用空间，就把发送的元素放至队列中。

如果满了就会等待接收者接收数据。

### 直接发送  sendDirect

前面说了如果存在等待的接收者时，我们是直接向发送数据的；具体是通过 [runtime.send](https://github.com/golang/go/blob/master/src/runtime/chan.go#L293) 实现的。

```go
func send(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
    ...
	if sg.elem != nil {
		sendDirect(c.elemtype, sg, ep)
		sg.elem = nil
	}
	gp := sg.g
	unlockf()
	gp.param = unsafe.Pointer(sg)
	sg.success = true
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}
	goready(gp, skip+1)
}

func sendDirect(t *_type, sg *sudog, src unsafe.Pointer) {
	// src is on our stack, dst is a slot on another stack.

    // 一旦我们从 sg 获取 elem 元素，如果目标堆栈被拷贝（缩小），它就不在更新
    // 所以我们要确保在读和使用之间不会发生抢占点
	dst := sg.elem
	typeBitsBulkBarrier(t, uintptr(dst), uintptr(src), t.size)
    // 不需要 cgo 写屏障检查，因为 dst 总是在 go 内存中。
	memmove(dst, src, t.size)
}
```

在未缓冲或空缓冲的通道上发送(send)和接收(recv)是一个正在运行的 goroutine 向另一个正在运行的 goroutine 的堆栈写入数据的唯一操作。GC 假设堆栈写操作只在 goroutine 运行时发生，并且只由 goroutine 完成。使用写屏障(write barrier)就足以弥补违反这个假设，但是写屏障必须要工作。`typedmemmove` 将调用`bulkBarrierPreWrite`，但目标字节不在堆中，所以这不会有帮助。所以我们调用 `memmove` 和 `typeBitsBulkBarrier`。

然后再给等待接收数据的 g 通过 cas 设置标记成可运行的状态，然后将该 g 压入 runq 可运行队列中等待下一轮调度时立刻唤醒数据的接收方。

### 缓冲区发送数据

```go
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	...
	// 缓冲区
	if c.qcount < c.dataqsiz {
		// Space is available in the channel buffer. Enqueue the element to send.
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			racenotify(c, c.sendx, nil)
		}
		typedmemmove(c.elemtype, qp, ep)
		c.sendx++
		// 发送数据的位置索引 = 数据队列的长度时
		// 就说明索引已经到环形数组的队尾，就意味着要重新开始
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		c.qcount++
		unlock(&c.lock)
		return true
	}
	...
}
```

如果缓冲通道还有空间时，就会从 sendx 位置获取下一个元素，并通过 `typedmemmove` 将发送的数据拷贝到缓冲区中，并增加 sendx 索引值以及计数器 qcount 的值。

如果当 Channel 没有接收者能够处理数据时，就会产生阻塞，如果选择的是非阻塞，则会立即返回。否则会阻塞实现下面的代码逻辑。

```go
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	if !block {
		unlock(&c.lock)
		return false
	}

	// Block on the channel. Some receiver will complete our operation for us.
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// No stack splits between assigning elem and enqueuing mysg
	// on gp.waiting where copystack can find it.
	mysg.elem = ep
	mysg.waitlink = nil
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.waiting = mysg
	gp.param = nil
	c.sendq.enqueue(mysg)
	// Signal to anyone trying to shrink our stack that we're about
	// to park on a channel. The window between when this G's status
	// changes and when we set gp.activeStackChans is not safe for
	// stack shrinking.
	atomic.Store8(&gp.parkingOnChan, 1)
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceEvGoBlockSend, 2)
	// Ensure the value being sent is kept alive until the
	// receiver copies it out. The sudog has a pointer to the
	// stack object, but sudogs aren't considered as roots of the
	// stack tracer.
	KeepAlive(ep)

	// someone woke us up.
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	gp.activeStackChans = false
	closed := !mysg.success
	gp.param = nil
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	mysg.c = nil
	releaseSudog(mysg)
	if closed {
		if c.closed == 0 {
			throw("chansend: spurious wakeup")
		}
		panic(plainError("send on closed channel"))
	}
	return true
}
```

阻塞当前通道，并获取当前的 g 以及通过信号量获取 sudog 对象，并为 sudog 设置上下文信息（如数据、当前 g、当前 g 是否在 select 中的标识以及当前阻塞的 channel）。然后我们可以通过给当前的 `g.waiting = sudog` 就能获取当前等待的信息。最后将 sudog 对象入队列到等待发送队列中。

接着通过原子操作 `atomic.Store8(&gp.parkingOnChan, 1)` 通知那些所有想缩小堆栈的人，我们已经在通道上休眠了（park）。并通过 `gopark` 函数将当前的 g 状态设置为等待状态等待唤醒。

唤醒之后会执行一些清理工作，将 sudog 对象的一些属性手动设置为 nil 释放空间。

最后返回 true 表示成功向 Channel 发送数据。

### 总结

Channel 发送信息总共有三个分支：

- 如果当前的 Channel 已存在等待的接收者，那么就会选择直接将数据发送给这个接收者并设置下一个等待执行的接收者 `c.recvq.dequeue()`
- 如果 Channel 的缓冲区没满，则会阻塞当前 Goroutine，并将数据发送给缓冲区等待其它接收者 Goroutine 执行，并记录下当前的位置以及数据计数器 `c.sendx,c.qcount`
- 如果都不满上面的分支，则会阻塞当前 Goroutine，并创建新的 sudog 对象并设置当前 Goroutine 相关的上下文信息发送至发送者等待队列中，等待其他 Goroutine 接收数据

## 接收数据

接收数据 Go 内部是通过运行时调用 [runtime.chanrecv](https://github.com/golang/go/blob/master/src/runtime/chan.go#L455) 实现的，它返回两个 bool 类型的参数。其主要的设计思路为：

chanrecv 从 Channel 上接收数据并写入数据至 ep 指针（对象变量）。ep 可能为空，在这种情况下接收的数据可以忽略。这里有三种情况：

- 如果是非阻塞情况下（`block == false`）以及没有可用的元素，则直接返回 false,false。
- 如果 Channel 被关闭了，指针 ep 就会置零并返回 true,false。
- 如果上面两种情况都不满足，会选择将数据元素填至 ep 指针并返回 true,true

一个非 null 的 ep 指针必须是指向堆或者是调用者的栈地址。

runtime.chanrecv 的实现与 runtime.chansend 的实现一致：

- 在一个空 Channel 接收数据时会直接调用 gopark 休眠让出处理器的使用权；
- 如果 Channel 已经关闭并且缓冲区不存在任何数据时，则会清楚 ep 指针的数据并立即返回；

除去上面这两个特殊的情况，同样的，接收数据逻辑与发送数据类似：

- 当检测到发送者等待队列中还有空闲的发送者，就会直接从发送者接收消息。
- 如果 Channel 缓冲区还有数据时，从 Channel 的缓冲区获取接收数据。
- 当 Channel 缓冲区不存在数据时，等待其它 Goroutine 向 Channel 发送数据。

### 直接接收数据

直接接收数据的逻辑与发送信息的逻辑几乎一样：

```go
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
	...
	if sg := c.sendq.dequeue(); sg != nil {
		// Found a waiting sender. If buffer is size 0, receive value
		// directly from sender. Otherwise, receive from head of queue
		// and add sender's value to the tail of the queue (both map to
		// the same buffer slot because the queue is full).
		recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true, true
	}
	...

}
```

如果发现 Channel 中的发送等待着队列中还有 Goroutine，则直接取出来并通过 [runtime.recv](https://github.com/golang/go/blob/master/src/runtime/chan.go#L608) 发送数据。

```
func recv(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
		if c.dataqsiz == 0 {
		if raceenabled {
			racesync(c, sg)
		}
		if ep != nil {
			// copy data from sender
			recvDirect(c.elemtype, sg, ep)
		}
	} else {	
		qp := chanbuf(c, c.recvx)
		if raceenabled {
			racenotify(c, c.recvx, nil)
			racenotify(c, c.recvx, sg)
		}
		// copy data from queue to receiver
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
		// copy data from sender to queue
		typedmemmove(c.elemtype, qp, sg.elem)
		c.recvx++
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		c.sendx = c.recvx // c.sendx = (c.sendx+1) % c.dataqsiz
	}
	sg.elem = nil
	gp := sg.g
	unlockf()
	gp.param = unsafe.Pointer(sg)
	sg.success = true
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}
	goready(gp, skip+1)
}
```

recv 在满 Channel 上处理接收操作，包括两部分：

1. 发送方 sg 发送的值被放入 Channel，并且发送方被唤醒，来继续往后执行。
2. 当前 G 的接收者将接收值写入 ep

对于同步通道操作，这两个值时相同的。

对于异步通道操作，接收方从通道缓冲区获取数据，发送方的数据放在通道缓冲区中。

[runtime.recv](https://github.com/golang/go/blob/master/src/runtime/chan.go#L608) 函数做了两件事：

1. 判断当前 Channel 的缓冲大小如果为 0，即不存在缓冲区，则直接通过 [runtime.recvDirect](https://github.com/golang/go/blob/181e8cde301cd8205489e746334174fee7290c9b/src/runtime/chan.go#L347) 直接从发送者接收消息；具体通过 `typeBitsBulkBarrier` 和 `memmove` 将发送者的数据拷贝到接收者内存地址中。
2. 如果 Channel 存在缓冲区，那么就进而分两个部分：
   1. 通过 `typedmemmove` 将队列中的数据拷贝至接收者内存中
   2. 将发送队列头的数据拷贝到缓冲区中，并更新接收者的位置（同样，环形队列长度与接收位置索引相等，则表示从头开始）

最后就会调用 [runtime.goready](https://github.com/golang/go/blob/master/src/runtime/proc.go#L375) 将当前发送者等待处理程序的 g 将状态从 waiting 状态变更为可运行状态；并入队列到当前处理器（_p_）的 run next（`_g_m.p.ptr().runnext`），等待下一轮调度器将阻塞的发送方唤醒。

### 缓冲区

如果 Channel 缓冲区中存在数据，接收数据操作会直接从 Channel 缓冲区提取 recvx 位置的数据进行后续处理：

```go
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
	...
	if c.qcount > 0 {
		// Receive directly from queue
		qp := chanbuf(c, c.recvx)
		if raceenabled {
			racenotify(c, c.recvx, nil)
		}
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
		typedmemclr(c.elemtype, qp)
		c.recvx++
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		c.qcount--
		unlock(&c.lock)
		return true, true
	}
	...
}
```

最后就会调用 [runtime.typedmemclr](https://github.com/golang/go/blob/181e8cde301cd8205489e746334174fee7290c9b/src/runtime/mbarrier.go#L306) 清理队列中的数据类型，即更新计数器和索引以及释放 Channel 锁。

接下来就是除去上面两种情况，必须要阻塞接收数据

### 阻塞接收

当 Channel 不存在等待的发送者队列处理程序，Channel 的缓冲区也不存在任何数据时，从 Channel 接收数据就会被阻塞；当然也可以通过 block 与 select 语句控制不阻塞：

```go
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
  ...
	if !block {
		unlock(&c.lock)
		return false, false
	}
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	mysg.elem = ep
	mysg.waitlink = nil
	gp.waiting = mysg
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.param = nil
	c.recvq.enqueue(mysg)
	atomic.Store8(&gp.parkingOnChan, 1)
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanReceive, traceEvGoBlockRecv, 2)
  
  if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	gp.activeStackChans = false
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	success := mysg.success
	gp.param = nil
	mysg.c = nil
	releaseSudog(mysg)
	return true, success
}
```

在没有任何发送者可用的情况下，创建新的 sudog，并将接收内存地址与当前 goroutine 的上下文保留，并将其插入接收者等待队列中去。入队列之后就会通知其它程序该 g 进入休眠状态，让出处理器的使用权并等待下次调度。

最后进行一些清除扫尾工作，释放锁、内存等。

### 总结

发送数据分五种情况：

1. 在一个空 Channel 接收数据时会直接调用 gopark 休眠让出处理器的使用权；
2. 如果 Channel 已经关闭并且缓冲区不存在任何数据时，则会清除 ep 指针的数据并立即返回；
3. 当检测到发送者等待队列中还有空闲的发送者（goroutine），如果不存在缓冲区就会直接从发送者接收消息。存在缓冲区就会将发送等待者的数据拷贝到缓冲区的 recvx 位置接收者的空间地址；
4. 如果 Channel 缓冲区中存在数据，接收数据操作会直接从 Channel 缓冲区提取 recvx 位置的数据
5. 一般情况会阻塞当前的 Goroutine，新建 sudog 结构压入接收着等待者队列中等待下次调度器唤醒。 

## 关闭 Channel

关闭 Channel 是通过 [runtime.closechan](https://github.com/golang/go/blob/4847c47cb8a93b56e1df8c249700e25f527d4ba3/src/runtime/chan.go#L356) 实现的，主要逻辑就是清理那些还在队列里等待被唤醒的接收者与发送者。因为要更改 Channel 的属性以及队列里面的数据所以必须要加锁保证其它线程更改数据。

在具体执行关闭操作之前，先上锁并设置了关闭标识字段 closed = 1；

具体的关闭操作分为三部分：

1. 释放 Channel 中所有等待唤醒的接收者 Goroutine
2. 释放 Channel 中所有等待唤醒的发送者 Goroutine
3. 将所有等待的 Goroutine 压入 glist 列表中，并将所有 Goroutine 状态变更为可运行状态放至到 runq 队列中的 runnext 等待下一次调度器调度执行。

```go
func closechan(c *hchan) {
	...
	lock(&c.lock)
	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("close of closed channel"))
	}
	...
	c.closed = 1
	var glist gList
	// release all readers
	for {
		sg := c.recvq.dequeue()
		if sg == nil {
			break
		}
		if sg.elem != nil {
			typedmemclr(c.elemtype, sg.elem)
			sg.elem = nil
		}
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
		gp := sg.g
		gp.param = unsafe.Pointer(sg)
		sg.success = false
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
		glist.push(gp)
	}

	// release all writers (they will panic)
	for {
		sg := c.sendq.dequeue()
		if sg == nil {
			break
		}
		sg.elem = nil
		if sg.releasetime != 0 {
			sg.releasetime = cputicks()
		}
		gp := sg.g
		gp.param = unsafe.Pointer(sg)
		sg.success = false
		if raceenabled {
			raceacquireg(gp, c.raceaddr())
		}
		glist.push(gp)
	}
	unlock(&c.lock)

	// Ready all Gs now that we've dropped the channel lock.
	for !glist.empty() {
		gp := glist.pop()
		gp.schedlink = 0
		goready(gp, 3)
	}
}
```

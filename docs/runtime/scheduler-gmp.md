# 调度器——GMP 调度模型

Goroutine 调度器，它是负责在工作线程上分发准备运行的 goroutines。

首先在讲 GMP 调度模型之前，我们先了解为什么会有这个模型，之前的调度模型是什么样子的？为什么要改成现在的模式？

我们从当初的[Goroutine 调度设计文档](https://docs.google.com/document/d/1TTj4T2JO42uD5ID9e89oa0sLKhJYD0Y_kqxDv3I3XMw/edit#)得知之前采用了 GM 的调度模型，并且在高并发测试下性能不高。文中提到测试显示 Vtocc 服务器在 8 核机器上的CPU最高为70%，而文件显示 `rutime.futex()` 就消耗了14%。通常，在性能至关重要的情况下，调度器可能会禁止用户使用惯用的细粒度并发。

那么是什么原因导致这些问题呢？Dmitry Vyukov 总结四个原因：

- 使用了一个全局互斥锁 mutex 处理整个与 goroutine 相关的操作（创建，完成，再调度等）。
- 频繁的 Goroutine 切换。工作线程会在那些可运行的 goroutine 之间频繁切换，这就导致了增加延迟以及额外的开销。
- 每个线程M都需要处理内存缓存（每个M的缓存与运行 G 所需要的缓存比例差距太大，100:1），这就导致了大量的内存占用影响了数据局部性。
- 系统调用(syscall)会导致工作线程频繁阻塞以及解除阻塞，这会导致大量的开销。

为了解决这个问题，于是就引入了 Processor 这个概念。引入了这个对象并不会因为多了一个对象开销性能都会有影响，反而这方面开销都下降了。P 其实负责的是 M 与 G 之间的调度相关的操作，在执行 G 时 P 一定要与 M 绑定。并且把 M，schedule 里面的对象都转移到 P 中去了，所以 M 与 调度器原来的操作反而变得更干净了。如[调度设计文档](https://docs.google.com/document/d/1TTj4T2JO42uD5ID9e89oa0sLKhJYD0Y_kqxDv3I3XMw/edit#)中提到的：当 M 准备执行 Go 代码时会从集合表中弹出一个 P；当执行代码结束后就会将 P 推进集合中。所以当 M 需要执行 Go 代码时，必须要与 P 绑定。而新增的这个机制，就是为了替代原来调度器中的 sched.atomic(mcpu/mcpumax)。

现在的调度模型主要分为三个概念：

- Goroutine(G)，表示待执行的任务
- 工作线程(M)，表示操作系统线程
- 处理器(P)，执行 Go 代码所需要的一种资源

P 必须要绑定到 M 上来执行具体的 Go 代码。



在讲 GMP 调度模型之前，我们先来了解以下 G、M、P 这三个对象有哪些核心变量。

## G

Goroutine 是建立在 M 内核线程之上的称为协程的一个执行单元。在切换 G 时都是直接在用户态发生的，所以开销很小。所占用的内存也比原来小了很多，从前面的内容我们知道，我们把其中某些元素放入至新引入的 P 中了。虽然占用的内存不大，但是里面的变量却非常多。我们目前了解其中相对重要的部分，其它的字段想进一步了解，可以直接查看 [runtime2.g 源码](https://github.com/golang/go/blob/master/src/runtime/runtime2.go#L403).

```
type g struct {
	stack       stack   // offset known to runtime/cgo
	stackguard0 uintptr // offset known to liblink
	stackguard1 uintptr // offset known to liblink
	...
}
```

- stack 开头的三个变量，都是与栈相关的变量。stack 表示当前 g 所占用的实际栈内存大小：[stack.lo, stack.hi]。
- stackguard0 这个是在 Go 栈内存增长，参与栈内存计算比较的时候要用到的，并且会**参与抢占调度**。
- stackguard1 这个是在 C 栈内存增长时，参与栈内存计算比较。

```
type g struct {
	...
	m            *m
	sched        gobuf
	param        unsafe.Pointer
	atomicstatus uint32
	schedlink    guintptr
	gopc           uintptr
	startpc        uintptr
	waiting        *sudog
	...
}
```

- g 中包含 m 这个不意外，因为 g 的执行最终是通过当前的 m 来执行的。
- sched 这是一个 [gobuf](https://github.com/golang/go/blob/aa4e0f528e1e018e2847decb549cfc5ac07ecf20/src/runtime/runtime2.go#L312) 对象，里面包含栈指针 sp、程序计数器 pc、返回值 ret以及其它上下文信息 ctxt 等，在 g 发生调度的时候（如系统调用）就会靠着 shced 来执行或恢复之前操作相关的数据
- param g 在活动期间传递的参数
- atomicstatus 这个就表示当前 g 的运行状态了
- schedlink，这个有点意思，它表示是调度器的链接器。通过它 g 就可以在 [schedt](https://github.com/golang/go/blob/aa4e0f528e1e018e2847decb549cfc5ac07ecf20/src/runtime/runtime2.go#L751) 调度器的全局的 runq 队列中定位到可用的 g。
- gopc 是开启 goroutine 那个时刻的程序计数器（就是 Go 代码中的 go func() 的位置）
- startpc 即 go 开启协程执行的函数的程序计数器 pc
- waiting 表示 go 在等待，是一个 sudog 指针类型，sudog 表示的是一个等待 g 的集合列表。还有等待相关的参数如（waitsince、waitreson）这里就不介绍了。

```
type g struct {
	...
	preempt       bool // 抢占信号, 与 stackguard0 = stackpreempt 一样
	preemptStop   bool // 抢占状态更改为 _Gpreempted
	preemptShrink bool // shrink stack at synchronous safe point
	...
}
```

这三个变量是跟抢占调度相关的。

## M

M 是指操作系统线程，Go 在启动时会根据 CPU 的核心数分配 M 的个数。最多会开启 10000 个线程，并且这里面大多数都不会执行用户代码。最多只有 GOMAXPROCS 个活跃的线程执行用户代码。默认的设置一般都是 CPU 的核心数，这样是为了在调度的时候防止线程频繁的发生上下文切换。而在调度 G 的所有过程都是在用户态进行的，较于操作系统级的线程 M 切换来说开销会小的多。

我们来看一下 M 的主要核心对象：

```
type m struct {
    g0      *g
    curg          *g
    ...
}
```

g0 是个特殊的 goroutine，它是持有调度栈的，它会参与调度的过程。如创建 m，创建 g 以及执行一些内存分配。

```
type m struct {
    p             puintptr // attached p for executing go code (nil if not executing go code)
	nextp         puintptr
	oldp          puintptr // the p that was attached before executing a syscall
    ...
}
```

这三个字段是与 P 处理器相关的；

- p，当前正在执行 go 代码的 P，如果没有执行代码就为 nil
- nextp，下一个要执行的 p
- oldp，在执行系统调用之前的 p

## P

是新引入的在 G 和 M 之间的调度层。它负责调度 runq 等待队列中的待运行的协程，在关键的操作时候可以选择让出线程，提高线程利用率。

P 内部也包含了大量对象，同样我们主要了解其中相对重要的字段，与调度那些等待的 g 密切相关的内容。

```
type p struct {
	m           muintptr
	// Queue of runnable goroutines. Accessed without lock.
	runqhead uint32
	runqtail uint32
	runq     [256]guintptr
	runnext guintptr
	...
}
```

前面提到了在执行具体 go 代码时，p 一定要与 m 相关联。后续的字段都是与运行的 goroutine 相关。

- runq，是一个长度固定为 256 的数组结构
- runqhead，runqtail 表示的当前 runq 队列中的首尾的位置
- runnext，表示的下一个可运行的 g（注意，不一定表示一定会在下一轮唤醒运行，如果没有弹出执行就又会回到队列的头部位置）

其实从这四个字段我们就能看出，runq 其实一个**由数组结构加双指针构成的一个环形队列结构**。

除此之外，还有一个匿名结构类型的字段需要注意，那就是 gFree。这个对象内部是由 gList、n 组成的一个链表对象，用来存放空闲 g 的。

```
gFree struct {
    gList
    n int32
}
```

## 启动 Schedule 调度器

在调度 GMP 之前我们必须还要知道调度器是如何启动的。

调度器启动在 [runtime.schedinit](https://github.com/golang/go/blob/aa4e0f528e1e018e2847decb549cfc5ac07ecf20/src/runtime/proc.go#L5947) 可以看得到。除去初始化锁的顺序信息和其它必要的信息（如gc、栈、系统参数与环境变量等），我们主要看下面几个变量：

```go
func schedinit() {
	...
	_g_ := getg()
	if raceenabled {
		_g_.racectx, raceprocctx0 = raceinit()
	}
	sched.maxmcount = 10000
	...
	lock(&sched.lock)
	sched.lastpoll = uint64(nanotime())
	procs := ncpu
	if n, ok := atoi32(gogetenv("GOMAXPROCS")); ok && n > 0 {
		procs = n
	}
	if procresize(procs) != nil {
		throw("unknown runnable goroutine during bootstrap")
	}
	unlock(&sched.lock)
	...
}
```

调度器最多只能开启 10000 个线程。如果设置了 GOMAXPROCS 则替换默认的 cpu 核心数。之后就会调用 procresize 对 proc 进行更改。这个时候调度器必须要上锁，不会执行任何 goroutine 代码。procresize 函数内部对全局变量 `allp` 的期望容量 `capcity` 与 procs 进行判断。如果目标值要比期望值大，则会进行扩容给。否则直接追加即可：

```go
func procresize(nprocs int32) *p {
	if nprocs > int32(len(allp)) {
		lock(&allpLock)
		if nprocs <= int32(cap(allp)) {
			allp = allp[:nprocs]
		} else {
			nallp := make([]*p, nprocs)
			copy(nallp, allp[:cap(allp)])
			allp = nallp
		}
		...
		unlock(&allpLock)
	}
	...
}
```

扩容之后就会循环初始化 p（初始化期间对 pp 的id、status 以及内存缓存赋值）， 并调用底层系统的写屏障(write barrier)确保安全的对 allp 进行覆盖。

在初始化阶段，p 的状态此时是 _Pgcstop。在初始化之后如果当前的 p 的序号是小于之前设置的 nproc 目标数时，就会将当前的 g.m.p 的状态更改为 _Prunning。如果不满足上述条件，则会恒定取全局的 allp 中的第一个，并将状态设置为 _Pidle。

设置完当前的 g.m.p 信息之后就会对一些不再引用的对象进行清理、压缩以及将除 allp 集合中的第一个 p 之外将状态全部置为 _Pidle，并将其放入调度器 sched.pidle 全局空闲队列中去。 

```go
func procresize(nprocs int32) *p {
	...
	mcache0 = nil
	// release resources from unused P's
	for i := nprocs; i < old; i++ {
		p := allp[i]
		p.destroy()
		// can't free P itself because it can be referenced by an M in syscall
	}
	// Trim allp.
	if int32(len(allp)) != nprocs {
		lock(&allpLock)
		allp = allp[:nprocs]
		idlepMask = idlepMask[:maskWords]
		timerpMask = timerpMask[:maskWords]
		unlock(&allpLock)
	}
    var runnablePs *p
	for i := nprocs - 1; i >= 0; i-- {
		p := allp[i]
		if _g_.m.p.ptr() == p {
			continue
		}
		p.status = _Pidle
		if runqempty(p) {
			pidleput(p)
		} else {
			p.m.set(mget())
			p.link.set(runnablePs)
			runnablePs = p
		}
	}
	...
    return runnablePs
}
```

### 小结

关于调度器启动总起来就是如下步骤：

- 程序启动，编译器调用 runtime.schedinit，初始化系统信息、gc 初始化以及其它相关的信息
- 调用 `getg` 获取当前
- 设置调度器相关的信息（如 maxmcount、procs），其中设置 procs 还会涉及扩充 resize。
- 设置完 procs 就是要对当前的 g 绑定对应的 p，所以就会初始化 p，将全局的 allp[0] 绑定到当前的 g 下并将其余的全部推送到调度器全局空闲队列中。

## 新建 Goroutine

其实我们可以从一个例子着手，查看 go 是如何启动一个 goroutine 的

```go
func startg() {
	go func() {
		fmt.Println("start g")
	}()
}
```

在启动 main.go 的时候，runtime 会执行 proc.go.main 方法创建主协程，并初始化一些信息以及 gc 相关的标识等操作。我们可以通过

`go build -gcflag -S startg.go` 命令能查看，编译器调用了 [runtime.newproc(SB)](https://github.com/golang/go/blob/aa4e0f528e1e018e2847decb549cfc5ac07ecf20/src/runtime/proc.go#L4250)，这个方法有两个参数，一个是参数的大小，另一个是 goroutine 要执行的函数体。

```go
func newproc(siz int32, fn *funcval) {
	argp := add(unsafe.Pointer(&fn), sys.PtrSize)
	gp := getg()
	pc := getcallerpc()
	systemstack(func() {
		newg := newproc1(fn, argp, siz, gp, pc)
		_p_ := getg().m.p.ptr()
		runqput(_p_, newg, true)
		if mainStarted {
			wakep()
		}
	})
}
```

newproc 方法主要就是保存这两个参数的信息以及对应的程序计数器 pc。然后会根据这些变量来新生成一个 g，然后把这个新生成的 g 推送到当前 g 上的线程的处理器 p 的局部 runq 队列中，然后根据特定的条件（mainStarted）来决定是否唤醒。

newproc1 除了一些栈空间大小的判断以及参数、调度器的内存地址拷贝之外，主要执行了如下功能：

- 创建新的 g
- 给 g 分配栈
- 更改 g 的状态
- 更改 g 的属性（调度器的指针，计数器等）

```go
func newproc1(fn *funcval, argp unsafe.Pointer, narg int32, callergp *g, callerpc uintptr) *g {
	_g_ := getg()
	...
	acquirem()
	...
	_p_ := _g_.m.p.ptr()
	newg := gfget(_p_)
	if newg == nil {
		newg = malg(_StackMin)
		casgstatus(newg, _Gidle, _Gdead)
		allgadd(newg) // publishes with a g->status of Gdead so GC scanner doesn't look at uninitialized stack.
	}
	...
    if narg > 0 {
		memmove(unsafe.Pointer(spArg), argp, uintptr(narg))
		if writeBarrier.needed && !_g_.m.curg.gcscandone {
			f := findfunc(fn.fn)
			stkmap := (*stackmap)(funcdata(f, _FUNCDATA_ArgsPointerMaps))
			if stkmap.nbit > 0 {
				// We're in the prologue, so it's always stack map index 0.
				bv := stackmapdata(stkmap, 0)
				bulkBarrierBitmap(spArg, spArg, uintptr(bv.n)*sys.PtrSize, 0, bv.bytedata)
			}
		}
	}
    ...
	newg.sched.sp = sp
	newg.stktopsp = sp
	newg.sched.pc = funcPC(goexit) + sys.PCQuantum // +PCQuantum so that previous instruction is in same function
	newg.sched.g = guintptr(unsafe.Pointer(newg))
	gostartcallfn(&newg.sched, fn)
	newg.gopc = callerpc
	newg.ancestors = saveAncestors(callergp)
	newg.startpc = fn.fn
	if _g_.m.curg != nil {
		newg.labels = _g_.m.curg.labels
	}
	if isSystemGoroutine(newg, false) {
		atomic.Xadd(&sched.ngsys, +1)
	}
	casgstatus(newg, _Gdead, _Grunnable)
	...
	releasem(_g_.m)
	return newg
}
```

上面的代码我省略了其它不在考虑的代码。在创建新 g 之前献给 m 上锁了防止被抢占，因为后续要对当前的 m 相关的 p 下的局部队列保存 g。

在创建 newg 的时候首先会调用 [gfget(_p_)](https://github.com/golang/go/blob/aa4e0f528e1e018e2847decb549cfc5ac07ecf20/src/runtime/proc.go#L4468) 从当前 p 下的 gFree 局部队列中获取空闲的 g（状态为 Gdead），如果局部队列中没有的话，就从调度器 sched 的全局队列中窃取空闲的 g。

> 如果发生了窃取，那么就会在第一次窃取时就把调度器 sched 中的空闲 g 的批次的全部窃取到自己的局部队列中直到局部队列满（n = 32）。

如果全局队列中也没有 g 的话。那么就会调用 [malg(_StackMin)](https://github.com/golang/go/blob/aa4e0f528e1e018e2847decb549cfc5ac07ecf20/src/runtime/proc.go#L4219) 根据传入的栈大小生成 newg。然后就会调用 memmove 指令将数据以及 fn 信息拷贝到栈上。最后就会将前面保存的栈指针以及 fn 程序计数器等信息保存在 newg 上，并更改 newg 状态由 _Gdead 转变为 _Grunnable。最后释放 _g_.m 并返回 newg。

新生成并返回的 newg 最终会由 [runqput](https://github.com/golang/go/blob/aa4e0f528e1e018e2847decb549cfc5ac07ecf20/src/runtime/proc.go#L5947) 推送至当前 p 下的队列中。

```go
func runqput(_p_ *p, gp *g, next bool) {
	...
	if next {
	retryNext:
		oldnext := _p_.runnext
		if !_p_.runnext.cas(oldnext, guintptr(unsafe.Pointer(gp))) {
			goto retryNext
		}
		if oldnext == 0 {
			return
		}
		// Kick the old runnext out to the regular run queue.
		gp = oldnext.ptr()
	}
retry:
	h := atomic.LoadAcq(&_p_.runqhead) // load-acquire, synchronize with consumers
	t := _p_.runqtail
	if t-h < uint32(len(_p_.runq)) {
		_p_.runq[t%uint32(len(_p_.runq))].set(gp)
		atomic.StoreRel(&_p_.runqtail, t+1) // store-release, makes the item available for consumption
		return
	}
	if runqputslow(_p_, gp, h, t) {
		return
	}
	// the queue is not full, now the put above must succeed
	goto retry
}
```

runqnext 会根据传入的 next 参数决定走两个分支：

- true：直接将 g 传给当前 p 下的 runnext
- false：会判断 p 中 runq 中的元素是否已满，如果没满则插入队尾；否则则走 runqputslow。

runqputslow 在 p 的局部队列满的情况下，负责取出队列中的一部分以及待加入的新 g 添加到调度器的全局运行队列上。

```go
func runqputslow(_p_ *p, gp *g, h, t uint32) bool {
	var batch [len(_p_.runq)/2 + 1]*g
	// First, grab a batch from local queue.
	n := t - h
	n = n / 2
	for i := uint32(0); i < n; i++ {
		batch[i] = _p_.runq[(h+i)%uint32(len(_p_.runq))].ptr()
	}
	...
	batch[n] = gp
	// Link the goroutines.
	for i := uint32(0); i < n; i++ {
		batch[i].schedlink.set(batch[i+1])
	}
	...
	var q gQueue
	q.head.set(batch[0])
	q.tail.set(batch[n])

	// Now put the batch on global queue.
	lock(&sched.lock)
	globrunqputbatch(&q, int32(n+1))
	unlock(&sched.lock)
	return true
}
```

要注意，globrunqputbatch 在添加全局队列 sched.runq 前后是要加锁的防止并发修改这个共享的全局变量。

> 注意：关于 p 本地运行队列 runq 和调度器 sched 的运行队列 runq 同样都是链表，但是组成的结构完全不一样。
>
> p 的 runq 是通过数组+双指针形成的环形队列。
>
> sched 的 runq 就是单纯的链表结构

## 调度

在执行完 `schedinit` 之后就会调用创建 M 的入口函数 [runtime.mstart](https://github.com/golang/go/blob/aa4e0f528e1e018e2847decb549cfc5ac07ecf20/src/runtime/proc.go#L1339)，前者内部会调用 [runtime.mstart1](https://github.com/golang/go/blob/aa4e0f528e1e018e2847decb549cfc5ac07ecf20/src/runtime/proc.go#L1380)。前者主要初始化 g0 的 stackguard0 和 stackguard1 字段。

























## 参考链接

- https://docs.google.com/document/d/1TTj4T2JO42uD5ID9e89oa0sLKhJYD0Y_kqxDv3I3XMw/edit#
- https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-goroutine/
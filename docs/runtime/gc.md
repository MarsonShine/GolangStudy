# Go 垃圾回收算法

GC 与 mutator 线程并发运行，类型准确（又名精确）（ is type accurate (aka precise)），允许多个 GC 线程并行运行。它是一个使用写屏障的**并发标记和扫描**。它是非分代和非压缩的。内存分配是通过使用每个 P 分配区域来完成的，彼此是隔离的，以最大程度地减少碎片，同时消除常见情况下的锁定。（Allocation is done using size segregated per P allocation areas to minimize fragmentation while eliminating locks in the common case.）

该算法分解为几个步骤。

> 这是对所使用的算法的一个高级描述。要了解 GC 的概况，可以从 Richard Jones 的[gchandbook.org](https://gchandbook.org/)开始。
>
> 该算法的知识遗产包括Dijkstra的即时算法（on-the-fly algorithm），见 Edsger W. Dijkstra, Leslie Lamport, A. J. Martin, C. S. Scholten 和 E. F. M. Steffens, 1978。
>
> 即时垃圾收集:一种合作练习。Commun. ACM 21, 11(1978年11月)，966-975。
>
> 有关这些步骤完整、正确和终止的期刊质量证明，请参见 Hudson, R. 和 Moss, J.E.B. Copying Garbage ollection without stopping the world.
>
> 并发与计算:实践与经验15(3-5)，2003。

1. GC 执行扫描终止阶段
   1. STW。这将导致所有的 p 到达一个 GC 安全点（safe-point）。
   2. 扫描所有未扫过的块（span）。这种情况只有在预期时间之前强制执行此 GC 周期时，才会有未扫描块。
2. GC 执行标记阶段
   1. 通过将 `gcphase` 设置为 `_GCmark`（来自 _GCoff）、启用写屏障、启用 mutator 辅助和将根标记作业排队，为标记阶段做准备。在所有 P 都启用写屏障之前，不能扫描任何对象，这是使用 STW 完成的。
   2. 世界恢复执行。从此刻开始，GC 工作由调度器启动的 mark workers 和作为分配的一部分执行的 assists 完成的。对于任何指针写操作，写屏障将覆盖指针和新指针的值遮暗（有关详细信息，请参见 mbarrier.go）。**新分配的对象立即被标记为黑色**。
   3. GC 执行根标记作业。这包括扫描所有堆栈，对所有全局变量进行着色，以及对堆外运行时数据结构中的任何堆指针进行着色。扫描堆栈会停止 goroutine，隐藏在其堆栈上找到的所有指针，然后恢复 goroutine。
   4. GC 清空**灰色对象**的工作队列，将扫描每个灰色对象将其置为黑色，并对对象中找到的所有指针进行着色（这反过来可能会将这些指针添加到工作队列中）。
   5. 由于 GC 工作分布在本地缓存中，因此 GC 使用分布式终止算法（distributed termination algorithm）来检测何时不再有根标记作业或灰色对象（请参阅 gcMarkDone）。 此时，GC 过渡到标记终止阶段。
3. GC 执行标记终止阶段
   1. STW
   2. 设置 `gcphase` 为 `_GCmarktermination`，并禁止 worker 和 assists
   3. 执行诸如刷新 mcache 之类的内务管理操作
4. GC执行清除阶段。
   1. 通过将 `gcphase` 设置为 `_GCoff`、设置清除状态并禁用写屏障来准备清楚阶段。
   2. 世界恢复运行。从此刻开始，新分配的对象是**白色**的，如有必要，在使用前分配扫描块（sweeps spans）。
   3. GC 会在后台并行清除并响应分配。请参阅下面的说明
5. 当有足够的分配时，会重复执行上面的序列1开始的步骤。请参阅下面对 GC 率的讨论。

## 并发清除

扫描阶段与正常的程序执行同时进行。堆都是逐个惰性扫描的，（当一个 goroutine 需要另一个块（span）时）以及在后台 goroutine 中同时进行（这有助于不受 CPU 限制的程序）。在 STW 标志终止的末尾，所有块被标记为“需要清扫”。

后台清除器只是简单的一个个清除。


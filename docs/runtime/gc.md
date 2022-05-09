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
   2. 扫描所有未扫过的 span。这种情况只有在预期时间之前强制执行此 GC 周期时，才会有未扫描 span。
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
   2. 世界恢复运行。从此刻开始，新分配的对象是**白色**的，如有必要，在使用前分配扫描 span（sweeps spans）。
   3. GC 会在后台并行清除并响应分配。请参阅下面的说明
5. 当有足够的分配时，会重复执行上面的序列1开始的步骤。请参阅下面对 GC 率的讨论。

## 并发清除

扫描阶段与正常的程序执行同时进行。堆都是逐个惰性扫描的，（当一个 goroutine 需要另一个 span 时）以及在后台 goroutine 中同时进行（这有助于不受 CPU 限制的程序）。在 STW 标志终止的末尾，所有的 span 被标记为“需要清扫”。

后台清除器只是简单的一个个清除。

为了避免在存在未扫描的 span 时请求更多的 OS 内存，当 goroutine 需要另一个 span 时，它首先尝试通过扫描来回收那么多内存。 当一个 goroutine 需要分配一个新的小对象 span 时，它会扫描相同大小的小对象 span，直到它释放至少一个对象。 当一个 goroutine 需要从堆中分配大对象 span 时，它会扫描 span，直到它至少将那么多页释放到堆中。 在一种情况下这可能不够：如果一个 goroutine 扫描并释放两个不相邻的一页 span 到堆中，它将分配一个新的两页 span，但仍然可能存在其他一页未扫描的 span，这可能是合并成一个两页的 span。

确保在未扫描的 span 上不进行任何操作（这会破坏 GC 位图中的标记位）至关重要。 在 GC 期间，所有 mcache 都被刷新到中心缓存中，因此它们是空的。 当一个 goroutine 将一个新的 span 抓取到 mcache 中时，GC 会清扫它。 当 goroutine 显式释放对象或设置终结器时，它要确保扫描了 span（通过扫描它，或等待并发扫描完成）。终结器 goroutine 仅在扫描所有 span 时才启动。当下一次 GC 开始时，它会扫描所有尚未扫描的 span（如果有的话）。

GC 频率

下一次 GC 是在我们分配了与已经使用的内存量成正比的额外内存量之后。 该比例由 GOGC 环境变量控制（默认为 100）。 如果 GOGC=100 并且我们正在使用 4M，我们将在达到 8M 时再次进行 GC（此标记在 gcController.heapGoal 变量中跟踪）。这使 GC 成本与分配成本成线性比例。调整 GOGC 只会改变线性常数（以及使用的额外内存量）。

Oblets

为了防止在扫描大型对象时出现长时间的停顿并提高并行性，垃圾收集器将对大于 `maxObletBytes` 的对象的扫描作业分解为最多为 `maxObletBytes` 的“oblets”。当扫描遇到一个大对象的开头时，它只扫描第一个 oblet 并将剩余的 oblet 作为新的扫描作业排入队列。

## 内存分配器

早期的内存分配器是基于[tcmalloc](http://goog-perftools.sourceforge.net/doc/tcmalloc.html)。

主分配器工作运行在页面中。小对象分配大小（最多 32 kB）被舍入到大约 70 个大小类之一，每个大小类都有自己的自由对象集，大小正好是该大小。任何空闲的内存页都可以拆分为一组大小相同的对象，然后使用空闲 bitmap 对其进行管理。

allocator 的数据结构如下：

cheap: malloc 堆，以页粒度管理。

mspan: 由 mheap 管理的正在运行的使用的页。

mcentral: 收集给定大小类的所有 span。

mcache: 每个 P 的 span 缓存都有空闲的空间。

mstats: 分配统计

### 小对象的分配逻辑

分配一个小对象将沿着缓存层次结构向上进行：

1. 当对象大小上升到小对象尺寸时，就会此 P 的 mcache 中查看相应的 mspan。扫描 mspan的空闲 bitmap，找到空闲的 slot。如果有空闲的 slot，分配它。这一切都可以在不获取锁的情况下完成。
2. 如果 mspan 没有空闲 slot，则从 mcentral 的具有可用空间的所需指定大小级别的 mspan 列表中获取新的 mspan。获得整个 span 分摊了锁定 mcentral 的成本。
3. 如果 mcentral 的 mspan 列表为空，则从 mheap 获取一系列页面以用于 mspan。
4. 如果 mheap 是空的或是没有足够大的页运行，则从操作系统分配一组新的页（至少 1MB）。分配大量运行页可以摊销与操作系统对话的成本。

扫描一个 mspan 并在其上释放对象会沿类似的层次结构进行：

1. 如果 mspan 在响应分配的时候处于扫描阶段，则将其返回到 mcache 以满足分配。 
2. 否则，如果 mspan 里面还有分配的对象，它被放置在 mspan 的 size 类的 mccentral 空闲列表中。
3. 否则，如果 mspan 中的所有对象都是空闲的，则将 mspan 的页返回到 mheap 并且 mspan 立即死亡。

### 大对象的分配逻辑

分配和释放大对象直接使用 mheap，绕过 mcache 和 mcentral。如果 mspan.needzero 为 false，则 mspan 中的空闲对象 slot 已准备清零。否则，如果 needzero 为 true，则在分配对象时将其清零。以这种方式延迟清零有以下好处： 

1. 栈帧分配可以完全避免归零。
2. 它表现出更好的时间局部性，因为程序可能即将写入内存。 
3. 我们不会将永不复用的页归零。

### 虚拟内存布局

堆由一组 arenas 组成，在 64 位上为 64MB，在 32 位上为 4MB (heapArenaBytes)。每个 arena 的起始地址也与 arena 大小对齐。

每个 arena 都有一个关联的 heapArena 对象，用于存储该 arena 的元数据：arena 中所有单词的 heap bitmap 和 arena 中所有页的 span map。heapArena 对象本身是在堆外分配的。

由于 arena 是对齐的，地址空间可以被视为一系列 arena 帧。arena map (mheap_.arenas) 从 arena 帧号映射到 *heapArena，或者对于 Go 堆不支持的部分地址空间的映射为 nil。arena map 结构为两级数组，由“L1” arena map 和许多“L2” arena map 组成；然而，由于 arena 很大，因此在许多架构上，arena map 由单个大型 L2 map 组成。

arena map 覆盖了整个可能的地址空间，允许 Go 堆使用地址空间的任何部分。分配器尝试保持 arenas 连续，以便大 span（以及大对象）可以跨越 arenas。

### mspan

### bitmap

### arena

# 相关拓展阅读

https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html

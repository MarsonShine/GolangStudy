# Goroutine 与运行时

Go 是根据 fork-join 模式执行并发的。Forks 是当 goroutine 开始之时，而联合点（join points）是当两个或更多 goroutine 通过 channel 或 sync 类同步之时。工作窃取算法依据一些基本的规则。给定一个执行的线程：

1. 在 fork 点，添加任务到与该线程相关联队列的尾部。
2. 如果这个线程是空闲的，就从**其他相关的随机线程的队列头部**窃取工作。
3. 在一个此时还不能实现的连接点（join point）上（例如与它同步的 goroutine 还没有完成），在自己线程的队列尾部出栈工作
4. 如果这个线程队列是空的：
   1. 在 join 处排队
   2. 从一个随机线程的相关队列的头部窃取工作

这有一点抽象，那就来看看代码，看下里面操作的算法，它递归计算一个 Fibonacci 队列：

```go
var fib func(n int) <-chan int

func sample() {
	fib = func(n int) <-chan int {
		result := make(chan int)
		go func() {
			defer close(result)
			if n <= 2 {
				result <- 1
				return
			}
			result <- <-fib(n-1) + <-fib(n-2)
		}()
		return result
	}

	fmt.Printf("fib(4) = %d", <-fib(4))
}
```

我们来看一下这个版本的工作窃取算法将会执行什么操作。这段程序我们假设是在有两个核芯处理器的机器上。我们将会为每个处理器分配一个 OS 线程，分别命名为 T1 和 T2。回到这个例子，我将从 T1 转到 T2 需要尽量提供一些结构。事实上，这些都不是决定性的。

我们程序开始。初始阶段，我们开始了一个 goroutine，即 main goroutine，并且我们假设是调度到了处理器 T1：

| T1 call stack  | T1 work deque | T2 call stack | T2 work deque |
| -------------- | ------------- | ------------- | ------------- |
| main goroutine |               |               |               |

下一步，我们到达调用 fib(4) 的地方。这个 goroutine 将会被调度并且放置在 T1 工作队列的尾部，并且父级 goroutine 会继续持续处理：

| T1 call stack  | T1 work deque | T2 call stack | T2 work deque |
| -------------- | ------------- | ------------- | ------------- |
| main goroutine | fib(4)        |               |               |

到此时此刻，这会依赖于时间，这里有两件事会任意一个会触发：T1 或是 T2 将窃取正在调用 fib(4) 的 goroutine。举个例子，为了更清晰的分析这个算法，我们假设 T1 成功窃取；但是，要值得注意的是其他线程也有可能赢得窃取。

| T1 call stack                                       | T1 work deque | T2 call stack | T2 work deque |
| --------------------------------------------------- | ------------- | ------------- | ------------- |
| (main grouting)(unrealalized join point)<br/>fib(4) |               |               |               |

fib(4) 运行在 T1 上并且由于从左右到添加的操作顺序——f(3) 入栈然后再是 fib(2) 进入到这个队列的尾部：

| T1 call stack                                      | T1 work deque     | T2 call stack | T2 work deque |
| -------------------------------------------------- | ----------------- | ------------- | ------------- |
| (main goroutine)(unrealized join point)<br/>fib(4) | fib(3)<br/>fib(2) |               |               |

此时，T2 仍是空闲的，所以它能将 fib(3) 从 T1 队列的头部解救出来。要注意这里的 fib(2)——最后的 fib(4) 推入到队列，因此 T1 第一件事就是要留在 T1 计算。我们稍后讨论这是为什么。

| T1 call stack                                      | T1 work deque | T2 call stack | T2 work deque |
| -------------------------------------------------- | ------------- | ------------- | ------------- |
| (main goroutine)(unrealized join point)<br/>fib(4) | fib(2)        | fib(3)        |               |

与此同时，T1 到达在 fib(4) 上无法继续工作的地方，因为它要等待 channel 返回从 fib(3) 和 fib(2) 的结果。这就是 unrealized join poin 在我们这个算法的第 3 阶段。正因如此，它从自己队列的尾部弹出工作，这里是 fib(2)：

| T1 call stack                                                | T1 work deque | T2 call stack | T2 work deque |
| ------------------------------------------------------------ | ------------- | ------------- | ------------- |
| (main goroutine)(unrealized join point)<br/>fib(4) (unrealized join point)<br/>fib(2) |               | fib(3)        |               |

到这里就有点迷惑了。因为我们没有在递归算法中使用回溯（backtracking），我们打算调度另一个 goroutine 来计算 fib(2)。这是一个新且单独的，从之前调度到 T1 的新 goroutine。这个是之前刚调度到 T1 的，是调用 fib(4) 的一部分（如 4-2）；新的 goroutine 是调用的 fib(3) 的部分（如 3-1）。下面是调用 fib(3) 时新调度的 goroutine：

| T1 call stack                                                | T1 work deque | T2 call stack | T2 work deque     |
| ------------------------------------------------------------ | ------------- | ------------- | ----------------- |
| (main goroutine) (unrealized join point)<br/>fib(4)(unrealized join point)<br/>fib(2) |               | fib(3)        | fib(2)<br/>fib(1) |

接下来，T1 到达我们递归算法的条件终止部分（n <= 2）和返回 1：

| T1 call stack                                                | T1 work deque | T2 call stack | T2 work deque    |
| ------------------------------------------------------------ | ------------- | ------------- | ---------------- |
| (main goroutine) (unrealized join point)<br/> fib(4) (unrealized join point) <br/>(returns 1) |               | fib(3)        | fib(2)<br>fib(1) |

然后 T1 再一次空闲，所以它从 T2 队列的头部窃取获取 fib(2)：

| T1 call stack                                                | T1 work deque | T2 call stack                           | T2 work deque |
| ------------------------------------------------------------ | ------------- | --------------------------------------- | ------------- |
| (main goroutine) (unrealized join point)<br>fib(4) (unrealized join point)<br>fib(2) |               | fib(3)(unrealized join point)<br>fib(1) |               |

T2 然后在此到达（n<=2）和返回 1：

| T1 call stack                                                | T1 work deque | T2 call stack                                 | T2 work deque |
| ------------------------------------------------------------ | ------------- | --------------------------------------------- | ------------- |
| (main goroutine) (unrealized join point)<br>fib(4) (unrealized join point)<br>fib(2) |               | fib(3) (unrealized join point)<br>(returns 1) |               |

下一步，T1 到达基础判断以及返回 1：

| T1 call stack                                                | T1 work deque | T2 call stack                                 | T2 work deque |
| ------------------------------------------------------------ | ------------- | --------------------------------------------- | ------------- |
| (main goroutine) (unrealized join point)<br>fib(4) (unrealized join point)<br>return 1 |               | fib(3) (unrealized join point)<br>(returns 1) |               |

T2 调用 fib(3) 时有两个 realized join point；那就是说都会调用 fib(2) 和 fib(1) 已经在它们的 channel 中返回结果，并且两个 goroutine 会连接回它们的父 goroutine——fib(3) 的调用者。它执行加法（1+1=2）并且返回它们 channel 上的结果：

| T1 call stack                                                | T1 work deque | T2 call stack | T2 work deque |
| ------------------------------------------------------------ | ------------- | ------------- | ------------- |
| (main goroutine) (unrealized join point)<br>fib(4) (unrealized join point)<br> |               | (return 2)    |               |

相同的事情在此发生：寄宿在调用 fib(4) 的 goroutine 也有两个 unrealized join point：fib(3) 和 fib(2)。我们在前面的步骤中刚刚完成了fib(3)的连接，当最后一个任务 T2 完成时，到 fib(2) 的连接也完成了。再次，执行了（2+1=3）并且返回了在 channel 上 fib(4) 的结果：

| T1 call stack                                           | T1 work deque | T2 call stack | T2 work deque |
| ------------------------------------------------------- | ------------- | ------------- | ------------- |
| (main goroutine) (unrealized join point)<br> (return 3) |               |               |               |

这时候，我们已经在 main goroutine 的实现了连接点（<- fib(4)），并且 main goroutine 得以继续。它会打印结果：

| T1 call stack | T1 work deque | T2 call stack | T2 work deque |
| ------------- | ------------- | ------------- | ------------- |
| print 3       |               |               |               |

现在，我们来解释一下这个算法中有趣的部分。回想一下，执行线程既会入栈，也会(在必要时)从其工作队列的尾部出栈。放在队列的尾部有一些有趣的属性：

- *这是完成父节点的连接最可能需要的工作*

  更快地完成连接意味着我们的程序可能会执行得更好，并且在内存中保留的东西也更少

- *它可能仍在处理器缓存中工作*

  因为这是线程在当前工作之前所做的最后工作，这个信息很可能会保存在执行线程的 CPU 的缓存中。这意味着会更少的缓存丢失！

总之，以这种方式调用工作有很多隐含的性能好处。

## 窃取任务或延续（Stealing Tasks or Continuations）

有一件事我们忽略了，那就是我们入队列的问题和窃取的问题是什么。在 fork-join 范例下，这里有亮点：任务与延续。为了确保能够清晰理解在 Go 的任务和延续，我们再来看下 fibonacci 程序：

```go
var fib func(n int) <-chan int 
fib = func(n int) <-chan int {
  result := make(chan int) 
  go func() { // 1
  	defer close(result) 
  	if n<=2{
  			result <- 1
  			return
		}
		result <- <-fib(n-1) + <-fib(n-2) 
	}()
	return result // 2
}
fmt.Printf("fib(4) = %d", <-fib(4))
```

1. 在 Go 中，goroutine 就是任务
2. 在一个 goroutine 被调用的之后的所有一切都是延续

在我们前面的分布式队列工作窃取算法演练中，我们正在对任务(或 goroutines )进行排队。因为 goroutine 能很好的封装工作体，这是一种很自然的方式；但是，这里实际不是 Go 的工作窃取算法的工作方式。Go 的工作窃取算法是入队列以及窃取延续。

为什么这很重要呢？排队和窃取延续做了什么是我们排队和窃取任务所没有做的呢？为了回答这个问题，我们来看一下连接点（join point）。

在我们的算法之下，当一个线程到达未实现连接点（unrealized join point）时，这个线程必须暂停执行并完成对一个任务的窃取。由于正当查看这个任务在做的事时要停止连接，这被称为 stalling join（暂停连接）。在窃取任务和窃取延续算法都会暂停连接，但是在发生暂停的频率上是有很大意义的。

试想一下：当创建一个 goroutine，这很像是你在程序中开启一个函数在 go func 中执行。从 goroutine 的延续将在某个点想与 goroutine 连接，这是很合理的可能。延续在 goroutine 完成之前尝试连接的情况并不少见。根据这些原理，当调度一个 goroutine 时，它就要立即开始工作这是很有意义的。

现在回想一下线程入栈和出栈工作到/从队列的尾部的属性，以及其它线程从头部出栈工作。 如果我们压入一个延续至队列的尾部，这最不可能被另一个线程从队列的头部弹出被窃取，因此当我们完成 goroutine 的执行时，我们很有可能会将它捡起来从而避免了暂停。这也使得分叉（fork）的任务看起来很像一个函数调用：线程跳转执行 goroutine 并且在它完成之后返回给延续。

让我们看看如何将延续窃取应用到斐波那契程序中。由于表示延续（continuation）要比任务（tasks）更不明确，我们将使用以下约定:

- 当一个 continuation 在工作队列中排队时，我们将它列为 con. of X。
- 当一个 continuation 弹出执行时，我们将显式的转换为 continuation 到下一个 fib 调用

下面是 Go 运行时正在做的事情的一个更详细的表示。

| T1 call stack | T1 work deque | T2 call stack | T2 work deque |
| ------------- | ------------- | ------------- | ------------- |
| main          |               |               |               |

主 goroutine 调用 fib4 并且这个调用的延续入栈到 T1 的工作队列的尾部:

| T1 call stack | T1 work deque | T2 call stack | T2 work deque |
| ------------- | ------------- | ------------- | ------------- |
| fib4          | cont. of main |               |               |

T2 是空闲的，所以它会窃取 main 的延续：

| T1 call stack | T1 work deque | T2 call stack | T2 work deque |
| ------------- | ------------- | ------------- | ------------- |
| fib4          |               | cont. of main |               |

调用 fib4 然后调度 fib3，它是立即执行的，并且 T1 入栈 fib4 的延续到它的队列尾部：

| T1 call stack | T1 work deque | T2 call stack | T2 work deque |
| ------------- | ------------- | ------------- | ------------- |
| fib3          | cont. of fib4 | cont. of main |               |

当 T2 企图执行 main 延续时，它到达了一个还未实现连接的点；因此，它会从 T1 的队列窃取更多的工作。这次，它是调用 fib4 的延续了：

| T1 call stack | T1 work deque | T2 call stack                                          | T2 work deque |
| ------------- | ------------- | ------------------------------------------------------ | ------------- |
| fib3          |               | cont. of main (unrealized join point)<br>cont. of fib4 |               |

接下来，T1 调用 fib3 的调用 goroutine 会调用 fib3 并立即调用。那 fib3 的延续就会被入栈到它的工作队列的尾部：

| T1 call stack | T1 work deque | T2 call stack                                          | T2 work deque |
| ------------- | ------------- | ------------------------------------------------------ | ------------- |
| fib2          | cont. of fib3 | cont. of main (unrealized join point)<br>cont. of fib4 |               |

T2 执行 fib4 的延续，从 T1 结束的地方开始，并且它调度 fib2，开始立即执行并再次入队列 fib4:

| T1 call stack | T1 work deque | T2 call stack                                  | T2 work deque |
| ------------- | ------------- | ---------------------------------------------- | ------------- |
| fib2          | cont. of fib3 | cont. of main (unrealized join point)<br/>fib2 | cont. of fib4 |

下一步，T1 调用 fib2 到达了递归算法的基础判断并返回 1：

| T1 call stack | T1 work deque | T2 call stack                                  | T2 work deque |
| ------------- | ------------- | ---------------------------------------------- | ------------- |
| return 1      | cont. of fib3 | cont. of main (unrealized join point)<br/>fib2 | cont. of fib4 |

然后 T2 也达到了判断条件并返回 1:

| T1 call stack | T1 work deque | T2 call stack                                          | T2 work deque |
| ------------- | ------------- | ------------------------------------------------------ | ------------- |
| return 1      | cont. of fib3 | cont. of main (unrealized join point)<br/>（return 1） | cont. of fib4 |

T1 随后从它自己的队列中窃取并开始执行 fib1。注意，T1 上的调用链是：fib3 -> fib2 -> fib1。这在我们早之前就讨论过窃取延续的好处。

| T1 call stack | T1 work deque | T2 call stack                                          | T2 work deque |
| ------------- | ------------- | ------------------------------------------------------ | ------------- |
| fib1          |               | cont. of main (unrealized join point)<br/>（return 1） | cont. of fib4 |

T2 就到达了 fib4 的延续的最终点，但是只能一个连接点到达：fib2。fib3 的调用仍在 T1 处理器上。T2 由于没有可窃取的工作，于是就空闲了：

| T1 call stack | T1 work deque | T2 call stack                                                | T2 work deque |
| ------------- | ------------- | ------------------------------------------------------------ | ------------- |
| fib1          |               | cont. of main (unrealized join point)<br/>fib(4) (unrealized join point) |               |

T1 现在到达了延续的最终点，fib3，并且 fib2 和 fib1 两个连接点都满足了。T1 返回 2:

| T1 call stack | T1 work deque | T2 call stack                                          | T2 work deque |
| ------------- | ------------- | ------------------------------------------------------ | ------------- |
| return 2      |               | cont. of main (unrealized join point)<br/>（return 2） |               |

现在 fib4，fib3 以及 fib2 都已经满足。T2 也就能够执行计算并返回结果了（2+1=3）：

| T1 call stack | T1 work deque | T2 call stack                                          | T2 work deque |
| ------------- | ------------- | ------------------------------------------------------ | ------------- |
|               |               | cont. of main (unrealized join point)<br/>（return 3） |               |

最后主 goroutine 的连接点已经实现并接收到了调用 fib4 的结果并打印结果：

| T1 call stack | T1 work deque | T2 call stack    | T2 work deque |
| ------------- | ------------- | ---------------- | ------------- |
|               |               | main（prints 3） |               |

当我们走这个过程时，我们简要地看到了延续如何帮助在 T1 上连续地执行事情。如果我们看看这个运行(使用延续偷窃)与使用任务偷窃的运行的统计数据，就会发现更清晰的好处：

| Statistics       | Continuation stealing | Task stealing         |
| ---------------- | --------------------- | --------------------- |
| # Steps          | 14                    | 15                    |
| Max Deque Length | 2                     | 2                     |
| # Stalled Joins  | 2（都在空闲的线程上） | 3（都在繁忙的线程上） |
| 调用堆栈的大小   | 2                     | 3                     |

这些统计数据可能看起来很接近，但是如果我们从更大的项目中推断，我们就可以开始看到持续窃取是如何提供显著的好处的。

让我们看一下当只有一个线程运行是什么样子：

| T1 call stack | T1 work  duque |
| ------------- | -------------- |
| main          |                |

| T1 call stack | T1 work  duque |
| ------------- | -------------- |
| fib4          | main           |

| T1 call stack | T1 work  duque |
| ------------- | -------------- |
| fib3          | main           |
|               | cont.of fib4   |

| T1 call stack | T1 work  duque |
| ------------- | -------------- |
| fib2          | main           |
|               | cont. of fib4  |
|               | cont. of fib3  |

| T1 call stack | T1 work  duque |
| ------------- | -------------- |
| return 1      | main           |
|               | cont. of fib4  |
|               | cont. of fib3  |

| T1 call stack | T1 work  duque |
| ------------- | -------------- |
| fib1          | main           |
|               | cont. of fib4  |

| T1 call stack | T1 work  duque |
| ------------- | -------------- |
| return 1      | main           |
|               | cont. of fib4  |

| T1 call stack | T1 work  duque |
| ------------- | -------------- |
| return 2      | main           |
|               | cont. of fib4  |

| T1 call stack | T1 work  duque |
| ------------- | -------------- |
| fib2          | main           |

| T1 call stack | T1 work  duque |
| ------------- | -------------- |
| return 1      | main           |

| T1 call stack | T1 work  duque |
| ------------- | -------------- |
| return 3      | main           |

| T1 call stack  | T1 work  duque |
| -------------- | -------------- |
| main (print 3) |                |

在单线程上的运行时使用 goroutine 与我们刚刚调用的函数是相同的。这是窃取延续的另一个好处。

考虑到所有这些，窃取延续被认为在理论上要优先于窃取任务的，因此最好就是延续排队而不是 goroutine。就像下表看到的，窃取延续有这些好处：

|          | Continuation | Child  |
| -------- | ------------ | ------ |
| 队列大小 | 有边界       | 无边界 |
| 执行顺序 | 串型顺序     | 无序   |
| 连接点   | 非暂缓       | 窃取   |

那么为什么不是所有窃取工作的算法都实现了延续窃取呢？好，延续窃取通常需要编译器的支持。幸运的是，Go 有自己的编译器，而且延续窃取是 Go 的工作窃取算法的实现方式。通常没有这种语言能奢侈的实现任务，或者称为 “孩子”，作为窃取类库。

正当这个模型接近 Go 算法时，但它仍然不能代表整体。Go 会执行额外的优化。在我们分析之前，让我们开始使用源代码中列出的 Go 调度器的命名法来设置阶段。

Go 的调度器主要有三个概念：

- G：一个 goroutine
- M：一个操作系统线程（在源代码中以图灵机引用）
- P：上下文（context）（在源代码中以处理器引用）

在我们关于窃取的讨论中，M 等价于 T，P 等价于工作队列（更改 GOMAXPROCS 可以更改想要分配的具体数目）。G 是一个 goroutine，但是要保留表示当前状态的的 goroutine，最显著的就是程序计数器（PC）。它允许一个 G 去表示一个延续，所以 Go 能做到窃取延续。

在 Go 的运行时，这些 M 寄宿在 P 上开始，然后由 P 来调度 G 在 M 上运行。

就我个人而言，当只使用这种符号时，我发现很难理解这种算法是如何工作的，所以我将使用全名来进行分析。好了，让我们看看 Go 的调度程序是如何工作的！

就像我提到的，GOMAXRPOCS 可以设置在运行时有多少可用的上下文。默认设置是在主机上每个上下文对应一个逻辑 CPU。不像上下文，可能会有比内核更多或更少的 OS 线程来帮助 Go 的运行时管理垃圾收集和 goroutines 之类的事情。我提到这个是因为在运行时有一个非常重要的保证：至少总是会有足够的 OS 线程来处理每个上下文。这就允许在运行时做一些重要的优化。运行期也包含了线程池，用于当前未被使用的线程。现在我们来谈一下那些优化。

考虑一下，如果任何一个 goroutines 被输入/输出阻塞，或者在 Go 的运行时之外进行系统调用阻塞，会发生什么情况。承载 goroutine 的操作系统线程也会被阻塞，将无法取得进展或承载任何其他 goroutine。从逻辑上讲，这没什么问题，但是从性能的角度来看，Go 可以做更多的工作来保持机器上的处理器尽可能被利用。

Go 在这种情况下所做的是将上下文与 OS 线程分离，这样上下文就可以被移交给另一个未阻塞的 OS 线程。这就允许上下文进一步调度 goroutines，这允许运行时保持主机的 cpu 活动。阻塞的 goroutine 保留了与阻塞线程相关联的。

当 goroutine 最终变得不阻塞时，主机操作系统线程就会企图将上下文从另一个操作系统线程窃取回来，以至于能够让它继续执行上次阻塞的 goroutine。然而，有时这并不总是可行的。在这种情况下，线程将把它的 goroutine 放在全局上下文中，线程将进入睡眠状态，并且它将被放入运行时的线程池以供将来使用(例如，如果goroutine再次被阻塞)。

我们刚才提到的全局背景不适合我们之前讨论的抽象工作窃取算法。这是 Go 如何优化 CPU 利用率所必需的实现细节。为了确保放置在全局上下文中的 goroutines 不会永远存在，一些额外的步骤被添加到工作窃取算法中。一个上下文将定期检查全局上下文，看是否有任何 goroutines 在那里，当一个上下文的队列是空的，它将首先检查全局上下文，以窃取工作，然后检查其他 OS 线程的上下文。

除了输入/输出和系统调用，Go 还允许 goroutines 在任何函数调用期间被抢占。这与 Go 的思想一致，即通过确保运行时能够有效地调度工作，而偏爱非常细粒度的并发任务。 一个值得注意的例外是团队一直试图解决的 goroutines 不执行输入/输出，系统调用，或函数调用。目前，这些 goroutine 是不可抢占的，可能会导致严重的问题，如长时间的GC等待，甚至死锁。幸运的是，从轶事（anecdotal）的角度来看，这只是微不足道的小事。

## 向开发人员展示所有这些内容

既然您已经理解了goroutines 是如何在幕后工作的，让我们再次回顾并重申开发人员是如何使用 go 关键字进行交互的。就是这样。

在函数或闭包之前加上 “go” 这个词，您就自动安排了一个任务，它将以对运行它的机器最有效的方式运行。作为开发者，我们仍然用我们熟悉的原语思考：函数。我们不需要了解一种新的做事方式、复杂的数据结构或调度算法。

可伸缩性、效率和简单性。这就是 goroutines 如此吸引人的原因。


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

为什么这很重要呢？
# Go 语言设计与实现

go 编译的相关指令文档：https://golang.org/cmd/compile/

用命令查看 Go 语言源代码编译成汇编语言

```cmd
$ go build -gcflags -S main.go
```

获取 Go 语言更详细的编译过程，可以输入以下命令

```cmd
$ GOSSAFUNC=main go build main.go
```

Go 反编译指令得到汇编治理

```bash
go tool compile -S -N -l main.go	// -N -l 是阻止编译器优化汇编代码
```

Go 指令运行测试代码

```bash
go test -gcflags=-N -benchmem -test.count=3 -test.cpu=1 -test.benchtime=1s -bench=.
```



# 目录

- 编译原理
- 数据结构
- 语言基础
- 常用关键字
- 并发
- 内存管理
  - 垃圾收集器
- 元编程
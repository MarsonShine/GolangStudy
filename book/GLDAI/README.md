# Go 语言设计与实现

用命令查看 Go 语言源代码编译成汇编语言

```cmd
$ go build -gcflags -S main.go
```

获取 Go 语言更详细的编译过程，可以输入以下命令

```cmd
$ GOSSAFUNC=main go build main.go
```


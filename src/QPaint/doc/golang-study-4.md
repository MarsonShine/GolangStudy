# golang 自学系列（四）——（调试）VSCode For Debug

这里如何装 vscode 我就不说了

这里如何在 vscode 正常写代码我也不说了

在能正常用 vscode 写 go 语言的前提下（何为正常？就是写代码有智能提示的那种）

在 终端/cmd/iterm 输出以下命令

```cmd
xcode-select --install	// vscode 第一次运行这个命令会弹出一个提示是否安装这个软件，点击是即可
go install github.com/derekparker/delve/cmd/dlv
```

在执行第二条命了的时候，如果你没有获取指定仓库代码，就会有如下异常

```
can't load package: package github.com/derekparker/delve/cmd/dlv: cannot find package "github.com/derekparker/delve/cmd/dlv" in any of:
        /xxx/xxx/src/github.com/derekparker/delve/cmd/dlv (from $GOROOT)
        /xxx/xxx/github.com/derekparker/delve/cmd/dlv (from $GOPATH)
```

这个时候只要执行以下命令即可

```
go get github.com/derekparker/delve/cmd/dlv
```

然后 F5 调试即可成功
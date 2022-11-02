package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/marsonshine/mscli/cmd/exec"
	"github.com/marsonshine/mscli/cmd/new"
	"github.com/spf13/cobra"
)

var (
	info string
)

func NewDefaultMsCliCommand() *cobra.Command {
	return NewDefaultMsCliCommandWithArgs(os.Args, os.Stdin, os.Stdout, os.Stderr)
}

func NewDefaultMsCliCommandWithArgs(args []string, in io.Reader, out, errout io.Writer) *cobra.Command {
	cmd := NewMsCliCommand(in, out, errout)

	return cmd
}

func NewMsCliCommand(in io.Reader, out, errout io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "msctl",
		Short: "msctl 是一个学习 go cli 实现的练习项目",
		Long:  "主要是用来熟悉 cobra 以及 go 语言的使用，方便日后编写自己的 cli 工具",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello msctl")
		},
	}
	// 绑定可选参数配置
	flags := cmds.PersistentFlags()
	flags.StringVar(&info, "info", "basic description", "")

	groups := CommandGroups{
		{
			Message: "Basic Commands Proxy:",
			Commands: []*cobra.Command{
				exec.NewCmdExec(in, out, errout),
			},
		},
		{
			Message: "New Project Cli",
			Commands: []*cobra.Command{
				new.NewCmdNew(in, out, errout),
			},
		},
	}
	groups.Add(cmds)

	return cmds
}

func Execute() {
	if err := NewDefaultMsCliCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

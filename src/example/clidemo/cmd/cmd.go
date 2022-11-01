package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
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
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	return cmds
}

func Execute() {
	if err := NewDefaultMsCliCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

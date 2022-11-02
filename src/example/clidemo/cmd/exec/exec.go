package exec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func NewCmdExec(in io.Reader, out, errout io.Writer) *cobra.Command {
	execCmd := &cobra.Command{
		Use:   "exec (command) [flags] -- COMMAND [args...]",
		Short: "exec 转发系统已安装的命令",
		Long:  "exec 目前只支持转发系统已安装的命令，后续会增加其它的子命令",
		Run:   run,
	}

	return execCmd
}

func run(cmd *cobra.Command, args []string) {
	fd := exec.Command(args[0], args[1:]...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	if err := fd.Run(); err != nil {
		fmt.Println(err)
	}
	bf := new(bytes.Buffer)
	bf.WriteTo(fd.Stdout)
	fmt.Printf("执行结果：\n\r %s", bf.String())
}

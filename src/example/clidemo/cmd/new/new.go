package new

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/marsonshine/mscli/util"
	"github.com/spf13/cobra"
)

var (
	repository  string // 模板项目路径地址
	projectName string // 项目名称
	output      string // 输出路径
)

func NewCmdNew(in io.Reader, out, errout io.Writer) *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new (command) [flags] -- COMMAND [args...]",
		Short: "new 创建新的项目结构，mscli new -r http://repository/template",
		Long:  "new 创建新的项目结构，mscli new -r http://repository/template",
		Run:   run,
	}

	initOptionFlags(newCmd)

	return newCmd
}

func run(cmd *cobra.Command, args []string) {
	// fd := exec.Command(args[0], args[1:]...)
	// fd.Stdout = os.Stdout
	// fd.Stderr = os.Stderr
	if repository != "" {
		gitclone(repository, output, projectName)
		begin := strings.LastIndexFunc(repository, func(r rune) bool {
			return r == rune('/')
		})
		end := strings.LastIndexFunc(repository, func(r rune) bool {
			return r == rune('.')
		})
		folderName := repository[begin+1 : end]
		// rename
		renameGitFolder(folderName, projectName)
		// rename mod name
		renameProjectPackageName(projectName)
	}

	// bf := new(bytes.Buffer)
	// bf.WriteTo(fd.Stdout)
	// fmt.Printf("执行结果：\n\r %s", bf.String())
}

func initOptionFlags(newCmd *cobra.Command) {
	newCmd.Flags().StringVarP(&repository, "repository", "r", "", "mscli new -r http://repository/template")
	newCmd.Flags().StringVarP(&projectName, "project", "p", "", "mscli new -r http://repository/template -p projectName")
	newCmd.Flags().StringVarP(&output, "output", "o", "", "mscli new -r http://repository/template -p projectName -o YourLocalPath")

	newCmd.MarkFlagRequired("repository")
	newCmd.MarkFlagRequired("project")
}

func gitclone(repository, output, targetFolder string) {
	cloneToFolder := output
	if cloneToFolder != "" {
		currentFolder, _ := os.Getwd()
		cloneToFolder += "/" + currentFolder
	}
	gitCmd := exec.Command("git", "clone", repository)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	if err := gitCmd.Run(); err != nil {
		printError(gitCmd)
	} else {
		bf := new(bytes.Buffer)
		bf.WriteTo(gitCmd.Stdout)
		fmt.Println(bf.String())
		fmt.Println("git clone finished...")
	}
}

func printError(cmd *exec.Cmd) {
	bf := new(bytes.Buffer)
	bf.WriteTo(cmd.Stderr)
	fmt.Printf("错误信息：\n\r %s", bf.String())
}

func renameGitFolder(folderName string, newFolderName string) {
	if exist, _ := util.Exists(folderName); !exist {
		fmt.Println("重命名失败，文件名不存在")
		return
	}
	os.Rename(folderName, newFolderName)
}

func renameProjectPackageName(projectName string) {
	// TODO
}

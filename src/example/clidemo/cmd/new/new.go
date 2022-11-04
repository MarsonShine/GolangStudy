package new

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	if err := os.Rename(folderName, newFolderName); err != nil {
		fmt.Println(err)
	}
}

func renameProjectPackageName(projectName string) {
	// 首先要找出 mod 的包名
	var modName string
	files, err := ioutil.ReadDir(projectName)
	if err != nil {
		fmt.Printf("包重命名失败：%v", err)
		return
	}
	for _, file := range files {
		if file.Name() == "go.mod" {
			modName = readModName(path.Join(projectName, file.Name()))
			break
		}
	}
	// 获取所有目标文件 .go
	// 递归查询所有文件和文件夹
	_ = filepath.Walk(projectName, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("包重命名失败：%v", err)
			return err
		}
		if !info.IsDir() && (filepath.Ext(info.Name()) == ".go") {
			// TODO packageName 变量化，目前是与projectName一致
			renamePackageName(path, modName, projectName)
			return nil
		}
		return nil
	})
}

func renamePackageName(filepath, modName, packageName string) {
	file, _ := os.Open(filepath)
	scanner := bufio.NewScanner(file)
	newBuffer := make([]byte, len(scanner.Bytes()))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "import") {
			newBuffer = append(newBuffer, bytes.ReplaceAll(scanner.Bytes(), []byte(modName), []byte(packageName))...)
			newBuffer = append(newBuffer, []byte("\n")...)
			// 找到import结束范围
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), ")") {
					break
				}
				newBuffer = append(newBuffer, bytes.ReplaceAll(scanner.Bytes(), []byte(modName), []byte(packageName))...)
				newBuffer = append(newBuffer, []byte("\n")...)
				continue
			}
			continue
		}
		newBuffer = append(newBuffer, scanner.Bytes()...)
		newBuffer = append(newBuffer, []byte("\n")...)
	}
	os.WriteFile(filepath, newBuffer, 0644)
}

func rewrite(filepath, modName, packageName string) {
	file, _ := os.Open(filepath)
	fi, _ := file.Stat()
	scanner := bufio.NewScanner(file)
	newBuffer := make([]byte, 0, fi.Size())
	for scanner.Scan() {
		newBuffer = append(newBuffer, scanner.Bytes()...)
		newBuffer = append(newBuffer, []byte("\n")...)
	}
	os.WriteFile(filepath, newBuffer, 0644)

	fmt.Printf("origin: %d;  new: %d", fi.Size(), len(newBuffer))
}

func readModName(filepath string) string {
	file, _ := os.Open(filepath)
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		modName := scanner.Text()
		return strings.Trim(strings.ReplaceAll(modName, "module", ""), " ") // only go1.18+
	}
	return ""
}

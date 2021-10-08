package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	p := filepath.Join("dir1", "dir2", "filename")
	fmt.Println("p:", p)

	fmt.Println(filepath.Join("dir//", "filename"))
	// Join 还会删除多余的分隔符和目录
	fmt.Println(filepath.Join("dir1/../dir1", "filename"))

	fmt.Println("Dir(p):", filepath.Dir(p))
	fmt.Println("Base(p):", filepath.Base(p))
	dir, file := filepath.Split(p)
	fmt.Println("Split(p): dir,file:", dir, file)

	// 获取文件路径的拓展名
	filename := "config.json"
	ext := filepath.Ext(filename)
	fmt.Println(ext)
	// 清除拓展名之后的文件名
	fmt.Println(strings.TrimSuffix(filename, ext))

	// Rel 寻找 basepath 与 targpath 之间的相对路径
	rel, err := filepath.Rel("a/b", "a/b/t/file")
	if err != nil {
		panic(err)
	}
	fmt.Println(rel)
}

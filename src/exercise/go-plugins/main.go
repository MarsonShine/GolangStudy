package main

import (
	"fmt"
	"go-plugins/content"
	"go-plugins/plugin"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	var pm plugin.Manager
	if err := pm.LoadPlugins("./plugin-binaries"); err != nil {
		log.Fatal("loading plugins:", err)
	}
	defer pm.Close()
	// 读取文件内容
	fileInfo, err := os.Open("./sampletext.txt")
	if err != nil {
		panic(err)
	}
	// 从stdin读取内容
	contents, err := io.ReadAll(fileInfo)
	if err != nil {
		log.Fatal("reading stdin:", err)
	}

	post := &content.Post{
		Id:       42,
		Author:   "joe",
		Contents: string(contents),
	}

	fmt.Printf("=== Text before htmlize:\n%s\n", post.Contents)
	result := htmlize(&pm, post)
	fmt.Printf("\n=== Text after htmlize:\n%s\n", result)
}

var rolePattern = regexp.MustCompile(":(\\w+):`([^`]*)`")

// htmlize将post.Contents的文本变成HTML并返回；它使用插件管理器对内容和其中的角色调用已加载的插件。
func htmlize(pm *plugin.Manager, post *content.Post) string {
	pcontents := insertParagraphs(post.Contents)
	pcontents = pm.ApplyContentHooks(pcontents, post)

	return rolePattern.ReplaceAllStringFunc(pcontents, func(match string) string {
		subm := rolePattern.FindStringSubmatch(match)
		if len(subm) != 3 {
			panic("expect match")
		}
		roleName := subm[1]
		roleText := subm[2]

		roleHtml, err := pm.ApplyRoleHooks(roleName, roleText, post)
		if err == nil {
			return roleHtml
		} else {
			return match
		}
	})
}

func insertParagraphs(s string) string {
	var b strings.Builder
	p := 0
	for p < len(s) {
		// Loop invariant: p is the index of the beginning of the current paragraph.
		var paragraph string
		nextBreak := strings.Index(s[p:], "\n\n")
		if nextBreak >= 0 {
			nextBreak = p + nextBreak
			paragraph = s[p:nextBreak]
			p = nextBreak + 1
		} else {
			paragraph = strings.TrimSpace(s[p:])
			p = len(s)
		}

		paragraph = strings.Join(strings.Split(paragraph, "\n"), " ")
		b.WriteString("<p>")
		b.WriteString(paragraph)
		b.WriteString("</p>\n\n")

		// Re-point p to the start of the next parapraph.
		for p < len(s) && s[p] == '\n' {
			p++
		}
	}
	return b.String()
}

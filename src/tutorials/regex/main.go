package main

import (
	"bytes"
	"fmt"
	"regexp"
)

func main() {
	match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
	fmt.Println(match)

	r, _ := regexp.Compile("p([a-z]+)ch")
	fmt.Println(r.FindString("peach punch"))
	fmt.Println(r.FindStringIndex("peach punch"))
	// 匹配字符串以及字串
	fmt.Println(r.FindStringSubmatch("peach punch"))
	fmt.Println(r.FindStringSubmatchIndex("peach punch"))

	// 匹配全部的目标字符串匹配项
	fmt.Println(r.FindAllString("peach punch pinch", -1))
	// 限制成功匹配的次数
	fmt.Println(r.FindAllString("peach punch pinch", 2))
	// 也可以用字节数组代替
	fmt.Println(r.Match([]byte("peach punch pinch")))

	r = regexp.MustCompile("p([a-z]+)ch")
	fmt.Println(r)

	// 将所有匹配的字符串变大写
	in := []byte("a peach")
	out := r.ReplaceAllFunc(in, bytes.ToUpper)
	fmt.Println(string(out))
}

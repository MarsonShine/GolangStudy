package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type Options struct {
	Id      string `json:"id,omitempty"` // omitempty，表示如果该字段是零值就不会参数序列化
	Verbose bool   `json:"verbose,omitempty"`
	Level   int    `json:"level,omitempty"`
	Power   int    `json:"power,omitempty"`
}

var jsonString string = `
{
  "id": "foobar",
  "bug": 42
}
`

func readJson() {
	jsonDeseializer := json.NewDecoder(bytes.NewReader([]byte(jsonString)))
	// jsonDeseializer.DisallowUnknownFields() // 不允许有未匹配的属性，如此次的bug属性就无法匹配Options中的字段，就会报错

	var opts Options
	if err := jsonDeseializer.Decode(&opts); err != nil {
		fmt.Println("Decode error:", err)
	} else {
		fmt.Println("deserializa result:", opts)
	}
}

func main() {
	readJson()
	readJsonWithOmitEmpty()
}

func readJsonWithOmitEmpty() {
	opts := Options{
		Id:    "baz",
		Level: 0,
	}
	r, _ := json.MarshalIndent(opts, "", " ") // 序列化之后的字符串不含Level
	fmt.Println(string(r))
}

// 配置设置默认值，在序列化的时候没有这个字段，默认配置是将该字段赋值为零值
// 但实际上我们可能想设置默认值为非零，如10，或字符串为”MarsonShine“呢？

func parseOptions(jsn []byte) Options {
	opts := Options{
		Verbose: false,
		Level:   0,
		Power:   10,
	}
	if err := json.Unmarshal(jsn, &opts); err != nil {
		log.Fatal(err)
	}
	return opts
}

// 自定义序列化
func (o *Options) UnmarshalJson(text []byte) error {
	type options Options
	opts := options{
		Power: 10,
	}
	if err := json.Unmarshal(text, &opts); err != nil {
		return err
	}
	*o = Options(opts)
	return nil
}

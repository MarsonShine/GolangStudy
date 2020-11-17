# 设置代理：

```cmd
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct
```

# 运行步骤

1. 打开项目根目录
2. 初始化mod `go mod init projectName`
3. 获取依赖包 `go get -u gorm.io/gorm`；`go get -u gorm.io/driver/sqlite`

## SQLite 运行必备的环境

1. 安装 gcc（gcc 版本一定要最新，否则可能无法运行）
2. 安装 sqlite（Linux 自带，如没有自带则自行官网下载）

# Go 目前发现的一些注意事项

接收前端传过来的 json 请求体并序列化：

```go
http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u User
	err := decoder.Decode(&u)
    if err != nil {
        panic("序列化失败")
    }
	//访问反序列化之后的实体 u.Name
}
```

**序列化的对象属性一定要大写**，如果是小写，序列化的时候会丢失属性信息

go 的实现接口，只需要包含对应接口所有的方法签名即可。与其它语言（如 java，c#）不一样，后者这些语言都必须要显式的用关键字 `TImplemenmt:Interface` 或 `TImplement implement Interface` 标明是实现类。
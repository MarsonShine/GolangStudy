# GolangStudy
golang自学仓库

《Go 语言程序设计》：https://books.studygolang.com/gopl-zh/

《Go 语言底层设计》：https://draveness.me/golang/

《Go 编程模式》：https://coolshell.cn/articles/21128.html（陈皓）

《Go 编程风格指南》：https://github.com/uber-go/guide/blob/master/style.md

《Go 编程风格指南》- 中译版：https://github.com/xxjwxc/uber_go_guide_cn

《Go-advice》：https://github.com/cristaloleg/go-advice/blob/master/README_ZH.md

《Go by Example 中文版》：https://gobyexample-cn.github.io/

# Go 项目部署

如果项目是在 windows 下编写的，则需要交叉编译

```ini
GOOS=linux GOARCH=amd64 go build flysnow.org/hello // gitbash

// cmd
go env -w GOOS=linux 
go build .
```

## 关于 supervisor 脚本部署

```bash
vim /etc/supervisor/conf.d/kratos_demo.conf
```

```ini
#/etc/supervisor/conf.d/kratos_demo.conf
[program:kratos_demo]
directory=/root/go/src/sites/kratos-demo/cmd
command=/root/go/src/sites/kratos-demo/cmd/cmd -conf ../configs 
autostart=true
autorestart=true
stderr_logfile=/var/log/kratos-demo.err
stdout_logfile=/var/log/kratos-demo.log
environment=CODENATION_ENV=prod
environment=GOPATH="/root/go"
```

## 本地文件上传至 Linux 上传文件

在跳板机或是无法直接用 ftp 连的时候，方便起见就可以直接运行 `rz -be` 命令，这个命令只能上传文件，所以如果碰到附带文件夹的话，建议还是压缩包文件再上传。

# VSCode-GO 

## 指定 protoc 路径

```json
"protoc": {  
    "options": [
        "--proto_path=${env.GOPATH}/pkg/mod/github.com/go-kratos/kratos/v2@v2.0.0-alpha6/third_party"
    ]
},
```

## 单元测试配置

### 设置环境变量

点击 `File -> Preferences -> Settings` 添加如下节点：

```json
"go.testEnvVars": {
    "Environment":"dev"
}
```

### 取消缓存

```json
"go.testFlags": [
    "-count=1"
],
```

### 设置超时时间

```json
"go.testTimeout": "300s"
```



# 参考资料

- https://golang.org/doc/install#install
- https://medium.com/@monirz/deploy-golang-app-in-5-minutes-ff354954fa8e
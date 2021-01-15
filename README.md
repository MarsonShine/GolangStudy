# GolangStudy
golang自学仓库

《Go 语言程序设计》：https://books.studygolang.com/gopl-zh/

《Go 语言底层设计》：https://draveness.me/golang/

《Go 编程模式》：https://coolshell.cn/articles/21128.html（陈皓）

《Go 编程风格指南》：https://github.com/uber-go/guide/blob/master/style.md

《Go 编程风格指南》- 中译版：https://github.com/xxjwxc/uber_go_guide_cn



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

## 参考资料

- https://golang.org/doc/install#install
- https://medium.com/@monirz/deploy-golang-app-in-5-minutes-ff354954fa8e
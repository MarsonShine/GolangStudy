# kratos 项目介绍

官网项目地址：https://github.com/go-kratos/kratos

# 预备环境

1. [golang](https://golang.org/)
2. [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
3. [protoc](https://grpc.io/docs/protoc-installation/)
   1. 手动下载可访问 https://github.com/google/protobuf/releases 下载对应操作系统版本安装

# 安装脚手架

安装 kratos 脚手架

```go
go get -u github.com/go-kratos/kratos/tool/kratos
```

创建项目

```powershell
cd yourpath/kratosdemo
kratos new kratos-demo
```

# 运行项目

```bash
kratos run
```

成功访问 http://localhost:8000/kratos-demo/start 即说明成功
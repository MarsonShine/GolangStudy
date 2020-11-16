# 运行步骤

1. 打开项目根目录
2. 初始化mod `go mod init projectName`
3. 获取依赖包 `go get -u gorm.io/gorm`；`go get -u gorm.io/driver/sqlite`

## SQLite 运行必备的环境

1. 安装 gcc（gcc 版本一定要最新，否则可能无法运行）
2. 安装 sqlite（Linux 自带，如没有自带则自行官网下载）
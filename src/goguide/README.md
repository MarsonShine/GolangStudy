# Go 命名规范

## 包名

- 全部小写，没有大写或下划线
- 简短而简洁。请记住，在每个使用的地方都完整标识了该名称。
- 不要用复数，如是 "net/http"，而不是 "net/https"
- 名称要具体化，不要用类似 "common"、"util"、"lib"、"shared"

## 函数名

[MixedCaps](https://golang.org/doc/effective_go.html#mixed-caps)
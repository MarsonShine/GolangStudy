# kratos 项目介绍

[详见 demo 项目](kratos-demo/README.md)

# kratos 目录结构解析

## 程序入口 

入口位置在 cmd/main.go：这里对程序必要的服务做注册，核心代码为：

```go
paladin.Init()
_, closeFunc, err := di.InitApp()
if err != nil {
    panic(err)
}
```

## 服务注册

整个服务注册是在文件夹 `internal/di` 中，服务包括 6 大模块

- Redis 模块
- Memcache 模块
- DB
- 业务服务模块（internal/service）
- HTTP 模块
- GRPC 模块

当然我们可以增加额外的模块，所以要注意一定是要写在这里的

## 服务配置

服务注册还依赖各个服务的注册文件，kratos 是采用 toml 文件通过 [paladin](https://go-kratos.dev/#/config-paladin) 实现管理的，这在 `configs` 文件夹下可以清楚的看到每个模块对应着不同的 .toml 文件。

## 撸码方式

我们在编写业务代码的时候，可以参照 kratos-demo 项目已有的内容。kratos 对业务编码有严格的限制要求，业务只能通过在 `Dao` 接口类型上编写业务方法对外暴露，然后在内部实现类 `dao` 实现业务接口方法。**整个过程是对 service 层是透明的**，我们无法直接在 service 层调用数据驱动示例（db，redis 等）。具体代码路径在 `internal/dao/dao.go`

# 功能介绍

## Redis 基本使用

Redis 是缓存数据库，所以按照 kratos 的约定，我们要把缓存逻辑放到 redis.go 文件中，存取方法都是写在 `dao` 这个对象中，就如前面所说，它是派生自 `Dao`。下面是例子：

```go
// 一定要写在 dao 这个类型对象中
func (d *dao) GetDemo(c context.Context, key string) (string, error) {
	conn := d.redis.Conn(c)
	defer conn.Close()
	// 如果没有就去数据库取
	if s, _ := redis.String(conn.Do("GET", key)); s == "" {
		// 取数据库
		if _, err := d.redis.Do(c, "SET", key, "marsonshine"); err != nil {
			log.Error("conn.Set(%s) error(%v)", key, err)
		}
	}
	return redis.String(conn.Do("GET", key))
}
```

## 日志模块

kratos 内置了参照 zap 编写的 log 模块，但是由于内部限制死了输出方式：[StdoutHandler](https://github.com/go-kratos/kratos/blob/master/pkg/log/stdout.go)，目前没有做到可以自定义 handler 的目的，所以这边集成了以 [zap](https://github.com/uber-go/zap) 为基础的日志库。

> kratos v2 版本对日志模块进行了破坏性变更，代码结构也发生了很大的变化，支持了自定义渲染器 StringFormat Render 以及 handler 的注册

封装的日志库使用方式与 kratos 内置的日志库使用几乎没有区别，下面是使用方式：

1. 先注册日志中间件

```golang
// 在 internal/server/http 目录下注册中间件
middleware.UseGLogger(engine)
```

2. 初始化日志实例

```go
log := glogger.CareateLog(glogger.GLoggerConfig{})
// 一般用法
log.Info("kratos-demo start")
log.Infof("get a signal %s", s.String())
// 附加字段用法
log.With(zap.String("field1","value1")).Info("logging...")
```

日志库内部封装了对初始字段的存取：requestId, platformId, userflag 以及接口调用接口的请求运行时间。

因为要记录接口的请求日志记录时间，所以在记录日志的时候**一定要先设置请求上下文**，这样才能正确有效的记录 duration：

```go
log.SetContext(ctx).Info("接口请求记录日志...") // 此 ctx 为 context.Context 类型
```

> 注意，初始化日志实例最好设置成单例模式，不要每个请求就初始化一次日志实例

## PB 序列化 JSON 问题

因为前端对大整型数据精度的问题，只能表现出 2^53 的数据，对于 64 位长整型数据会丢失精度，所以我们要在序列化层面要自动对 int64 类型的字段转换成字符串类型。

这里分两种情况，一是 protobuffer 转换 json；二是 struct 转换 json；

解决方法相同，就是通过 go 语言内置的 tag 设置，可以显式的设置类型字段的 json 序列化，代码如下所示：

```protobuf
// protobuffer
message HelloResp {
	string Content = 1 [(gogoproto.jsontag) = 'content'];
	int64 Id = 2 [(gogoproto.jsontag) = 'id,string'];	// 这里意思是说将属性 Id 序列化成属性名为 id，类型为 string
}
```

```go
// struct
type Article struct {
	ID      int64 `json:",string"`
	Content string
	Author  string
}
```

而将不同类型的 json 转换不需要设置类型转换器 `MarshalJSON` 是因为 go 内置对这几种特殊格式的转换功能，其 json mapping 详见 https://developers.google.com/protocol-buffers/docs/proto3#json。

# 未完待续...
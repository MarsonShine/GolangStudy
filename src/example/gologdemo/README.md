# Go 日志选型

Go 的日志类库有很多选择，社区活跃度与使用度很广的都有 [logrus](https://github.com/sirupsen/logrus)、[zap](https://github.com/uber-go/zap)、[zerolog](https://github.com/rs/zerolog)、[apex/log](https://github.com/apex/log)。

其中 logrus 星数最多，关于国内者的使用资料是最多的。功能全而强大，但是有一个很明显的问题，那就是性能问题。在各个日志包的benchmark 测试报告中 logrus 性能在末位。

zap 是 uber 出品，社区也很活跃，国内使用者能查到的资料也比较常见，并且性能很高，内存友好，号称零分配。

zerolog 与 zap 类似，也是优先性能，零分配日志组件。



> 有意思的是，zap 和 zerolog 这两个库的性能测试报告中都号称超过了对方。

# Zap

本文采用 zap 做日志模块开发核心组件。

项目地址：https://github.com/uber-go/zap

本示例 demo 项目结构

```
--logcore		// 日志核心模块，待封装
--src			
--tmp
------logs
main.go			// 示例代码
```

## 基本使用

```go
// 设置输出日志文件
log.SetLevel(log.TraceLevel)
file, err := os.OpenFile("out.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
if err == nil {
    log.SetOutput(file)
}
defer file.Close()
// 记录日志
log.Info("log success")
// 记录额外地属性字段值
fields := log.Fields{"userId": 12, "requestId":"123456789"}
log.WithFields(fields).Info("log success")
```

zap 日志输出默认的格式以 `key=value` 的形式拼接的，如果想以 json 形式输出则需要设置格式化器：

```go
// 设置序列化方式
log.SetFormatter(&log.JSONFormatter{
    FieldMap: log.FieldMap{
        log.FieldKeyTime: "@timestamp",
        log.FieldKeyMsg:  "message",
    },
})
```

## 自定义日志配置

zap 还可以初始化日志配置参数：

```go
cfg := zap.Config{
    Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
    Development: true,
    Encoding:    "json",
    EncoderConfig: zapcore.EncoderConfig{
        TimeKey:        "t",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "msg",
        StacktraceKey:  "trace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     formatEncodeTime,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    },
    OutputPaths:      []string{"stdout", "./tmp/logs"},
    ErrorOutputPaths: []string{"stderr"},
    InitialFields: map[string]interface{}{
        "app": "test",
    },
}
logger, err := cfg.Build()
if err != nil {
    panic(err)
}
defer logger.Sync()
logger.Info("logger construction succeeded")
logger.Error("logger construction falied")
```

这里我挑几个重点参数讲，这与后面我们统一日志格式必需要用到的

- Level，设置记录日志的最低级别
- Encoding，转化的格式
- OutputPaths，输出的路径，第一个参数为 key，后一个参数对应具体的日志，zap 内置了两个基本标准输出，`stdout` 和 `stderr`
- InitialFields，初始化字段，在每次输出日志时，都会输出这个配置的属性值

## 高级应用

因为要统一日志格式，所以这无关语言。而本身日志格式没有统一规范，所以不同的日志组件其内部定义的输出格式是不一致的。必须要使用统一的格式。而公司内部定义的格式很可能组件时无法满足需求的，所以必须要实现自定义格式化器。

遗憾的是，zap 不支持自定义连接符，具体原因详见 [#825](https://github.com/uber-go/zap/issues/825)

解决方法很复杂，因为本身 zap 自身内置了两种格式：console 和 json。而当要增加一个格式化器需要调用 `RegisterEncoder` 函数来注册格式转码器进而实现接口 `zapcore.Encoder` 的方法 `EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error)`。具体相关知识详见 [#829](https://github.com/uber-go/zap/issues/829)，[#563](https://github.com/uber-go/zap/issues/563)

具体方案是拷贝[json_encodr.go 源码](https://github.com/uber-go/zap/blob/master/zapcore/json_encoder.go)，将其接口方法 `EncodeEntry` 中的连接符改成自己想要的连接符以及在将私有方法 `addKey` 以及 `addElementSeparator` 改成自己想要的连接符。

最后注册：

```go
zap.RegisterEncoder("key=value", keyValueEncoder)

func keyValueEncoder(c zapcore.EncoderConfig) (zapcore.Encoder, error) {
	return encoder.NewKVEncoder(c), nil
}
```

具体的自定义转码器详见代码：/logcore/encoder/kvencoder.go

## 推送 Logs 至 ElasticSearch

Go 日志推送 ElasticSearch 与 .NET 的日志库不同，Go 需要借助外力来实现才将日志输出到 ElasticSearch 的过程简化。而 .NET 的日志库是已经写好了提供给用户使用。

### FileBeat

在写好日志输出文件的功能之后，我们如何才能将日志内容输出到 ElasticSearch？这需要借助 Elastic 轻量化插件 —— FileBeat。简单来说这个插件有两个组成部分：prospector（采集器）和 harvesters（收割器）。这两个组件一起将您的日志文件输出到指定的输出（elasticsearch）。

harvesters 会扫描日志文件进行行读取，并将内容输出到 output 中。每个文件都会启动一个 harvesters，它负责打开文件和关闭文件。

FileBeat 工作原理：频繁地把文件状态（上一次 harvesters 读取问价内容时的位置）从注册表里更新到磁盘，然后保证能读取出来发送给 output。如果 output 不可用（如 elasticsearch 或 logstash 宕机不可用），FileBeat 则会将最近读取的位置信息保存下来，等到 output 可用时在恢复发送。

这里就会有一个潜在的问题：因为每个日志文件就会开启一个 prospector 和 harvesters，那么每天会产生大量的新日志文件，那么就会发现 FileBeat 的注册表文件会变得非常大，这个场景自然 elastic 也想到并提供了解决方案，详见：https://www.elastic.co/guide/en/beats/filebeat/current/reduce-registry-size.html

### 集成 FileBeat

日志系统所属服务器安装 [FileBeat](https://www.elastic.co/cn/downloads/beats/filebeat)。解压并更新 filebeat.yml 文件：

```yml
type: log
# 更改为 true 表示应用此配置节点
enabled: true
paths:
  - e:\repositories\GolangStudy\src\example\gologdemo\tmp\*
# 输出
output.elasticsearch:
  hosts: ["192.168.3.67:9200"]
  username: "elastic"
  password: "changeme"
```

然后启动 filebeat，切换到安装目录执行以下命令：

```bash
.\filebeat.cmd -e -c .\filebeat.yml
```

这样就完成了对日志内容输出到 ElasticSearch，业务代码不需要做任何更改。

> 当然，FileBeat 是中间件，可以使用在 Go 也能为 .NET 使用

## 推送 Logs 至 Logstash

方法与 [推送 Logs 至 ElasticSearch](#推送 Logs 至 ElasticSearch) 一样，采用节点 `output.logstash` 即可。

# 封装平台 Log 包

//TODO
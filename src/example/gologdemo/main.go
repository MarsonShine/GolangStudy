package main

// https://logz.io/blog/golang-logs/
import (
	"bytes"
	"fmt"
	"golog/logcore/encoder"
	"golog/src/fzlog"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// logrusExample()
	// zaplogExample()
	fzlogExample()
}

func logrusExample() {
	// log.WithFields(log.Fields{
	// 	"animal": "walrus",
	// }).Info("A walrus appears")
	// 设置序列化方式
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime: "@timestamp",
			log.FieldKeyMsg:  "message",
		},
	})
	log.SetLevel(log.TraceLevel)
	file, err := os.OpenFile("out.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	defer file.Close()

	fields := log.Fields{"userId": 12, "requestId": "123456789"}
	log.WithFields(fields).Info("user logged in")
}

func zaplogExample() {
	// basicZaplogExample()
	addCustomEncoder()
	// registerConfig()
	basicConfigurationExample()
}

func fzlogExample() {
	logger := fzlog.CreateLog()
	logger.Info("fzlog Info...")
	logger.Debug("fzlog Debug...")
	logger.Error("fzlog Error...")
}

func basicZaplogExample() {
	const url = "example.com"
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar() // 在性能不错但是还不是关键因素的情况下用 sugar，如果性能遇到瓶颈，可以直接调用 logger
	sugar.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}

func basicConfigurationExample() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "key=value",
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
}

func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}

func addCustomEncoder() {
	zap.RegisterEncoder("key=value", keyValueEncoder)
}

func registerConfig() {
	buf := bytes.NewBuffer(nil)
	zap.RegisterSink("filebeat", func(u *url.URL) (zap.Sink, error) {
		return privateServerSink{
			zapcore.AddSync(buf),
		}, nil
	})
}

type privateServerSink struct{ zapcore.WriteSyncer }

func (privateServerSink) Close() error { return nil }

func keyValueEncoder(c zapcore.EncoderConfig) (zapcore.Encoder, error) {
	return encoder.NewKVEncoder(c), nil
}

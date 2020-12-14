package fzlog

import (
	"fmt"
	"golog/logcore/encoder"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FzLog struct {
	log *zap.Logger
}

func (fzlog FzLog) Info(msg string, fields ...zap.Field) {
	fzlog.log.Info(msg, fields...)
}

func (fzlog FzLog) Error(msg string, fields ...zap.Field) {
	fzlog.log.Error(msg, fields...)
}

func (fzlog FzLog) Debug(msg string, fields ...zap.Field) {
	fzlog.log.Debug(msg, fields...)
}

var config zap.Config

func CreateLog() FzLog {
	initDefaultConfig()
	logger, err := config.Build()
	if err != nil {
		logger.Error("logger construction falied")
		panic(err)
	}
	defer logger.Sync()
	logger.Info("logger construction succeeded")
	return FzLog{
		logger,
	}
}

func initDefaultConfig() {
	registerEncoder()
	config = zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Encoding:    "kvpare",
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
	}
}

func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}

func registerEncoder() {
	zap.RegisterEncoder("kvpare", func(c zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return encoder.NewKVEncoder(c), nil
	})
}

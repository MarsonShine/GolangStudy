package fzlog

import (
	"context"
	"fmt"
	"golog/logcore/encoder"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// std is the name of the standard logger in stdlib `log`
	std = CreateLog()
)

const (
	RequestID  string = "requestId"
	PlatformID string = "platformId"
	UserFlag   string = "userFlag"
	Duration   string = "duration"
	Size       string = "size"
)

type FzLog struct {
	log     *zap.Logger
	context context.Context
}

func (fzlog FzLog) initDefaultFields() FzLog {
	ctx := fzlog.context
	d := ctx.Value(Duration)
	r := ctx.Value(RequestID)
	p := ctx.Value(PlatformID)
	fmt.Print(d, r, p)
	start, ok := ctx.Value(Duration).(time.Time)
	var duration = ""
	if ok {
		duration = strconv.FormatInt(time.Since(start).Milliseconds(), 10)
	}

	fzlog.log.With(
		zap.String(RequestID, ctx.Value(RequestID).(string)),
		zap.String(UserFlag, ctx.Value(UserFlag).(string)),
		zap.String(PlatformID, ctx.Value(UserFlag).(string)),
		zap.String(Duration, duration),
		zap.Int64(Size, ctx.Value(Size).(int64)),
	)
	return fzlog
}

func Info(msg string, fields ...zap.Field) {
	std.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	std.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	std.Debug(msg, fields...)
}

func (fzlog FzLog) Info(msg string, fields ...zap.Field) {
	fzlog.initDefaultFields().log.Info(msg, fields...)
}

func (fzlog FzLog) With(fields ...zap.Field) FzLog {
	fzlog.initDefaultFields().log.With(fields...)
	return fzlog
}

func (fzlog FzLog) Error(msg string, fields ...zap.Field) {
	fzlog.initDefaultFields().log.Error(msg, fields...)
}

func (fzlog FzLog) Debug(msg string, fields ...zap.Field) {
	fzlog.initDefaultFields().log.Debug(msg, fields...)
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
		log:     logger,
		context: nil,
	}
}

func withContext(ctx context.Context) FzLog {
	if ctx == nil {
		panic("context is not null")
	}
	log := CreateLog()
	log.context = ctx
	return log
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
		InitialFields: map[string]interface{}{
			"requestId":  "",
			"userflag":   "",
			"platformId": "",
		},
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

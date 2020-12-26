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
	context *context.Context
	sugar   *zap.SugaredLogger
}

func (fzlog FzLog) initDefaultFields() FzLog {
	ctx := *fzlog.context
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
	// fzlog.initDefaultFields().log.Info(msg, fields...)
	fields = fzlog.appendFields(fields...)
	fzlog.log.Info(msg, fields...)
}

func (fzlog FzLog) appendFields(fields ...zap.Field) []zap.Field {
	if nil == fzlog.context {
		return nil
	}
	ctx := *fzlog.context
	start, ok := ctx.Value(Duration).(time.Time)
	var duration = ""
	if ok {
		duration = strconv.FormatInt(time.Since(start).Milliseconds(), 10)
	}
	fileds2 := []zapcore.Field{
		zap.String(RequestID, ctx.Value(RequestID).(string)),
		zap.String(UserFlag, ctx.Value(UserFlag).(string)),
		zap.String(PlatformID, ctx.Value(PlatformID).(string)),
		zap.String(Duration, duration),
		zap.Int64(Size, ctx.Value(Size).(int64)),
	}
	fields = append(fields, fileds2...)
	return fields
}

func (fzlog FzLog) With(fields ...zap.Field) FzLog {
	fzlog.log.With(fields...)
	return fzlog
}

func (fzlog FzLog) Error(msg string, fields ...zap.Field) {
	fields = fzlog.appendFields(fields...)
	fzlog.log.Error(msg, fields...)
}

func (fzlog FzLog) Debug(msg string, fields ...zap.Field) {
	fields = fzlog.appendFields(fields...)
	fzlog.log.Debug(msg, fields...)
}

func (fzlog FzLog) Infof(msg string, args ...interface{}) {
	a := fzlog.defaultLogData()
	if a != nil {
		fzlog.sugar.With(a).Infof(msg, args...)
	} else {
		fzlog.sugar.Infof(msg, args...)
	}
}
func (fzlog FzLog) Errorf(msg string, args ...interface{}) {
	a := fzlog.defaultLogData()
	if a != nil {
		fzlog.sugar.With(a).Errorf(msg, args...)
	} else {
		fzlog.sugar.Errorf(msg, args...)
	}
}
func (fzlog FzLog) Warnf(msg string, args ...interface{}) {
	a := fzlog.defaultLogData()
	if a != nil {
		fzlog.sugar.With(a).Warnf(msg, args...)
	} else {
		fzlog.sugar.Warnf(msg, args...)
	}
}

func (fzlog FzLog) defaultLogData() interface{} {
	if nil == fzlog.context {
		return nil
	}
	ctx := *fzlog.context
	start, ok := ctx.Value(Duration).(time.Time)
	var duration = ""
	if ok {
		duration = strconv.FormatInt(time.Since(start).Milliseconds(), 10)
	}
	args := struct {
		RequestID  string
		UserFlag   string
		PlatformID string
		Duration   string
		Size       int64
	}{
		ctx.Value(RequestID).(string),
		ctx.Value(UserFlag).(string),
		ctx.Value(PlatformID).(string),
		duration,
		ctx.Value(Size).(int64),
	}

	return args
}

var config zap.Config

func CreateLog() FzLog {
	initDefaultConfig()
	logger, err := config.Build()
	zap.ReplaceGlobals(logger)
	if err != nil {
		logger.Error("logger construction falied")
		panic(err)
	}
	defer logger.Sync()
	logger.Info("logger construction succeeded")
	return FzLog{
		log:     logger,
		context: nil,
		sugar:   logger.Sugar(),
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
		// InitialFields: map[string]interface{}{
		// 	"requestId":  "",
		// 	"userflag":   "",
		// 	"platformId": "",
		// },
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

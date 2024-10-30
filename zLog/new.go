package zLog

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger = MustNew()

func MustNew(opts ...Option) *zap.Logger {
	logger, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return logger
}

func New(opts ...Option) (*zap.Logger, error) {
	o := evaluateOptions(opts)
	// 一些定制化字段
	var fields []zap.Field
	for _, kv := range o.kvPairs {
		fields = append(fields, zap.String(kv.K, kv.V))
	}

	//zap原生的选项，下面指定了一些默认值，可覆盖
	zapOpts := append([]zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.DPanicLevel),
		zap.Fields(fields...),
	}, o.zapOptions...)

	// 标准输出
	writeSyncers := make([]zapcore.WriteSyncer, 0, 4)
	writeSyncers = append(writeSyncers, o.writeSyncers...)

	if len(writeSyncers) == 0 {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}

	encfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	if o.levelColor {
		encfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	tee := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encfg),
			//zapcore.NewConsoleEncoder(encfg),
			zapcore.NewMultiWriteSyncer(writeSyncers...),
			o.atomicLevel,
		),
	)
	ins := zap.New(tee, zapOpts...)
	return ins, nil
}

func SetDefaultLogger(logger *zap.Logger) {
	defaultLogger = logger
}

func GetDefaultLogger() *zap.Logger {
	return defaultLogger
}

// Deprecated: use TraceXXX eg TraceDebug.
func WithContext(ctx context.Context) *zap.Logger {
	// notice logger will be clone much times by this func
	return defaultLogger.With(TraceIdFromCtx(ctx))
}

func With(fields ...zap.Field) *zap.Logger {
	// notice logger will be clone much times by this func
	return defaultLogger.With(fields...)
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return defaultLogger.WithOptions(opts...)
}

func ErrorOrInfo(msg string, err error, fields ...zap.Field) {
	if err != nil {
		defaultLogger.Error(msg, append([]zap.Field{zap.Error(err)}, fields...)...)
	} else {
		defaultLogger.Info(msg, fields...)
	}
}

func Sync() error {
	return defaultLogger.Sync()
}

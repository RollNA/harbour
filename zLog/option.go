package zLog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	bizKey     = "biz"
	traceIdKey = "traceId"
)

type KVPair struct {
	K string
	V string
}

//"b": "api",
//"id": "3987960976",

type Option func(o *options)

type options struct {
	kvPairs []KVPair // 预设不变字段

	ctxTraceKeyName   interface{} // ctx中链路追踪的key
	ctxContextKeyName interface{} // ctx中上下文的key

	atomicLevel zap.AtomicLevel
	levelColor  bool
	zapOptions  []zap.Option

	logPath      string
	writeSyncers []zapcore.WriteSyncer

	rotate *lumberjack.Logger
}

func defaultOptions() options {
	return options{
		atomicLevel: zap.NewAtomicLevelAt(zap.InfoLevel),
		rotate: &lumberjack.Logger{
			Filename:   "./logs/server.log",
			MaxSize:    20,
			MaxAge:     1,
			MaxBackups: 100,
			LocalTime:  false,
			Compress:   false,
		},
	}
}

func evaluateOptions(opts []Option) options {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func OptLogKVPair(kvPairs []KVPair) Option {
	return func(o *options) {
		if kvPairs != nil {
			o.kvPairs = kvPairs
		}
	}
}

func OptCtxTraceKeyName(key interface{}) Option {
	return func(o *options) {
		o.ctxTraceKeyName = key
	}
}

func OptCtxContextKeyName(key interface{}) Option {
	return func(o *options) {
		o.ctxContextKeyName = key
	}
}

// 日志等级颜色输出
func OptLevelColor(levelColor bool) Option {
	return func(o *options) {
		o.levelColor = levelColor
	}
}

// 最低日志输出等级
func OptLevel(level zapcore.Level) Option {
	return func(o *options) {
		o.atomicLevel = zap.NewAtomicLevelAt(level)
	}
}

// 同一日志多个io输出
func OptWriteSyncers(ws ...zapcore.WriteSyncer) Option {
	return func(o *options) {
		o.writeSyncers = ws
	}
}

func OptZapOptions(opts ...zap.Option) Option {
	return func(o *options) {
		o.zapOptions = opts
	}
}

func OptFilePath(filePath string) Option {
	return func(o *options) {
		o.logPath = filePath
	}
}

func OptRotate(rotate *lumberjack.Logger) Option {
	return func(o *options) {
		o.rotate = rotate
	}
}

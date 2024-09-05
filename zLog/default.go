package zLog

import "go.uber.org/zap"

func Debug(msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	defaultLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	defaultLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	defaultLogger.Fatal(msg, fields...)
}

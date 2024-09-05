package zLog

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
)

func SpanIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}
	return ""
}

func TraceIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return ""
}

func TraceIdFromCtx(ctx context.Context) zap.Field {
	traceId := TraceIdFromContext(ctx)
	if traceId == "" {
		return zap.Skip()
	}
	return TraceId(traceId)
}

func TraceId(traceId string) zap.Field {
	return zap.String("traceId", traceId)
}

func TraceDebug(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, append(fields, TraceIdFromCtx(ctx))...)
}

func TraceInfo(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, append(fields, TraceIdFromCtx(ctx))...)
}

func TraceWarn(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, append(fields, TraceIdFromCtx(ctx))...)
}

func TraceError(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, append(fields, TraceIdFromCtx(ctx))...)
}

func TraceDPanic(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.DPanic(msg, append(fields, TraceIdFromCtx(ctx))...)
}

func TracePanic(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Panic(msg, append(fields, TraceIdFromCtx(ctx))...)
}

func TraceFatal(ctx context.Context, msg string, fields ...zap.Field) {
	defaultLogger.Fatal(msg, append(fields, TraceIdFromCtx(ctx))...)
}

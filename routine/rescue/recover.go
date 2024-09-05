package rescue

import (
	"context"

	"github.com/quincy0/harbour/zLog"
	"go.uber.org/zap"
)

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		zLog.Error("recover failed", zap.Any("p", p), zap.Stack("stack"))
	}
}

func RecoverCtx(ctx context.Context, cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}
	if p := recover(); p != nil {
		zLog.TraceError(ctx, "recover failed", zap.Any("p", p), zap.Stack("stack"))
	}
}

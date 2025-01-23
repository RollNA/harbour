package routine

import (
	"context"
	"time"

	"github.com/RollNA/harbour/routine/rescue"
	"go.opentelemetry.io/otel/trace"
)

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func GoSafeCtx(ctx context.Context, fn func()) {
	go RunSafeCtx(ctx, fn)
}

func RunSafe(fn func()) {
	defer rescue.Recover()
	fn()
}

func RunSafeCtx(ctx context.Context, fn func()) {
	defer rescue.RecoverCtx(ctx)
	fn()
}

func NewContextWithTimeout(c context.Context, seconds int64) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(seconds))
	span := trace.SpanFromContext(c)
	ctx = trace.ContextWithSpan(ctx, span)
	// copy c value
	return ctx, cancel
}

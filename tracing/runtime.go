package tracing

import (
	"context"
	"runtime/trace"
)

var _ Tracer = (*RuntimeTracer)(nil)

type RuntimeTracer struct{}

func (r *RuntimeTracer) WithRegion(ctx context.Context, regionName string, fn func()) {
	trace.WithRegion(ctx, regionName, fn)
}

func (r *RuntimeTracer) Logf(ctx context.Context, category, format string, args ...any) {
	trace.Logf(ctx, category, format, args...)
}

func (r *RuntimeTracer) NewTask(ctx context.Context, taskType string) (context.Context, Task) {
	return trace.NewTask(ctx, taskType)
}

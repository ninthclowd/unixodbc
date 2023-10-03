package tracing

import (
	"context"
)

type Tracer interface {
	NewTask(ctx context.Context, taskType string) (context.Context, Task)
	WithRegion(ctx context.Context, regionName string, fn func())
	Logf(ctx context.Context, category, format string, args ...any)
}

type Task interface {
	End()
}

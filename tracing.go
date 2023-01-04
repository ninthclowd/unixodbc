package unixodbc

import (
	"context"
	"runtime/trace"
)

var tracer Tracer = new(RuntimeTracer)

type Tracer interface {
	Start(ctx context.Context, name string) Trace
}

type Trace interface {
	Start(name string) Trace
	End()
}

var _ Tracer = (*NoOpTracer)(nil)

type NoOpTracer struct{}

func (n *NoOpTracer) Start(ctx context.Context, name string) Trace { return new(NoOpTrace) }

var _ Trace = (*NoOpTrace)(nil)

type NoOpTrace struct{}

func (n *NoOpTrace) Start(name string) Trace { return new(NoOpTrace) }

func (n *NoOpTrace) End() {}

var _ Tracer = (*RuntimeTracer)(nil)

type RuntimeTracer struct{}

func (g *RuntimeTracer) Start(ctx context.Context, name string) Trace {
	ctx2, t := trace.NewTask(ctx, name)
	return &GoTrace{ctx2, t}
}

var _ Trace = (*GoTrace)(nil)

type GoTrace struct {
	ctx context.Context
	t   *trace.Task
}

func (g *GoTrace) Start(name string) Trace {
	ctx, t := trace.NewTask(g.ctx, name)
	return &GoTrace{ctx, t}
}

func (g *GoTrace) End() {
	g.t.End()
}

package unixodbc

import (
	"github.com/ninthclowd/unixodbc/tracing"
)

var Tracer tracing.Tracer = new(tracing.RuntimeTracer)

package zaptrace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Context constructs traceId and spanId field from context.Context if a
// trace.SpanContext is present in the context value.
func Context(ctx context.Context) zap.Field {
	spanCtx := trace.SpanContextFromContext(ctx)

	return SpanContext(spanCtx)
}

type spanContextMarshaler struct {
	SpanContext trace.SpanContext
}

func (m spanContextMarshaler) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if m.SpanContext.HasTraceID() {
		enc.AddString("traceId", m.SpanContext.TraceID().String())
	}

	if m.SpanContext.HasSpanID() {
		enc.AddString("spanId", m.SpanContext.SpanID().String())
	}

	return nil
}

// SpanContext constructs traceId and spanId field from a trace.SpanContext if
// it has valid traceId and spanId.
func SpanContext(spanCtx trace.SpanContext) zap.Field {
	return zap.Inline(spanContextMarshaler{
		SpanContext: spanCtx,
	})
}

package zaptrace_test

import (
	"context"
	"crypto/rand"
	"testing"

	rec "github.com/adzil/zapf/internal/fieldrecorder"
	"github.com/adzil/zapf/zaptrace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestTrace(t *testing.T) {
	type Context struct {
		Field        zap.Field
		AssertResult func(result rec.Object)
	}

	for k, v := range map[string]func(t *testing.T, tc *Context){
		"with no span context": func(t *testing.T, tc *Context) {
			tc.Field = zaptrace.Context(context.Background())

			tc.AssertResult = func(result rec.Object) {
				assert.Len(t, result, 0, "encoded result should be empty")
			}
		},

		"with randomized span context value": func(t *testing.T, tc *Context) {
			config := trace.SpanContextConfig{}

			_, err := rand.Read(config.TraceID[:])
			require.NoError(t, err, "rand read should return nil error")

			_, err = rand.Read(config.SpanID[:])
			require.NoError(t, err, "rand read should return nil error")

			spanCtx := trace.NewSpanContext(config)
			ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)

			tc.Field = zaptrace.Context(ctx)

			expected := rec.Object{
				"traceId": rec.String(config.TraceID.String()),
				"spanId":  rec.String(config.SpanID.String()),
			}

			tc.AssertResult = func(result rec.Object) {
				assert.Equal(t, expected, result, "trace and span id should match")
			}
		},
	} {
		t.Run(k, func(t *testing.T) {
			tc := &Context{}
			v(t, tc)

			assert.Equal(t, zapcore.InlineMarshalerType, tc.Field.Type)

			om, ok := tc.Field.Interface.(zapcore.ObjectMarshaler)
			require.True(t, ok, "field interface must be an object marshaler")

			enc := rec.NewObjectEncoder(t)
			err := om.MarshalLogObject(enc)

			assert.NoError(t, err, "marshal log object should return nil error")
			tc.AssertResult(enc.Result())
		})
	}
}

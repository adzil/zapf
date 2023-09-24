package zapproto_test

import (
	"testing"

	rec "github.com/adzil/zapf/internal/fieldrecorder"
	marshalerpb "github.com/adzil/zapf/internal/gen/go/marshaler"
	"github.com/adzil/zapf/zapproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	marshalerpbMessageFullName = (&marshalerpb.Message{}).ProtoReflect().Descriptor().FullName()
)

func TestMessage_Marshal(t *testing.T) {
	type Context struct {
		Input     zap.Field
		FieldName string
		Expects   rec.Object
	}

	for k, v := range map[string]func(t *testing.T, tc *Context){
		"Message": func(t *testing.T, tc *Context) {
			tc.Input = zapproto.Message(tc.FieldName, &marshalerpb.Message{
				Text: "hello",
			})

			tc.Expects = rec.Object{
				"text": rec.String("hello"),
			}
		},

		"TypedMessage": func(t *testing.T, tc *Context) {
			tc.Input = zapproto.TypedMessage(tc.FieldName, &marshalerpb.Message{
				Text: "hello",
			})

			tc.Expects = rec.Object{
				"@type": rec.String("type.googleapis.com/" + marshalerpbMessageFullName),
				"text":  rec.String("hello"),
			}
		},
	} {
		t.Run(k, func(t *testing.T) {
			tc := &Context{
				FieldName: "message",
			}
			v(t, tc)

			assert.Equal(t, tc.FieldName, tc.Input.Key)
			assert.Equal(t, zapcore.ObjectMarshalerType, tc.Input.Type)

			om, ok := tc.Input.Interface.(zapcore.ObjectMarshaler)
			require.True(t, ok, "field should have object marshaler set")

			enc := rec.NewObjectEncoder(t)
			err := om.MarshalLogObject(enc)

			assert.NoError(t, err, "marshal log object should return nil error")
			assert.Equal(t, tc.Expects, enc.Result(), "encoded object should match")
		})
	}
}

func TestMessages_Marshal(t *testing.T) {
	type Context struct {
		Input     zap.Field
		FieldName string
		Expects   rec.Array
	}

	for k, v := range map[string]func(t *testing.T, tc *Context){
		"Messages": func(t *testing.T, tc *Context) {
			tc.Input = zapproto.Messages(tc.FieldName, []*marshalerpb.Message{
				{
					Text: "hello",
				},
				{
					Text: "world",
				},
			})

			tc.Expects = rec.Array{
				rec.Object{
					"text": rec.String("hello"),
				},
				rec.Object{
					"text": rec.String("world"),
				},
			}
		},

		"TypedMessages": func(t *testing.T, tc *Context) {
			tc.Input = zapproto.TypedMessages(tc.FieldName, []*marshalerpb.Message{
				{
					Text: "hello",
				},
				{
					Text: "world",
				},
			})

			tc.Expects = rec.Array{
				rec.Object{
					"@type": rec.String("type.googleapis.com/" + marshalerpbMessageFullName),
					"text":  rec.String("hello"),
				},
				rec.Object{
					"@type": rec.String("type.googleapis.com/" + marshalerpbMessageFullName),
					"text":  rec.String("world"),
				},
			}
		},
	} {
		t.Run(k, func(t *testing.T) {
			tc := &Context{
				FieldName: "message",
			}
			v(t, tc)

			assert.Equal(t, tc.FieldName, tc.Input.Key)
			assert.Equal(t, zapcore.ArrayMarshalerType, tc.Input.Type)

			am, ok := tc.Input.Interface.(zapcore.ArrayMarshaler)
			require.True(t, ok, "field should have array marshaler set")

			enc := rec.NewArrayEncoder(t)
			err := am.MarshalLogArray(enc)

			assert.NoError(t, err, "marshal log array should return nil error")
			assert.Equal(t, tc.Expects, enc.Result(), "encoded object should match")
		})
	}
}

func TestMessages_Errors(t *testing.T) {
	field := zapproto.Messages("message", []*anypb.Any{{}})

	am, ok := field.Interface.(zapcore.ArrayMarshaler)
	require.True(t, ok, "field should have array marshaler set")

	enc := rec.NewArrayEncoder(t)
	err := am.MarshalLogArray(enc)

	assert.Error(t, err, "marshal log array should return an error")
	assert.Len(t, enc.Result(), 0, "there should be no encoded result")
}

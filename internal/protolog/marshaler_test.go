package protolog_test

import (
	"encoding/base64"
	"testing"

	rec "github.com/adzil/zapf/internal/fieldrecorder"
	marshalerpb "github.com/adzil/zapf/internal/gen/go/marshaler"
	"github.com/adzil/zapf/internal/protolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestOptions_MarshalerOf_MarshalLogObject(t *testing.T) {
	type Context struct {
		Input        proto.Message
		Options      protolog.Options
		AssertErr    func(error)
		AssertObject func(rec.Object)
	}

	for k, v := range map[string]func(t *testing.T, tc *Context){
		"message with various types": func(t *testing.T, tc *Context) {
			tc.Input = &marshalerpb.Marshaler{
				Map: map[string]string{
					"hello": "world",
				},
				Array:    []string{"hello", "world"},
				Bool:     true,
				String_:  "hello",
				Bytes:    []byte(`world`),
				Enum:     marshalerpb.Choice_CHOICE_ONE,
				Float:    1.0,
				Double:   -2.0,
				Int32:    -3,
				Int64:    -4,
				Uint32:   5,
				Uint64:   6,
				Sint32:   -7,
				Sint64:   -8,
				Fixed32:  9,
				Fixed64:  10,
				Sfixed32: -11,
				Sfixed64: -12,
				Payload: &marshalerpb.Marshaler_Message{
					Message: &marshalerpb.Message{
						Text: "hello",
					},
				},
			}

			expected := rec.Object{
				"map": rec.Object{
					"hello": rec.String("world"),
				},
				"array":    rec.Array{rec.String("hello"), rec.String("world")},
				"bool":     rec.Bool(true),
				"string":   rec.String("hello"),
				"bytes":    rec.String(base64.StdEncoding.EncodeToString([]byte(`world`))),
				"enum":     rec.String(marshalerpb.Choice_CHOICE_ONE.String()),
				"float":    rec.Float32(1.0),
				"double":   rec.Float64(-2.0),
				"int32":    rec.Int32(-3),
				"int64":    rec.Int64(-4),
				"uint32":   rec.Uint32(5),
				"uint64":   rec.Uint64(6),
				"sint32":   rec.Int32(-7),
				"sint64":   rec.Int64(-8),
				"fixed32":  rec.Uint32(9),
				"fixed64":  rec.Uint64(10),
				"sfixed32": rec.Int32(-11),
				"sfixed64": rec.Int64(-12),
				"message": rec.Object{
					"text": rec.String("hello"),
				},
			}

			tc.AssertObject = func(o rec.Object) {
				assert.Equal(t, expected, o, "object should reflect filled fields")
			}
		},

		"message with any field": func(t *testing.T, tc *Context) {
			anyPayload, err := anypb.New(&marshalerpb.Message{
				Text: "hello",
			})
			require.NoError(t, err, "anypb new must return no error")

			tc.Input = &marshalerpb.Marshaler{
				Payload: &marshalerpb.Marshaler_Any{
					Any: anyPayload,
				},
			}

			expected := rec.Object{
				"any": rec.Object{
					"@type": rec.String(anyPayload.TypeUrl),
					"text":  rec.String("hello"),
				},
			}

			tc.AssertObject = func(o rec.Object) {
				assert.Equal(t, expected, o, "object should contain any")
			}
		},

		"failed unmarshal empty any": func(t *testing.T, tc *Context) {
			tc.Input = &anypb.Any{}

			tc.AssertErr = func(err error) {
				assert.ErrorContains(t, err, "empty type URL", "marshal with empty anypb.Any should return error")
			}
		},
	} {
		t.Run(k, func(t *testing.T) {
			tc := &Context{}
			v(t, tc)

			enc := rec.NewObjectEncoder(t)
			err := tc.Options.MarshalerOf(tc.Input).MarshalLogObject(enc)

			switch {
			case tc.AssertErr != nil && tc.AssertObject == nil:
				tc.AssertErr(err)

			case tc.AssertObject != nil:
				require.NoError(t, err, "marshal log object should return no error")
				tc.AssertObject(enc.Result())

			default:
				t.Error("AssertErr and AssertObj cannot be both empty nor both defined")
			}
		})
	}
}

func TestNewMarshalerOf_MarshalLogObject(t *testing.T) {
	enc := rec.NewObjectEncoder(t)
	err := protolog.MarshalerOf(nil).MarshalLogObject(nil)

	assert.NoError(t, err, "new marshaler's marshal log object should return no error")
	assert.Nil(t, enc.Result(), "enc result should be nil")
}

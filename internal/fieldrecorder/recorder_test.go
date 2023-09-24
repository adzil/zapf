package rec_test

import (
	"errors"
	"testing"

	rec "github.com/adzil/zapf/internal/fieldrecorder"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

type objectMarshalerFunc func(enc zapcore.ObjectEncoder) error

func (fn objectMarshalerFunc) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return fn(enc)
}

type arrayMarshalerFunc func(enc zapcore.ArrayEncoder) error

func (fn arrayMarshalerFunc) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	return fn(enc)
}

func TestObjectEncoder_Values(t *testing.T) {
	expected := rec.Object{
		"object": rec.Object{
			"string": rec.String("test"),
		},
		"array": rec.Array{
			rec.String("test"),
		},
		"bool":    rec.Bool(true),
		"float32": rec.Float32(1.0),
		"float64": rec.Float64(-2.0),
		"int32":   rec.Int32(-3),
		"int64":   rec.Int64(-4),
		"uint32":  rec.Uint32(5),
		"uint64":  rec.Uint64(6),
	}

	enc := rec.NewObjectEncoder(t)

	err := enc.AddObject("object", objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		enc.AddString("string", "test")

		return nil
	}))
	assert.NoError(t, err, "add array should return no error")

	err = enc.AddArray("array", arrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		enc.AppendString("test")

		return nil
	}))
	assert.NoError(t, err, "add array should return no error")

	enc.AddBool("bool", true)
	enc.AddFloat32("float32", 1.0)
	enc.AddFloat64("float64", -2.0)
	enc.AddInt32("int32", -3)
	enc.AddInt64("int64", -4)
	enc.AddUint32("uint32", 5)
	enc.AddUint64("uint64", 6)

	assert.Equal(t, expected, enc.Result(), "encode result should equal with the expected result")
}

func TestObjectEncoder_AddErrors(t *testing.T) {
	enc := rec.NewObjectEncoder(t)

	expectedError := errors.New("test error")

	err := enc.AddObject("object", objectMarshalerFunc(func(_ zapcore.ObjectEncoder) error {
		return expectedError
	}))
	assert.ErrorIs(t, err, expectedError, "add object should return expected error")

	err = enc.AddArray("array", arrayMarshalerFunc(func(_ zapcore.ArrayEncoder) error {
		return expectedError
	}))
	assert.ErrorIs(t, err, expectedError, "add object should return expected error")

	assert.Len(t, enc.Result(), 0, "encoder should be empty")
}

func TestArrayEncoder_Values(t *testing.T) {
	expected := rec.Array{
		rec.Object{
			"string": rec.String("test"),
		},
		rec.Array{
			rec.String("test"),
		},
		rec.Bool(true),
		rec.Float32(1.0),
		rec.Float64(-2.0),
		rec.Int32(-3),
		rec.Int64(-4),
		rec.Uint32(5),
		rec.Uint64(6),
	}

	enc := rec.NewArrayEncoder(t)

	err := enc.AppendObject(objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		enc.AddString("string", "test")

		return nil
	}))
	assert.NoError(t, err, "add array should return no error")

	err = enc.AppendArray(arrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		enc.AppendString("test")

		return nil
	}))
	assert.NoError(t, err, "add array should return no error")

	enc.AppendBool(true)
	enc.AppendFloat32(1.0)
	enc.AppendFloat64(-2.0)
	enc.AppendInt32(-3)
	enc.AppendInt64(-4)
	enc.AppendUint32(5)
	enc.AppendUint64(6)

	assert.Equal(t, expected, enc.Result(), "encode result should equal with the expected result")
}

func TestArrayEncoder_AddErrors(t *testing.T) {
	enc := rec.NewArrayEncoder(t)

	expectedError := errors.New("test error")

	err := enc.AppendObject(objectMarshalerFunc(func(_ zapcore.ObjectEncoder) error {
		return expectedError
	}))
	assert.ErrorIs(t, err, expectedError, "add object should return expected error")

	err = enc.AppendArray(arrayMarshalerFunc(func(_ zapcore.ArrayEncoder) error {
		return expectedError
	}))
	assert.ErrorIs(t, err, expectedError, "add object should return expected error")

	assert.Len(t, enc.Result(), 0, "encoder should be empty")
}

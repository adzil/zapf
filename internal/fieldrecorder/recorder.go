package rec

import (
	"testing"

	"github.com/adzil/zapf/internal/mocks"
	"go.uber.org/zap/zapcore"
)

type Value interface {
	isValue()
}

type (
	Object  map[string]Value
	Array   []Value
	String  string
	Bool    bool
	Float32 float32
	Float64 float64
	Int32   int32
	Uint32  uint32
	Int64   int
	Uint64  uint
)

func (Object) isValue()  {}
func (Array) isValue()   {}
func (String) isValue()  {}
func (Bool) isValue()    {}
func (Float32) isValue() {}
func (Float64) isValue() {}
func (Int32) isValue()   {}
func (Uint32) isValue()  {}
func (Int64) isValue()   {}
func (Uint64) isValue()  {}

type ObjectEncoder struct {
	// Mock is only used to fill out unimplemented methods
	*mocks.ObjectEncoder
	t      *testing.T
	result Object
}

func (enc *ObjectEncoder) Result() Object {
	return enc.result
}

func NewObjectEncoder(t *testing.T) *ObjectEncoder {
	return &ObjectEncoder{
		ObjectEncoder: mocks.NewObjectEncoder(t),
		t:             t,
	}
}

func (enc *ObjectEncoder) addValue(k string, v Value) {
	if enc.result == nil {
		enc.result = make(Object)
	}

	enc.result[k] = v
}

func (enc *ObjectEncoder) AddObject(key string, ms zapcore.ObjectMarshaler) error {
	nenc := NewObjectEncoder(enc.t)
	if err := ms.MarshalLogObject(nenc); err != nil {
		return err
	}

	enc.addValue(key, nenc.Result())

	return nil
}

func (enc *ObjectEncoder) AddArray(key string, ms zapcore.ArrayMarshaler) error {
	nenc := NewArrayEncoder(enc.t)
	if err := ms.MarshalLogArray(nenc); err != nil {
		return err
	}

	enc.addValue(key, nenc.Result())

	return nil
}

func (enc *ObjectEncoder) AddBool(key string, b bool) {
	enc.addValue(key, Bool(b))
}

func (enc *ObjectEncoder) AddString(key string, s string) {
	enc.addValue(key, String(s))
}

func (enc *ObjectEncoder) AddFloat32(key string, f float32) {
	enc.addValue(key, Float32(f))
}

func (enc *ObjectEncoder) AddFloat64(key string, f float64) {
	enc.addValue(key, Float64(f))
}

func (enc *ObjectEncoder) AddInt32(key string, i int32) {
	enc.addValue(key, Int32(i))
}

func (enc *ObjectEncoder) AddUint32(key string, i uint32) {
	enc.addValue(key, Uint32(i))
}

func (enc *ObjectEncoder) AddInt64(key string, i int64) {
	enc.addValue(key, Int64(i))
}

func (enc *ObjectEncoder) AddUint64(key string, i uint64) {
	enc.addValue(key, Uint64(i))
}

type ArrayEncoder struct {
	// Mock is only used to fill out unimplemented methods
	*mocks.ArrayEncoder
	t      *testing.T
	result Array
}

func NewArrayEncoder(t *testing.T) *ArrayEncoder {
	return &ArrayEncoder{
		ArrayEncoder: mocks.NewArrayEncoder(t),
		t:            t,
	}
}

func (enc *ArrayEncoder) appendValue(v Value) {
	enc.result = append(enc.result, v)
}

func (enc *ArrayEncoder) Result() Array {
	return enc.result
}

func (enc *ArrayEncoder) AppendObject(ms zapcore.ObjectMarshaler) error {
	nenc := NewObjectEncoder(enc.t)
	if err := ms.MarshalLogObject(nenc); err != nil {
		return err
	}

	enc.appendValue(nenc.Result())

	return nil
}

func (enc *ArrayEncoder) AppendArray(ms zapcore.ArrayMarshaler) error {
	nenc := NewArrayEncoder(enc.t)
	if err := ms.MarshalLogArray(nenc); err != nil {
		return err
	}

	enc.appendValue(nenc.Result())

	return nil
}

func (enc *ArrayEncoder) AppendBool(b bool) {
	enc.appendValue(Bool(b))
}

func (enc *ArrayEncoder) AppendString(s string) {
	enc.appendValue(String(s))
}

func (enc *ArrayEncoder) AppendFloat32(f float32) {
	enc.appendValue(Float32(f))
}

func (enc *ArrayEncoder) AppendFloat64(f float64) {
	enc.appendValue(Float64(f))
}

func (enc *ArrayEncoder) AppendInt32(i int32) {
	enc.appendValue(Int32(i))
}

func (enc *ArrayEncoder) AppendUint32(i uint32) {
	enc.appendValue(Uint32(i))
}

func (enc *ArrayEncoder) AppendInt64(i int64) {
	enc.appendValue(Int64(i))
}

func (enc *ArrayEncoder) AppendUint64(i uint64) {
	enc.appendValue(Uint64(i))
}

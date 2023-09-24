package protolog

import (
	"encoding/base64"
	"fmt"

	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	anypbAnyFullName = (&anypb.Any{}).ProtoReflect().Descriptor().FullName()
)

type encoder interface {
	AppendObject(zapcore.ObjectMarshaler) error
	AppendArray(zapcore.ArrayMarshaler) error
	AppendBool(bool)
	AppendString(string)
	AppendFloat32(float32)
	AppendFloat64(float64)
	AppendInt32(int32)
	AppendUint32(uint32)
	AppendInt64(int64)
	AppendUint64(uint64)
}

type fieldEncoder struct {
	enc zapcore.ObjectEncoder
	key string
}

func (e fieldEncoder) AppendObject(ms zapcore.ObjectMarshaler) error {
	return e.enc.AddObject(e.key, ms)
}

func (e fieldEncoder) AppendArray(ms zapcore.ArrayMarshaler) error {
	return e.enc.AddArray(e.key, ms)
}

func (e fieldEncoder) AppendBool(b bool) {
	e.enc.AddBool(e.key, b)
}

func (e fieldEncoder) AppendString(s string) {
	e.enc.AddString(e.key, s)
}

func (e fieldEncoder) AppendFloat32(f float32) {
	e.enc.AddFloat32(e.key, f)
}

func (e fieldEncoder) AppendFloat64(f float64) {
	e.enc.AddFloat64(e.key, f)
}

func (e fieldEncoder) AppendInt32(i int32) {
	e.enc.AddInt32(e.key, i)
}

func (e fieldEncoder) AppendUint32(u uint32) {
	e.enc.AddUint32(e.key, u)
}

func (e fieldEncoder) AppendInt64(i int64) {
	e.enc.AddInt64(e.key, i)
}

func (e fieldEncoder) AppendUint64(u uint64) {
	e.enc.AddUint64(e.key, u)
}

func appendField[E encoder](enc E, fd protoreflect.FieldDescriptor, v protoreflect.Value) error {
	switch {
	case fd.IsMap():
		return enc.AppendObject(mapMarshaler{
			ValueDesc: fd.MapValue(),
			Map:       v.Map(),
		})

	case fd.IsList():
		return enc.AppendArray(listMarshaler{
			Desc: fd,
			List: v.List(),
		})
	}

	return appendValue(enc, fd, v)
}

func appendValue[E encoder](enc E, fd protoreflect.FieldDescriptor, v protoreflect.Value) error {
	switch fd.Kind() {
	case protoreflect.MessageKind:
		return enc.AppendObject(messageMarshaler{
			Message: v.Message(),
		})

	case protoreflect.BoolKind:
		enc.AppendBool(v.Bool())

	case protoreflect.BytesKind:
		enc.AppendString(base64.StdEncoding.EncodeToString(v.Bytes()))

	case protoreflect.EnumKind:
		enumName := fd.Enum().Values().ByNumber(v.Enum()).Name()
		enc.AppendString(string(enumName))

	case protoreflect.FloatKind:
		enc.AppendFloat32(float32(v.Float()))

	case protoreflect.DoubleKind:
		enc.AppendFloat64(v.Float())

	case protoreflect.StringKind:
		enc.AppendString(v.String())

	case protoreflect.Int32Kind, protoreflect.Sfixed32Kind, protoreflect.Sint32Kind:
		enc.AppendInt32(int32(v.Int()))

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		enc.AppendUint32(uint32(v.Uint()))

	case protoreflect.Int64Kind, protoreflect.Sfixed64Kind, protoreflect.Sint64Kind:
		enc.AppendInt64(v.Int())

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		enc.AppendUint64(v.Uint())

	default:
		return fmt.Errorf("cannot marshal value for protobuf field %s", fd.FullName())
	}

	return nil
}

type listMarshaler struct {
	Desc protoreflect.FieldDescriptor
	List protoreflect.List
}

func (m listMarshaler) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for i := 0; i < m.List.Len(); i++ {
		if err := appendValue(enc, m.Desc, m.List.Get(i)); err != nil {
			return err
		}
	}

	return nil
}

type mapMarshaler struct {
	ValueDesc protoreflect.FieldDescriptor
	Map       protoreflect.Map
}

func (m mapMarshaler) MarshalLogObject(enc zapcore.ObjectEncoder) (err error) {
	m.Map.Range(func(mk protoreflect.MapKey, v protoreflect.Value) bool {
		err = appendValue(fieldEncoder{enc, mk.String()}, m.ValueDesc, v)

		return err == nil
	})

	return err
}

type messageMarshaler struct {
	Typed   bool
	Message protoreflect.Message
}

func (m *messageMarshaler) marshalAny(enc zapcore.ObjectEncoder) error {
	anyMsg, ok := m.Message.Interface().(interface{ UnmarshalNew() (proto.Message, error) })
	if !ok {
		return fmt.Errorf("message %s does not implement UnmarshalNew method", anypbAnyFullName)
	}

	v, err := anyMsg.UnmarshalNew()
	if err != nil {
		return err
	}

	return messageMarshaler{
		Typed:   true,
		Message: v.ProtoReflect(),
	}.MarshalLogObject(enc)
}

func (m messageMarshaler) MarshalLogObject(enc zapcore.ObjectEncoder) (err error) {
	fullName := m.Message.Descriptor().FullName()

	if fullName == anypbAnyFullName {
		return m.marshalAny(enc)
	}

	if m.Typed {
		enc.AddString("@type", "type.googleapis.com/"+string(fullName))
	}

	m.Message.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		key := fd.JSONName()

		err = appendField(fieldEncoder{enc, key}, fd, v)

		return err == nil
	})

	return err
}

type objectMarshalerFunc func(enc zapcore.ObjectEncoder) error

func (fn objectMarshalerFunc) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return fn(enc)
}

type Options struct {
	Typed bool
}

func (opts Options) MarshalerOf(msg proto.Message) zapcore.ObjectMarshaler {
	if msg == nil {
		return objectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			return nil
		})
	}

	return messageMarshaler{
		Typed:   opts.Typed,
		Message: msg.ProtoReflect(),
	}
}

func MarshalerOf(msg proto.Message) zapcore.ObjectMarshaler {
	return Options{}.MarshalerOf(msg)
}

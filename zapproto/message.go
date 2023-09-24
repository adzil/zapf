package zapproto

import (
	"github.com/adzil/zapf/internal/protolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
)

// Message constructs a field with a given key and Protobuf message. It will
// serialize the Protobuf message lazily.
func Message(key string, msg proto.Message) zap.Field {
	return zap.Object(key, protolog.MarshalerOf(msg))
}

// TypedMessage constructs a field with a given key and Protobuf message. It
// will serialize the Protobuf message with its type URL lazily.
func TypedMessage(key string, msg proto.Message) zap.Field {
	opts := protolog.Options{
		Typed: true,
	}

	return zap.Object(key, opts.MarshalerOf(msg))
}

type protosMarshaler[M proto.Message] struct {
	Marshaler protolog.Options
	Messages  []M
}

func (m protosMarshaler[M]) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, msg := range m.Messages {
		if err := enc.AppendObject(m.Marshaler.MarshalerOf(msg)); err != nil {
			return err
		}
	}

	return nil
}

// Messages constructs a field with a given key and typed Protobuf messages. It
// will serialize the Protobuf messages lazily.
func Messages[M proto.Message](key string, msgs []M) zap.Field {
	return zap.Array(key, protosMarshaler[M]{
		Messages: msgs,
	})
}

// TypedMessages constructs a field with a given key and typed Protobuf
// messages. It will serialize the Protobuf messages with their type URL lazily.
func TypedMessages[M proto.Message](key string, msgs []M) zap.Field {
	return zap.Array(key, protosMarshaler[M]{
		Marshaler: protolog.Options{
			Typed: true,
		},
		Messages: msgs,
	})
}

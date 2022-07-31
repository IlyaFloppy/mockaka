package mockaka

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// matchesUnaryHandlerSignature returns true if provided method matches grpc unary call signature.
// It checks that arguments are (context.Context, *ProtobufMessage)
// And return values are (*ProtobufMessage, error).
func matchUnaryHandlerSignature(method reflect.Method) error {
	args := make([]reflect.Type, 0, method.Type.NumIn())
	rets := make([]reflect.Type, 0, method.Type.NumOut())
	for i := 0; i < method.Type.NumIn(); i++ {
		args = append(args, method.Type.In(i))
	}
	for i := 0; i < method.Type.NumOut(); i++ {
		rets = append(rets, method.Type.Out(i))
	}

	if len(args) != 2 || len(rets) != 2 {
		return errors.New("method's signature should be (context.Context, *ProtobufMessage) (*ProtobufMessage, error)")
	}

	// Check that second argument is context.Context.
	if _, ok := reflect.New(args[0]).Interface().(*context.Context); !ok {
		return errors.New("first argument is not context.Context")
	}

	// Check that third argument is *ProtobufMessage.
	if err := matchProtobufMessage(args[1]); err != nil {
		return errors.Wrap(err, "second argument is not a protobuf message")
	}

	// Check that first return value is *ProtobufMessage.
	if err := matchProtobufMessage(rets[0]); err != nil {
		return errors.Wrap(err, "first return value is not a protobuf message")
	}

	// Check that second return value is an error.
	if _, ok := reflect.New(rets[1]).Interface().(*error); !ok {
		return errors.New("second argument is not error")
	}

	return nil
}

// matchProtobufMessage checks that typ is a protobuf message.
// All protobuf messages are pointers and have these fields:
// - state         protoimpl.MessageState
// - sizeCache     protoimpl.SizeCache
// - unknownFields protoimpl.UnknownFields.
func matchProtobufMessage(typ reflect.Type) error {
	if typ.Kind() != reflect.Pointer {
		return errors.New(typ.String() + " is not a pointer")
	}

	checks := map[string]reflect.Type{
		"state":         reflect.TypeOf(protoimpl.MessageState{}),
		"sizeCache":     reflect.TypeOf(protoimpl.SizeCache(0)),
		"unknownFields": reflect.TypeOf(protoimpl.UnknownFields{}),
	}

	for name, t := range checks {
		if f, ok := typ.Elem().FieldByName(name); !ok || f.Type != t {
			return fmt.Errorf("%s does not have %q field of type %q", typ.String(), name, t.String())
		}
	}

	return nil
}

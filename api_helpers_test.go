package mockaka_test

import (
	"testing"

	"github.com/IlyaFloppy/mockaka"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type ProtobufMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (*ProtobufMessage) ProtoMessage()                        {}
func (x *ProtobufMessage) ProtoReflect() protoreflect.Message { return nil }

func TestPB(t *testing.T) {
	msg := &ProtobufMessage{
		sizeCache:     42,
		unknownFields: []byte{1, 2, 3, 4},
		Message:       "Persists",
	}

	converted := mockaka.PB(msg)

	require.Equal(t, "Persists", converted.Message)
	require.Zero(t, converted.sizeCache)
	require.Zero(t, converted.unknownFields)
}

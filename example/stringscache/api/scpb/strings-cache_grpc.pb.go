// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: strings-cache.proto

package scpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StringsCacheServiceClient is the client API for StringsCacheService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StringsCacheServiceClient interface {
	Reverse(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
	Invalidate(ctx context.Context, in *InvalidateRequest, opts ...grpc.CallOption) (*InvalidateResponse, error)
}

type stringsCacheServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStringsCacheServiceClient(cc grpc.ClientConnInterface) StringsCacheServiceClient {
	return &stringsCacheServiceClient{cc}
}

func (c *stringsCacheServiceClient) Reverse(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/scpb.StringsCacheService/Reverse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stringsCacheServiceClient) Invalidate(ctx context.Context, in *InvalidateRequest, opts ...grpc.CallOption) (*InvalidateResponse, error) {
	out := new(InvalidateResponse)
	err := c.cc.Invoke(ctx, "/scpb.StringsCacheService/Invalidate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StringsCacheServiceServer is the server API for StringsCacheService service.
// All implementations must embed UnimplementedStringsCacheServiceServer
// for forward compatibility
type StringsCacheServiceServer interface {
	Reverse(context.Context, *Message) (*Message, error)
	Invalidate(context.Context, *InvalidateRequest) (*InvalidateResponse, error)
	mustEmbedUnimplementedStringsCacheServiceServer()
}

// UnimplementedStringsCacheServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStringsCacheServiceServer struct {
}

func (UnimplementedStringsCacheServiceServer) Reverse(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reverse not implemented")
}
func (UnimplementedStringsCacheServiceServer) Invalidate(context.Context, *InvalidateRequest) (*InvalidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Invalidate not implemented")
}
func (UnimplementedStringsCacheServiceServer) mustEmbedUnimplementedStringsCacheServiceServer() {}

// UnsafeStringsCacheServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StringsCacheServiceServer will
// result in compilation errors.
type UnsafeStringsCacheServiceServer interface {
	mustEmbedUnimplementedStringsCacheServiceServer()
}

func RegisterStringsCacheServiceServer(s grpc.ServiceRegistrar, srv StringsCacheServiceServer) {
	s.RegisterService(&StringsCacheService_ServiceDesc, srv)
}

func _StringsCacheService_Reverse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StringsCacheServiceServer).Reverse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scpb.StringsCacheService/Reverse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StringsCacheServiceServer).Reverse(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _StringsCacheService_Invalidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InvalidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StringsCacheServiceServer).Invalidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/scpb.StringsCacheService/Invalidate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StringsCacheServiceServer).Invalidate(ctx, req.(*InvalidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StringsCacheService_ServiceDesc is the grpc.ServiceDesc for StringsCacheService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StringsCacheService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "scpb.StringsCacheService",
	HandlerType: (*StringsCacheServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Reverse",
			Handler:    _StringsCacheService_Reverse_Handler,
		},
		{
			MethodName: "Invalidate",
			Handler:    _StringsCacheService_Invalidate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "strings-cache.proto",
}

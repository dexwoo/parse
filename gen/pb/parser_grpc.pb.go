// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

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

// ParserClient is the client API for Parser service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ParserClient interface {
	GetCurrentBlock(ctx context.Context, in *GetCurrentBlockRequest, opts ...grpc.CallOption) (*GetCurrentBlockResponse, error)
}

type parserClient struct {
	cc grpc.ClientConnInterface
}

func NewParserClient(cc grpc.ClientConnInterface) ParserClient {
	return &parserClient{cc}
}

func (c *parserClient) GetCurrentBlock(ctx context.Context, in *GetCurrentBlockRequest, opts ...grpc.CallOption) (*GetCurrentBlockResponse, error) {
	out := new(GetCurrentBlockResponse)
	err := c.cc.Invoke(ctx, "/parser.parser/GetCurrentBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ParserServer is the server API for Parser service.
// All implementations should embed UnimplementedParserServer
// for forward compatibility
type ParserServer interface {
	GetCurrentBlock(context.Context, *GetCurrentBlockRequest) (*GetCurrentBlockResponse, error)
}

// UnimplementedParserServer should be embedded to have forward compatible implementations.
type UnimplementedParserServer struct {
}

func (UnimplementedParserServer) GetCurrentBlock(context.Context, *GetCurrentBlockRequest) (*GetCurrentBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrentBlock not implemented")
}

// UnsafeParserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ParserServer will
// result in compilation errors.
type UnsafeParserServer interface {
	mustEmbedUnimplementedParserServer()
}

func RegisterParserServer(s grpc.ServiceRegistrar, srv ParserServer) {
	s.RegisterService(&Parser_ServiceDesc, srv)
}

func _Parser_GetCurrentBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCurrentBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ParserServer).GetCurrentBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/parser.parser/GetCurrentBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ParserServer).GetCurrentBlock(ctx, req.(*GetCurrentBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Parser_ServiceDesc is the grpc.ServiceDesc for Parser service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Parser_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "parser.parser",
	HandlerType: (*ParserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCurrentBlock",
			Handler:    _Parser_GetCurrentBlock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "parser.proto",
}

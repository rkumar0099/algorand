// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package client

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

// RPCClientServiceClient is the client API for RPCClientService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RPCClientServiceClient interface {
	SendTx(ctx context.Context, in *ReqTx, opts ...grpc.CallOption) (*ResEmpty, error)
	SendResTx(ctx context.Context, in *ResTx, opts ...grpc.CallOption) (*ResEmpty, error)
}

type rPCClientServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRPCClientServiceClient(cc grpc.ClientConnInterface) RPCClientServiceClient {
	return &rPCClientServiceClient{cc}
}

func (c *rPCClientServiceClient) SendTx(ctx context.Context, in *ReqTx, opts ...grpc.CallOption) (*ResEmpty, error) {
	out := new(ResEmpty)
	err := c.cc.Invoke(ctx, "/RPCClientService/SendTx", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCClientServiceClient) SendResTx(ctx context.Context, in *ResTx, opts ...grpc.CallOption) (*ResEmpty, error) {
	out := new(ResEmpty)
	err := c.cc.Invoke(ctx, "/RPCClientService/SendResTx", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RPCClientServiceServer is the server API for RPCClientService service.
// All implementations must embed UnimplementedRPCClientServiceServer
// for forward compatibility
type RPCClientServiceServer interface {
	SendTx(context.Context, *ReqTx) (*ResEmpty, error)
	SendResTx(context.Context, *ResTx) (*ResEmpty, error)
	mustEmbedUnimplementedRPCClientServiceServer()
}

// UnimplementedRPCClientServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRPCClientServiceServer struct {
}

func (UnimplementedRPCClientServiceServer) SendTx(context.Context, *ReqTx) (*ResEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTx not implemented")
}
func (UnimplementedRPCClientServiceServer) SendResTx(context.Context, *ResTx) (*ResEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendResTx not implemented")
}
func (UnimplementedRPCClientServiceServer) mustEmbedUnimplementedRPCClientServiceServer() {}

// UnsafeRPCClientServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RPCClientServiceServer will
// result in compilation errors.
type UnsafeRPCClientServiceServer interface {
	mustEmbedUnimplementedRPCClientServiceServer()
}

func RegisterRPCClientServiceServer(s grpc.ServiceRegistrar, srv RPCClientServiceServer) {
	s.RegisterService(&RPCClientService_ServiceDesc, srv)
}

func _RPCClientService_SendTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqTx)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCClientServiceServer).SendTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RPCClientService/SendTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCClientServiceServer).SendTx(ctx, req.(*ReqTx))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCClientService_SendResTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResTx)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCClientServiceServer).SendResTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RPCClientService/SendResTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCClientServiceServer).SendResTx(ctx, req.(*ResTx))
	}
	return interceptor(ctx, in, info, handler)
}

// RPCClientService_ServiceDesc is the grpc.ServiceDesc for RPCClientService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RPCClientService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RPCClientService",
	HandlerType: (*RPCClientServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendTx",
			Handler:    _RPCClientService_SendTx_Handler,
		},
		{
			MethodName: "SendResTx",
			Handler:    _RPCClientService_SendResTx_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "client.proto",
}

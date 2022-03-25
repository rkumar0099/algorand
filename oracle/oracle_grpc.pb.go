// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package oracle

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

// RPCOracleServiceClient is the client API for RPCOracleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RPCOracleServiceClient interface {
	SendOPP(ctx context.Context, in *ReqOPP, opts ...grpc.CallOption) (*ResOPP, error)
}

type rPCOracleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRPCOracleServiceClient(cc grpc.ClientConnInterface) RPCOracleServiceClient {
	return &rPCOracleServiceClient{cc}
}

func (c *rPCOracleServiceClient) SendOPP(ctx context.Context, in *ReqOPP, opts ...grpc.CallOption) (*ResOPP, error) {
	out := new(ResOPP)
	err := c.cc.Invoke(ctx, "/RPCOracleService/SendOPP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RPCOracleServiceServer is the server API for RPCOracleService service.
// All implementations must embed UnimplementedRPCOracleServiceServer
// for forward compatibility
type RPCOracleServiceServer interface {
	SendOPP(context.Context, *ReqOPP) (*ResOPP, error)
	mustEmbedUnimplementedRPCOracleServiceServer()
}

// UnimplementedRPCOracleServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRPCOracleServiceServer struct {
}

func (UnimplementedRPCOracleServiceServer) SendOPP(context.Context, *ReqOPP) (*ResOPP, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendOPP not implemented")
}
func (UnimplementedRPCOracleServiceServer) mustEmbedUnimplementedRPCOracleServiceServer() {}

// UnsafeRPCOracleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RPCOracleServiceServer will
// result in compilation errors.
type UnsafeRPCOracleServiceServer interface {
	mustEmbedUnimplementedRPCOracleServiceServer()
}

func RegisterRPCOracleServiceServer(s grpc.ServiceRegistrar, srv RPCOracleServiceServer) {
	s.RegisterService(&RPCOracleService_ServiceDesc, srv)
}

func _RPCOracleService_SendOPP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqOPP)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCOracleServiceServer).SendOPP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RPCOracleService/SendOPP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCOracleServiceServer).SendOPP(ctx, req.(*ReqOPP))
	}
	return interceptor(ctx, in, info, handler)
}

// RPCOracleService_ServiceDesc is the grpc.ServiceDesc for RPCOracleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RPCOracleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RPCOracleService",
	HandlerType: (*RPCOracleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendOPP",
			Handler:    _RPCOracleService_SendOPP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "oracle.proto",
}

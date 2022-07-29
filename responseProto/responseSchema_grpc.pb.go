// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: responseProto/responseSchema.proto

package responseSchemapb

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

// CodeforcesClient is the client API for Codeforces service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CodeforcesClient interface {
	CfRequest(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type codeforcesClient struct {
	cc grpc.ClientConnInterface
}

func NewCodeforcesClient(cc grpc.ClientConnInterface) CodeforcesClient {
	return &codeforcesClient{cc}
}

func (c *codeforcesClient) CfRequest(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/responseSchemapb.Codeforces/CfRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CodeforcesServer is the server API for Codeforces service.
// All implementations must embed UnimplementedCodeforcesServer
// for forward compatibility
type CodeforcesServer interface {
	CfRequest(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedCodeforcesServer()
}

// UnimplementedCodeforcesServer must be embedded to have forward compatible implementations.
type UnimplementedCodeforcesServer struct {
}

func (UnimplementedCodeforcesServer) CfRequest(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CfRequest not implemented")
}
func (UnimplementedCodeforcesServer) mustEmbedUnimplementedCodeforcesServer() {}

// UnsafeCodeforcesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CodeforcesServer will
// result in compilation errors.
type UnsafeCodeforcesServer interface {
	mustEmbedUnimplementedCodeforcesServer()
}

func RegisterCodeforcesServer(s grpc.ServiceRegistrar, srv CodeforcesServer) {
	s.RegisterService(&Codeforces_ServiceDesc, srv)
}

func _Codeforces_CfRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeforcesServer).CfRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/responseSchemapb.Codeforces/CfRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeforcesServer).CfRequest(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Codeforces_ServiceDesc is the grpc.ServiceDesc for Codeforces service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Codeforces_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "responseSchemapb.Codeforces",
	HandlerType: (*CodeforcesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CfRequest",
			Handler:    _Codeforces_CfRequest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "responseProto/responseSchema.proto",
}
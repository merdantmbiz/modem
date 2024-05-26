// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: internal/token.proto

package auth

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

const (
	AuthService_AuthStream_FullMethodName = "/auth.AuthService/AuthStream"
)

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	AuthStream(ctx context.Context, opts ...grpc.CallOption) (AuthService_AuthStreamClient, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) AuthStream(ctx context.Context, opts ...grpc.CallOption) (AuthService_AuthStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &AuthService_ServiceDesc.Streams[0], AuthService_AuthStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &authServiceAuthStreamClient{stream}
	return x, nil
}

type AuthService_AuthStreamClient interface {
	Send(*AuthRequest) error
	Recv() (*AuthResponse, error)
	grpc.ClientStream
}

type authServiceAuthStreamClient struct {
	grpc.ClientStream
}

func (x *authServiceAuthStreamClient) Send(m *AuthRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *authServiceAuthStreamClient) Recv() (*AuthResponse, error) {
	m := new(AuthResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	AuthStream(AuthService_AuthStreamServer) error
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) AuthStream(AuthService_AuthStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method AuthStream not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_AuthStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AuthServiceServer).AuthStream(&authServiceAuthStreamServer{stream})
}

type AuthService_AuthStreamServer interface {
	Send(*AuthResponse) error
	Recv() (*AuthRequest, error)
	grpc.ServerStream
}

type authServiceAuthStreamServer struct {
	grpc.ServerStream
}

func (x *authServiceAuthStreamServer) Send(m *AuthResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *authServiceAuthStreamServer) Recv() (*AuthRequest, error) {
	m := new(AuthRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AuthStream",
			Handler:       _AuthService_AuthStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "internal/token.proto",
}

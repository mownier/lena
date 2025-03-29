// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpcendpoint

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LenaServiceClient is the client API for LenaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LenaServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	SignIn(ctx context.Context, in *SignInRequest, opts ...grpc.CallOption) (*SignInResponse, error)
	SignOut(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Verify(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error)
}

type lenaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLenaServiceClient(cc grpc.ClientConnInterface) LenaServiceClient {
	return &lenaServiceClient{cc}
}

func (c *lenaServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/grpcendpoints.LenaService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lenaServiceClient) SignIn(ctx context.Context, in *SignInRequest, opts ...grpc.CallOption) (*SignInResponse, error) {
	out := new(SignInResponse)
	err := c.cc.Invoke(ctx, "/grpcendpoints.LenaService/SignIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lenaServiceClient) SignOut(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/grpcendpoints.LenaService/SignOut", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lenaServiceClient) Verify(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/grpcendpoints.LenaService/Verify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lenaServiceClient) Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error) {
	out := new(RefreshResponse)
	err := c.cc.Invoke(ctx, "/grpcendpoints.LenaService/Refresh", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LenaServiceServer is the server API for LenaService service.
// All implementations must embed UnimplementedLenaServiceServer
// for forward compatibility
type LenaServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	SignIn(context.Context, *SignInRequest) (*SignInResponse, error)
	SignOut(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	Verify(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error)
	mustEmbedUnimplementedLenaServiceServer()
}

// UnimplementedLenaServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLenaServiceServer struct {
}

func (UnimplementedLenaServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedLenaServiceServer) SignIn(context.Context, *SignInRequest) (*SignInResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignIn not implemented")
}
func (UnimplementedLenaServiceServer) SignOut(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignOut not implemented")
}
func (UnimplementedLenaServiceServer) Verify(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Verify not implemented")
}
func (UnimplementedLenaServiceServer) Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
func (UnimplementedLenaServiceServer) mustEmbedUnimplementedLenaServiceServer() {}

// UnsafeLenaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LenaServiceServer will
// result in compilation errors.
type UnsafeLenaServiceServer interface {
	mustEmbedUnimplementedLenaServiceServer()
}

func RegisterLenaServiceServer(s grpc.ServiceRegistrar, srv LenaServiceServer) {
	s.RegisterService(&LenaService_ServiceDesc, srv)
}

func _LenaService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LenaServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpcendpoints.LenaService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LenaServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LenaService_SignIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignInRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LenaServiceServer).SignIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpcendpoints.LenaService/SignIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LenaServiceServer).SignIn(ctx, req.(*SignInRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LenaService_SignOut_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LenaServiceServer).SignOut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpcendpoints.LenaService/SignOut",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LenaServiceServer).SignOut(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LenaService_Verify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LenaServiceServer).Verify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpcendpoints.LenaService/Verify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LenaServiceServer).Verify(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LenaService_Refresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LenaServiceServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpcendpoints.LenaService/Refresh",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LenaServiceServer).Refresh(ctx, req.(*RefreshRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LenaService_ServiceDesc is the grpc.ServiceDesc for LenaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LenaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpcendpoints.LenaService",
	HandlerType: (*LenaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _LenaService_Register_Handler,
		},
		{
			MethodName: "SignIn",
			Handler:    _LenaService_SignIn_Handler,
		},
		{
			MethodName: "SignOut",
			Handler:    _LenaService_SignOut_Handler,
		},
		{
			MethodName: "Verify",
			Handler:    _LenaService_Verify_Handler,
		},
		{
			MethodName: "Refresh",
			Handler:    _LenaService_Refresh_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "endpoints/grpcendpoint/lena.proto",
}

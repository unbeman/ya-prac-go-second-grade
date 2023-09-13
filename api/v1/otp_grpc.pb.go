// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: otp.proto

package v1

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
	OtpService_OTPGenerate_FullMethodName = "/pass_keeper.OtpService/OTPGenerate"
	OtpService_OTPVerify_FullMethodName   = "/pass_keeper.OtpService/OTPVerify"
	OtpService_OTPValidate_FullMethodName = "/pass_keeper.OtpService/OTPValidate"
	OtpService_OTPDisable_FullMethodName  = "/pass_keeper.OtpService/OTPDisable"
)

// OtpServiceClient is the client API for OtpService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OtpServiceClient interface {
	OTPGenerate(ctx context.Context, in *OTPGenRequest, opts ...grpc.CallOption) (*OTPGenResponse, error)
	OTPVerify(ctx context.Context, in *OTPVerifyRequest, opts ...grpc.CallOption) (*OTPVerifyResponse, error)
	OTPValidate(ctx context.Context, in *OTPValidateRequest, opts ...grpc.CallOption) (*OTPValidateResponse, error)
	OTPDisable(ctx context.Context, in *OTPDisableRequest, opts ...grpc.CallOption) (*OTPDisableResponse, error)
}

type otpServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOtpServiceClient(cc grpc.ClientConnInterface) OtpServiceClient {
	return &otpServiceClient{cc}
}

func (c *otpServiceClient) OTPGenerate(ctx context.Context, in *OTPGenRequest, opts ...grpc.CallOption) (*OTPGenResponse, error) {
	out := new(OTPGenResponse)
	err := c.cc.Invoke(ctx, OtpService_OTPGenerate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *otpServiceClient) OTPVerify(ctx context.Context, in *OTPVerifyRequest, opts ...grpc.CallOption) (*OTPVerifyResponse, error) {
	out := new(OTPVerifyResponse)
	err := c.cc.Invoke(ctx, OtpService_OTPVerify_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *otpServiceClient) OTPValidate(ctx context.Context, in *OTPValidateRequest, opts ...grpc.CallOption) (*OTPValidateResponse, error) {
	out := new(OTPValidateResponse)
	err := c.cc.Invoke(ctx, OtpService_OTPValidate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *otpServiceClient) OTPDisable(ctx context.Context, in *OTPDisableRequest, opts ...grpc.CallOption) (*OTPDisableResponse, error) {
	out := new(OTPDisableResponse)
	err := c.cc.Invoke(ctx, OtpService_OTPDisable_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OtpServiceServer is the server API for OtpService service.
// All implementations must embed UnimplementedOtpServiceServer
// for forward compatibility
type OtpServiceServer interface {
	OTPGenerate(context.Context, *OTPGenRequest) (*OTPGenResponse, error)
	OTPVerify(context.Context, *OTPVerifyRequest) (*OTPVerifyResponse, error)
	OTPValidate(context.Context, *OTPValidateRequest) (*OTPValidateResponse, error)
	OTPDisable(context.Context, *OTPDisableRequest) (*OTPDisableResponse, error)
	mustEmbedUnimplementedOtpServiceServer()
}

// UnimplementedOtpServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOtpServiceServer struct {
}

func (UnimplementedOtpServiceServer) OTPGenerate(context.Context, *OTPGenRequest) (*OTPGenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OTPGenerate not implemented")
}
func (UnimplementedOtpServiceServer) OTPVerify(context.Context, *OTPVerifyRequest) (*OTPVerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OTPVerify not implemented")
}
func (UnimplementedOtpServiceServer) OTPValidate(context.Context, *OTPValidateRequest) (*OTPValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OTPValidate not implemented")
}
func (UnimplementedOtpServiceServer) OTPDisable(context.Context, *OTPDisableRequest) (*OTPDisableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OTPDisable not implemented")
}
func (UnimplementedOtpServiceServer) mustEmbedUnimplementedOtpServiceServer() {}

// UnsafeOtpServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OtpServiceServer will
// result in compilation errors.
type UnsafeOtpServiceServer interface {
	mustEmbedUnimplementedOtpServiceServer()
}

func RegisterOtpServiceServer(s grpc.ServiceRegistrar, srv OtpServiceServer) {
	s.RegisterService(&OtpService_ServiceDesc, srv)
}

func _OtpService_OTPGenerate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OTPGenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OtpServiceServer).OTPGenerate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OtpService_OTPGenerate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OtpServiceServer).OTPGenerate(ctx, req.(*OTPGenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OtpService_OTPVerify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OTPVerifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OtpServiceServer).OTPVerify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OtpService_OTPVerify_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OtpServiceServer).OTPVerify(ctx, req.(*OTPVerifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OtpService_OTPValidate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OTPValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OtpServiceServer).OTPValidate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OtpService_OTPValidate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OtpServiceServer).OTPValidate(ctx, req.(*OTPValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OtpService_OTPDisable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OTPDisableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OtpServiceServer).OTPDisable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OtpService_OTPDisable_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OtpServiceServer).OTPDisable(ctx, req.(*OTPDisableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OtpService_ServiceDesc is the grpc.ServiceDesc for OtpService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OtpService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pass_keeper.OtpService",
	HandlerType: (*OtpServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OTPGenerate",
			Handler:    _OtpService_OTPGenerate_Handler,
		},
		{
			MethodName: "OTPVerify",
			Handler:    _OtpService_OTPVerify_Handler,
		},
		{
			MethodName: "OTPValidate",
			Handler:    _OtpService_OTPValidate_Handler,
		},
		{
			MethodName: "OTPDisable",
			Handler:    _OtpService_OTPDisable_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "otp.proto",
}
// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: pkg/proto/monitaization.proto

package monit

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
	MonitizationService_HealthCheck_FullMethodName = "/monitization.MonitizationService/HealthCheck"
	MonitizationService_VideoReward_FullMethodName = "/monitization.MonitizationService/VideoReward"
)

// MonitizationServiceClient is the client API for MonitizationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MonitizationServiceClient interface {
	HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
	VideoReward(ctx context.Context, in *VideoRewardRequest, opts ...grpc.CallOption) (*VideoRewardResponse, error)
}

type monitizationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMonitizationServiceClient(cc grpc.ClientConnInterface) MonitizationServiceClient {
	return &monitizationServiceClient{cc}
}

func (c *monitizationServiceClient) HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, MonitizationService_HealthCheck_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *monitizationServiceClient) VideoReward(ctx context.Context, in *VideoRewardRequest, opts ...grpc.CallOption) (*VideoRewardResponse, error) {
	out := new(VideoRewardResponse)
	err := c.cc.Invoke(ctx, MonitizationService_VideoReward_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MonitizationServiceServer is the server API for MonitizationService service.
// All implementations must embed UnimplementedMonitizationServiceServer
// for forward compatibility
type MonitizationServiceServer interface {
	HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
	VideoReward(context.Context, *VideoRewardRequest) (*VideoRewardResponse, error)
	mustEmbedUnimplementedMonitizationServiceServer()
}

// UnimplementedMonitizationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMonitizationServiceServer struct {
}

func (UnimplementedMonitizationServiceServer) HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedMonitizationServiceServer) VideoReward(context.Context, *VideoRewardRequest) (*VideoRewardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VideoReward not implemented")
}
func (UnimplementedMonitizationServiceServer) mustEmbedUnimplementedMonitizationServiceServer() {}

// UnsafeMonitizationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MonitizationServiceServer will
// result in compilation errors.
type UnsafeMonitizationServiceServer interface {
	mustEmbedUnimplementedMonitizationServiceServer()
}

func RegisterMonitizationServiceServer(s grpc.ServiceRegistrar, srv MonitizationServiceServer) {
	s.RegisterService(&MonitizationService_ServiceDesc, srv)
}

func _MonitizationService_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitizationServiceServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MonitizationService_HealthCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitizationServiceServer).HealthCheck(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MonitizationService_VideoReward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VideoRewardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitizationServiceServer).VideoReward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MonitizationService_VideoReward_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitizationServiceServer).VideoReward(ctx, req.(*VideoRewardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MonitizationService_ServiceDesc is the grpc.ServiceDesc for MonitizationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MonitizationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "monitization.MonitizationService",
	HandlerType: (*MonitizationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _MonitizationService_HealthCheck_Handler,
		},
		{
			MethodName: "VideoReward",
			Handler:    _MonitizationService_VideoReward_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/proto/monitaization.proto",
}

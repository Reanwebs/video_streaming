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

// VideoServiceClient is the client API for VideoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VideoServiceClient interface {
	HealthCheck(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	UploadVideo(ctx context.Context, opts ...grpc.CallOption) (VideoService_UploadVideoClient, error)
	FindArchivedVideoByUserId(ctx context.Context, in *FindArchivedVideoByUserIdRequest, opts ...grpc.CallOption) (*FindArchivedVideoByUserIdResponse, error)
	FindUserVideo(ctx context.Context, in *FindUserVideoRequest, opts ...grpc.CallOption) (*FindUserVideoResponse, error)
	ArchiveVideo(ctx context.Context, in *ArchiveVideoRequest, opts ...grpc.CallOption) (*ArchiveVideoResponse, error)
	FetchAllVideo(ctx context.Context, in *FetchAllVideoRequest, opts ...grpc.CallOption) (*FetchAllVideoResponse, error)
	GetVideoById(ctx context.Context, in *GetVideoByIdRequest, opts ...grpc.CallOption) (*GetVideoByIdResponse, error)
	ToggleStar(ctx context.Context, in *ToggleStarRequest, opts ...grpc.CallOption) (*ToggleStarResponse, error)
}

type videoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVideoServiceClient(cc grpc.ClientConnInterface) VideoServiceClient {
	return &videoServiceClient{cc}
}

func (c *videoServiceClient) HealthCheck(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/pb.VideoService/HealthCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) UploadVideo(ctx context.Context, opts ...grpc.CallOption) (VideoService_UploadVideoClient, error) {
	stream, err := c.cc.NewStream(ctx, &VideoService_ServiceDesc.Streams[0], "/pb.VideoService/UploadVideo", opts...)
	if err != nil {
		return nil, err
	}
	x := &videoServiceUploadVideoClient{stream}
	return x, nil
}

type VideoService_UploadVideoClient interface {
	Send(*UploadVideoRequest) error
	CloseAndRecv() (*UploadVideoResponse, error)
	grpc.ClientStream
}

type videoServiceUploadVideoClient struct {
	grpc.ClientStream
}

func (x *videoServiceUploadVideoClient) Send(m *UploadVideoRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *videoServiceUploadVideoClient) CloseAndRecv() (*UploadVideoResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadVideoResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *videoServiceClient) FindArchivedVideoByUserId(ctx context.Context, in *FindArchivedVideoByUserIdRequest, opts ...grpc.CallOption) (*FindArchivedVideoByUserIdResponse, error) {
	out := new(FindArchivedVideoByUserIdResponse)
	err := c.cc.Invoke(ctx, "/pb.VideoService/FindArchivedVideoByUserId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) FindUserVideo(ctx context.Context, in *FindUserVideoRequest, opts ...grpc.CallOption) (*FindUserVideoResponse, error) {
	out := new(FindUserVideoResponse)
	err := c.cc.Invoke(ctx, "/pb.VideoService/FindUserVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) ArchiveVideo(ctx context.Context, in *ArchiveVideoRequest, opts ...grpc.CallOption) (*ArchiveVideoResponse, error) {
	out := new(ArchiveVideoResponse)
	err := c.cc.Invoke(ctx, "/pb.VideoService/ArchiveVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) FetchAllVideo(ctx context.Context, in *FetchAllVideoRequest, opts ...grpc.CallOption) (*FetchAllVideoResponse, error) {
	out := new(FetchAllVideoResponse)
	err := c.cc.Invoke(ctx, "/pb.VideoService/FetchAllVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) GetVideoById(ctx context.Context, in *GetVideoByIdRequest, opts ...grpc.CallOption) (*GetVideoByIdResponse, error) {
	out := new(GetVideoByIdResponse)
	err := c.cc.Invoke(ctx, "/pb.VideoService/GetVideoById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) ToggleStar(ctx context.Context, in *ToggleStarRequest, opts ...grpc.CallOption) (*ToggleStarResponse, error) {
	out := new(ToggleStarResponse)
	err := c.cc.Invoke(ctx, "/pb.VideoService/ToggleStar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VideoServiceServer is the server API for VideoService service.
// All implementations must embed UnimplementedVideoServiceServer
// for forward compatibility
type VideoServiceServer interface {
	HealthCheck(context.Context, *Request) (*Response, error)
	UploadVideo(VideoService_UploadVideoServer) error
	FindArchivedVideoByUserId(context.Context, *FindArchivedVideoByUserIdRequest) (*FindArchivedVideoByUserIdResponse, error)
	FindUserVideo(context.Context, *FindUserVideoRequest) (*FindUserVideoResponse, error)
	ArchiveVideo(context.Context, *ArchiveVideoRequest) (*ArchiveVideoResponse, error)
	FetchAllVideo(context.Context, *FetchAllVideoRequest) (*FetchAllVideoResponse, error)
	GetVideoById(context.Context, *GetVideoByIdRequest) (*GetVideoByIdResponse, error)
	ToggleStar(context.Context, *ToggleStarRequest) (*ToggleStarResponse, error)
	mustEmbedUnimplementedVideoServiceServer()
}

// UnimplementedVideoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedVideoServiceServer struct {
}

func (UnimplementedVideoServiceServer) HealthCheck(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedVideoServiceServer) UploadVideo(VideoService_UploadVideoServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadVideo not implemented")
}
func (UnimplementedVideoServiceServer) FindArchivedVideoByUserId(context.Context, *FindArchivedVideoByUserIdRequest) (*FindArchivedVideoByUserIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindArchivedVideoByUserId not implemented")
}
func (UnimplementedVideoServiceServer) FindUserVideo(context.Context, *FindUserVideoRequest) (*FindUserVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindUserVideo not implemented")
}
func (UnimplementedVideoServiceServer) ArchiveVideo(context.Context, *ArchiveVideoRequest) (*ArchiveVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArchiveVideo not implemented")
}
func (UnimplementedVideoServiceServer) FetchAllVideo(context.Context, *FetchAllVideoRequest) (*FetchAllVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchAllVideo not implemented")
}
func (UnimplementedVideoServiceServer) GetVideoById(context.Context, *GetVideoByIdRequest) (*GetVideoByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVideoById not implemented")
}
func (UnimplementedVideoServiceServer) ToggleStar(context.Context, *ToggleStarRequest) (*ToggleStarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToggleStar not implemented")
}
func (UnimplementedVideoServiceServer) mustEmbedUnimplementedVideoServiceServer() {}

// UnsafeVideoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VideoServiceServer will
// result in compilation errors.
type UnsafeVideoServiceServer interface {
	mustEmbedUnimplementedVideoServiceServer()
}

func RegisterVideoServiceServer(s grpc.ServiceRegistrar, srv VideoServiceServer) {
	s.RegisterService(&VideoService_ServiceDesc, srv)
}

func _VideoService_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.VideoService/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).HealthCheck(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_UploadVideo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(VideoServiceServer).UploadVideo(&videoServiceUploadVideoServer{stream})
}

type VideoService_UploadVideoServer interface {
	SendAndClose(*UploadVideoResponse) error
	Recv() (*UploadVideoRequest, error)
	grpc.ServerStream
}

type videoServiceUploadVideoServer struct {
	grpc.ServerStream
}

func (x *videoServiceUploadVideoServer) SendAndClose(m *UploadVideoResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *videoServiceUploadVideoServer) Recv() (*UploadVideoRequest, error) {
	m := new(UploadVideoRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _VideoService_FindArchivedVideoByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindArchivedVideoByUserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).FindArchivedVideoByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.VideoService/FindArchivedVideoByUserId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).FindArchivedVideoByUserId(ctx, req.(*FindArchivedVideoByUserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_FindUserVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindUserVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).FindUserVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.VideoService/FindUserVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).FindUserVideo(ctx, req.(*FindUserVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_ArchiveVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArchiveVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).ArchiveVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.VideoService/ArchiveVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).ArchiveVideo(ctx, req.(*ArchiveVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_FetchAllVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchAllVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).FetchAllVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.VideoService/FetchAllVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).FetchAllVideo(ctx, req.(*FetchAllVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_GetVideoById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVideoByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).GetVideoById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.VideoService/GetVideoById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).GetVideoById(ctx, req.(*GetVideoByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_ToggleStar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ToggleStarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).ToggleStar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.VideoService/ToggleStar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).ToggleStar(ctx, req.(*ToggleStarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VideoService_ServiceDesc is the grpc.ServiceDesc for VideoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VideoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.VideoService",
	HandlerType: (*VideoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _VideoService_HealthCheck_Handler,
		},
		{
			MethodName: "FindArchivedVideoByUserId",
			Handler:    _VideoService_FindArchivedVideoByUserId_Handler,
		},
		{
			MethodName: "FindUserVideo",
			Handler:    _VideoService_FindUserVideo_Handler,
		},
		{
			MethodName: "ArchiveVideo",
			Handler:    _VideoService_ArchiveVideo_Handler,
		},
		{
			MethodName: "FetchAllVideo",
			Handler:    _VideoService_FetchAllVideo_Handler,
		},
		{
			MethodName: "GetVideoById",
			Handler:    _VideoService_GetVideoById_Handler,
		},
		{
			MethodName: "ToggleStar",
			Handler:    _VideoService_ToggleStar_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadVideo",
			Handler:       _VideoService_UploadVideo_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/proto/video.proto",
}

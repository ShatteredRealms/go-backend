// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: chat.proto

package pb

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

// ChatServiceClient is the client API for ChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatServiceClient interface {
	ConnectChannel(ctx context.Context, in *ChannelIdMessage, opts ...grpc.CallOption) (ChatService_ConnectChannelClient, error)
	ConnectDirectMessage(ctx context.Context, in *CharacterName, opts ...grpc.CallOption) (ChatService_ConnectDirectMessageClient, error)
	SendChatMessage(ctx context.Context, in *SendChatMessageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SendDirectMessage(ctx context.Context, in *SendDirectMessageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetChannel(ctx context.Context, in *ChannelIdMessage, opts ...grpc.CallOption) (*ChatChannel, error)
	CreateChannel(ctx context.Context, in *CreateChannelMessage, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteChannel(ctx context.Context, in *ChannelIdMessage, opts ...grpc.CallOption) (*emptypb.Empty, error)
	EditChannel(ctx context.Context, in *UpdateChatChannelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	AllChatChannels(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ChatChannels, error)
	GetAuthorizedChatChannels(ctx context.Context, in *RequestAuthorizedChatChannels, opts ...grpc.CallOption) (*ChatChannels, error)
	AuthorizeUserForChatChannel(ctx context.Context, in *RequestChatChannelAuthChange, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeauthorizeUserForChatChannel(ctx context.Context, in *RequestChatChannelAuthChange, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) ConnectChannel(ctx context.Context, in *ChannelIdMessage, opts ...grpc.CallOption) (ChatService_ConnectChannelClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChatService_ServiceDesc.Streams[0], "/sro.chat.ChatService/ConnectChannel", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceConnectChannelClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChatService_ConnectChannelClient interface {
	Recv() (*ChatMessage, error)
	grpc.ClientStream
}

type chatServiceConnectChannelClient struct {
	grpc.ClientStream
}

func (x *chatServiceConnectChannelClient) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chatServiceClient) ConnectDirectMessage(ctx context.Context, in *CharacterName, opts ...grpc.CallOption) (ChatService_ConnectDirectMessageClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChatService_ServiceDesc.Streams[1], "/sro.chat.ChatService/ConnectDirectMessage", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceConnectDirectMessageClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChatService_ConnectDirectMessageClient interface {
	Recv() (*ChatMessage, error)
	grpc.ClientStream
}

type chatServiceConnectDirectMessageClient struct {
	grpc.ClientStream
}

func (x *chatServiceConnectDirectMessageClient) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chatServiceClient) SendChatMessage(ctx context.Context, in *SendChatMessageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.chat.ChatService/SendChatMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) SendDirectMessage(ctx context.Context, in *SendDirectMessageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.chat.ChatService/SendDirectMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetChannel(ctx context.Context, in *ChannelIdMessage, opts ...grpc.CallOption) (*ChatChannel, error) {
	out := new(ChatChannel)
	err := c.cc.Invoke(ctx, "/sro.chat.ChatService/GetChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) CreateChannel(ctx context.Context, in *CreateChannelMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.chat.ChatService/CreateChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) DeleteChannel(ctx context.Context, in *ChannelIdMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.chat.ChatService/DeleteChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) EditChannel(ctx context.Context, in *UpdateChatChannelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.chat.ChatService/EditChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) AllChatChannels(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ChatChannels, error) {
	out := new(ChatChannels)
	err := c.cc.Invoke(ctx, "/sro.chat.ChatService/AllChatChannels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetAuthorizedChatChannels(ctx context.Context, in *RequestAuthorizedChatChannels, opts ...grpc.CallOption) (*ChatChannels, error) {
	out := new(ChatChannels)
	err := c.cc.Invoke(ctx, "/sro.chat.ChatService/GetAuthorizedChatChannels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) AuthorizeUserForChatChannel(ctx context.Context, in *RequestChatChannelAuthChange, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.chat.ChatService/AuthorizeUserForChatChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) DeauthorizeUserForChatChannel(ctx context.Context, in *RequestChatChannelAuthChange, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/sro.chat.ChatService/DeauthorizeUserForChatChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServiceServer is the server API for ChatService service.
// All implementations must embed UnimplementedChatServiceServer
// for forward compatibility
type ChatServiceServer interface {
	ConnectChannel(*ChannelIdMessage, ChatService_ConnectChannelServer) error
	ConnectDirectMessage(*CharacterName, ChatService_ConnectDirectMessageServer) error
	SendChatMessage(context.Context, *SendChatMessageRequest) (*emptypb.Empty, error)
	SendDirectMessage(context.Context, *SendDirectMessageRequest) (*emptypb.Empty, error)
	GetChannel(context.Context, *ChannelIdMessage) (*ChatChannel, error)
	CreateChannel(context.Context, *CreateChannelMessage) (*emptypb.Empty, error)
	DeleteChannel(context.Context, *ChannelIdMessage) (*emptypb.Empty, error)
	EditChannel(context.Context, *UpdateChatChannelRequest) (*emptypb.Empty, error)
	AllChatChannels(context.Context, *emptypb.Empty) (*ChatChannels, error)
	GetAuthorizedChatChannels(context.Context, *RequestAuthorizedChatChannels) (*ChatChannels, error)
	AuthorizeUserForChatChannel(context.Context, *RequestChatChannelAuthChange) (*emptypb.Empty, error)
	DeauthorizeUserForChatChannel(context.Context, *RequestChatChannelAuthChange) (*emptypb.Empty, error)
	mustEmbedUnimplementedChatServiceServer()
}

// UnimplementedChatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChatServiceServer struct {
}

func (UnimplementedChatServiceServer) ConnectChannel(*ChannelIdMessage, ChatService_ConnectChannelServer) error {
	return status.Errorf(codes.Unimplemented, "method ConnectChannel not implemented")
}
func (UnimplementedChatServiceServer) ConnectDirectMessage(*CharacterName, ChatService_ConnectDirectMessageServer) error {
	return status.Errorf(codes.Unimplemented, "method ConnectDirectMessage not implemented")
}
func (UnimplementedChatServiceServer) SendChatMessage(context.Context, *SendChatMessageRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendChatMessage not implemented")
}
func (UnimplementedChatServiceServer) SendDirectMessage(context.Context, *SendDirectMessageRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendDirectMessage not implemented")
}
func (UnimplementedChatServiceServer) GetChannel(context.Context, *ChannelIdMessage) (*ChatChannel, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChannel not implemented")
}
func (UnimplementedChatServiceServer) CreateChannel(context.Context, *CreateChannelMessage) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateChannel not implemented")
}
func (UnimplementedChatServiceServer) DeleteChannel(context.Context, *ChannelIdMessage) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteChannel not implemented")
}
func (UnimplementedChatServiceServer) EditChannel(context.Context, *UpdateChatChannelRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditChannel not implemented")
}
func (UnimplementedChatServiceServer) AllChatChannels(context.Context, *emptypb.Empty) (*ChatChannels, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllChatChannels not implemented")
}
func (UnimplementedChatServiceServer) GetAuthorizedChatChannels(context.Context, *RequestAuthorizedChatChannels) (*ChatChannels, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthorizedChatChannels not implemented")
}
func (UnimplementedChatServiceServer) AuthorizeUserForChatChannel(context.Context, *RequestChatChannelAuthChange) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthorizeUserForChatChannel not implemented")
}
func (UnimplementedChatServiceServer) DeauthorizeUserForChatChannel(context.Context, *RequestChatChannelAuthChange) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeauthorizeUserForChatChannel not implemented")
}
func (UnimplementedChatServiceServer) mustEmbedUnimplementedChatServiceServer() {}

// UnsafeChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServiceServer will
// result in compilation errors.
type UnsafeChatServiceServer interface {
	mustEmbedUnimplementedChatServiceServer()
}

func RegisterChatServiceServer(s grpc.ServiceRegistrar, srv ChatServiceServer) {
	s.RegisterService(&ChatService_ServiceDesc, srv)
}

func _ChatService_ConnectChannel_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ChannelIdMessage)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServiceServer).ConnectChannel(m, &chatServiceConnectChannelServer{stream})
}

type ChatService_ConnectChannelServer interface {
	Send(*ChatMessage) error
	grpc.ServerStream
}

type chatServiceConnectChannelServer struct {
	grpc.ServerStream
}

func (x *chatServiceConnectChannelServer) Send(m *ChatMessage) error {
	return x.ServerStream.SendMsg(m)
}

func _ChatService_ConnectDirectMessage_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CharacterName)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServiceServer).ConnectDirectMessage(m, &chatServiceConnectDirectMessageServer{stream})
}

type ChatService_ConnectDirectMessageServer interface {
	Send(*ChatMessage) error
	grpc.ServerStream
}

type chatServiceConnectDirectMessageServer struct {
	grpc.ServerStream
}

func (x *chatServiceConnectDirectMessageServer) Send(m *ChatMessage) error {
	return x.ServerStream.SendMsg(m)
}

func _ChatService_SendChatMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendChatMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).SendChatMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.chat.ChatService/SendChatMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).SendChatMessage(ctx, req.(*SendChatMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_SendDirectMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendDirectMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).SendDirectMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.chat.ChatService/SendDirectMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).SendDirectMessage(ctx, req.(*SendDirectMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChannelIdMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.chat.ChatService/GetChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetChannel(ctx, req.(*ChannelIdMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_CreateChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateChannelMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).CreateChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.chat.ChatService/CreateChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).CreateChannel(ctx, req.(*CreateChannelMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_DeleteChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChannelIdMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).DeleteChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.chat.ChatService/DeleteChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).DeleteChannel(ctx, req.(*ChannelIdMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_EditChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateChatChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).EditChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.chat.ChatService/EditChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).EditChannel(ctx, req.(*UpdateChatChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_AllChatChannels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).AllChatChannels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.chat.ChatService/AllChatChannels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).AllChatChannels(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetAuthorizedChatChannels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestAuthorizedChatChannels)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetAuthorizedChatChannels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.chat.ChatService/GetAuthorizedChatChannels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetAuthorizedChatChannels(ctx, req.(*RequestAuthorizedChatChannels))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_AuthorizeUserForChatChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestChatChannelAuthChange)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).AuthorizeUserForChatChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.chat.ChatService/AuthorizeUserForChatChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).AuthorizeUserForChatChannel(ctx, req.(*RequestChatChannelAuthChange))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_DeauthorizeUserForChatChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestChatChannelAuthChange)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).DeauthorizeUserForChatChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sro.chat.ChatService/DeauthorizeUserForChatChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).DeauthorizeUserForChatChannel(ctx, req.(*RequestChatChannelAuthChange))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatService_ServiceDesc is the grpc.ServiceDesc for ChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sro.chat.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendChatMessage",
			Handler:    _ChatService_SendChatMessage_Handler,
		},
		{
			MethodName: "SendDirectMessage",
			Handler:    _ChatService_SendDirectMessage_Handler,
		},
		{
			MethodName: "GetChannel",
			Handler:    _ChatService_GetChannel_Handler,
		},
		{
			MethodName: "CreateChannel",
			Handler:    _ChatService_CreateChannel_Handler,
		},
		{
			MethodName: "DeleteChannel",
			Handler:    _ChatService_DeleteChannel_Handler,
		},
		{
			MethodName: "EditChannel",
			Handler:    _ChatService_EditChannel_Handler,
		},
		{
			MethodName: "AllChatChannels",
			Handler:    _ChatService_AllChatChannels_Handler,
		},
		{
			MethodName: "GetAuthorizedChatChannels",
			Handler:    _ChatService_GetAuthorizedChatChannels_Handler,
		},
		{
			MethodName: "AuthorizeUserForChatChannel",
			Handler:    _ChatService_AuthorizeUserForChatChannel_Handler,
		},
		{
			MethodName: "DeauthorizeUserForChatChannel",
			Handler:    _ChatService_DeauthorizeUserForChatChannel_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ConnectChannel",
			Handler:       _ChatService_ConnectChannel_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ConnectDirectMessage",
			Handler:       _ChatService_ConnectDirectMessage_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chat.proto",
}

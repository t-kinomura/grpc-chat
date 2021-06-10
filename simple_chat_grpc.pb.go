// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc_chat

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

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatClient interface {
	SimpleChat(ctx context.Context, opts ...grpc.CallOption) (Chat_SimpleChatClient, error)
}

type chatClient struct {
	cc grpc.ClientConnInterface
}

func NewChatClient(cc grpc.ClientConnInterface) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) SimpleChat(ctx context.Context, opts ...grpc.CallOption) (Chat_SimpleChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &Chat_ServiceDesc.Streams[0], "/chat.Chat/SimpleChat", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatSimpleChatClient{stream}
	return x, nil
}

type Chat_SimpleChatClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type chatSimpleChatClient struct {
	grpc.ClientStream
}

func (x *chatSimpleChatClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatSimpleChatClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatServer is the server API for Chat service.
// All implementations must embed UnimplementedChatServer
// for forward compatibility
type ChatServer interface {
	SimpleChat(Chat_SimpleChatServer) error
	mustEmbedUnimplementedChatServer()
}

// UnimplementedChatServer must be embedded to have forward compatible implementations.
type UnimplementedChatServer struct {
}

func (UnimplementedChatServer) SimpleChat(Chat_SimpleChatServer) error {
	return status.Errorf(codes.Unimplemented, "method SimpleChat not implemented")
}
func (UnimplementedChatServer) mustEmbedUnimplementedChatServer() {}

// UnsafeChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServer will
// result in compilation errors.
type UnsafeChatServer interface {
	mustEmbedUnimplementedChatServer()
}

func RegisterChatServer(s grpc.ServiceRegistrar, srv ChatServer) {
	s.RegisterService(&Chat_ServiceDesc, srv)
}

func _Chat_SimpleChat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServer).SimpleChat(&chatSimpleChatServer{stream})
}

type Chat_SimpleChatServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type chatSimpleChatServer struct {
	grpc.ServerStream
}

func (x *chatSimpleChatServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatSimpleChatServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Chat_ServiceDesc is the grpc.ServiceDesc for Chat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SimpleChat",
			Handler:       _Chat_SimpleChat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "simple_chat.proto",
}
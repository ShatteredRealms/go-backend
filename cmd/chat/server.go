package main

import (
	"context"
	chat "github.com/WilSimpson/ShatteredRealms/go-backend/cmd/chat/global"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/srv"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func NewServer(
	ctx context.Context,
	jwt service.JWTService,
	chatService service.ChatService,
) (*grpc.Server, *runtime.ServeMux, error) {
	var span trace.Span
	ctx, span = chat.Tracer.Start(ctx, "new-chat-server")

	grpcServer, gwmux, opts, err := srv.CreateGrpcServerWithAuth(
		ctx,
		jwt,
		conf.Accounts.Remote.Address(),
		"chat",
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	charactersConn, err := grpc.Dial(
		conf.Characters.Remote.Address(),
		srv.InsecureOtelGrpcDialOpts()...,
	)
	if err != nil {
		return nil, nil, err
	}

	csc := pb.NewCharactersServiceClient(charactersConn)

	chatServiceServer := srv.NewChatServiceServer(chatService, jwt, csc)
	pb.RegisterChatServiceServer(grpcServer, chatServiceServer)
	err = pb.RegisterChatServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Chat.Local.Address(),
		opts,
	)

	span.End()

	return grpcServer, gwmux, nil
}

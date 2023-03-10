package srv

import (
	"context"
	chat "github.com/ShatteredRealms/go-backend/cmd/chat/app"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type chatServiceServer struct {
	pb.UnimplementedChatServiceServer
	server *chat.ChatServerContext
}

func (s chatServiceServer) ConnectChannel(target *pb.ChannelTarget, server pb.ChatService_ConnectChannelServer) error {
	r := s.server.ChatService.ChannelMessagesReader(server.Context(), uint(target.ChannelId))

	for {
		msg, err := r.ReadMessage(server.Context())
		if err != nil {
			_ = r.Close()
			return err
		}

		err = server.Send(&pb.ChatMessage{
			CharacterName: string(msg.Key),
			Message:       string(msg.Value),
		})

		if err != nil {
			_ = r.Close()
			return err
		}
	}
}

func (s chatServiceServer) ConnectDirectMessage(name *pb.CharacterName, server pb.ChatService_ConnectDirectMessageServer) error {
	r := s.server.ChatService.DirectMessagesReader(server.Context(), name.Character)
	for {
		msg, err := r.ReadMessage(server.Context())
		if err != nil {
			_ = r.Close()
			return err
		}

		err = server.Send(&pb.ChatMessage{
			CharacterName: string(msg.Key),
			Message:       string(msg.Value),
		})

		if err != nil {
			_ = r.Close()
			return err
		}
	}
}

func (s chatServiceServer) SendChatMessage(ctx context.Context, request *pb.SendChatMessageRequest) (*emptypb.Empty, error) {
	if err := s.verifyUserOwnsCharacter(ctx, request.ChatMessage.CharacterName); err != nil {
		return nil, err
	}

	if err := s.server.ChatService.SendChannelMessage(ctx, request.ChatMessage.CharacterName, request.ChatMessage.Message, uint(request.ChannelId)); err != nil {
		log.WithContext(ctx).Errorf("send channel chat message: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to send message")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) SendDirectMessage(ctx context.Context, request *pb.SendDirectMessageRequest) (*emptypb.Empty, error) {
	if err := s.verifyUserOwnsCharacter(ctx, request.ChatMessage.CharacterName); err != nil {
		return nil, err
	}

	if err := s.server.ChatService.SendDirectMessage(ctx, request.ChatMessage.CharacterName, request.ChatMessage.Message, request.Character); err != nil {
		log.WithContext(ctx).Errorf("send direct chat message: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to send message")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) GetChannel(ctx context.Context, target *pb.ChannelTarget) (*pb.ChatChannel, error) {
	//TODO implement me
	panic("implement me")
}

func (s chatServiceServer) CreateChannel(ctx context.Context, message *pb.CreateChannelMessage) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s chatServiceServer) DeleteChannel(ctx context.Context, target *pb.ChannelTarget) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s chatServiceServer) EditChannel(ctx context.Context, request *pb.UpdateChatChannelRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s chatServiceServer) AllChatChannels(ctx context.Context, empty *emptypb.Empty) (*pb.ChatChannels, error) {
	//TODO implement me
	panic("implement me")
}

func (s chatServiceServer) GetAuthorizedChatChannels(ctx context.Context, channels *pb.RequestAuthorizedChatChannels) (*pb.ChatChannels, error) {
	//TODO implement me
	panic("implement me")
}

func (s chatServiceServer) AuthorizeUserForChatChannel(ctx context.Context, change *pb.RequestChatChannelAuthChange) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s chatServiceServer) DeauthorizeUserForChatChannel(ctx context.Context, change *pb.RequestChatChannelAuthChange) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s chatServiceServer) mustEmbedUnimplementedChatServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewChatServiceServer(server *chat.ChatServerContext) pb.ChatServiceServer {
	return &chatServiceServer{
		server: server,
	}
}

func (s chatServiceServer) verifyUserOwnsCharacter(ctx context.Context, characterName string) error {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "authentication required")
	}

	serverCtx := helpers.ContextAddClientAuth(ctx, s.server.GlobalConfig.Chat.Keycloak.ClientId, s.server.GlobalConfig.Chat.Keycloak.ClientSecret)
	chars, err := s.server.CharacterService.GetAllCharactersForUser(serverCtx, &pb.UserTarget{UserId: claims.Subject})
	if err != nil {
		log.WithContext(ctx).Errorf("chat character service get for user: %v", err)
		return status.Errorf(codes.Internal, "unable to verify character")
	}

	for _, c := range chars.Characters {
		if c.Name == characterName {
			return nil
		}
	}

	return status.Errorf(codes.Unauthenticated, "character not found")
}

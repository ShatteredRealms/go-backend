package srv

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type characterTarget struct {
	Id   uint64
	Name string
}

type chatServiceServer struct {
	pb.UnimplementedChatServiceServer
	chatService service.ChatService
	jwtService  service.JWTService

	charactersServiceClient pb.CharactersServiceClient
}

func (s chatServiceServer) ConnectChannel(request *pb.ChannelIdMessage, server pb.ChatService_ConnectChannelServer) error {
	r := s.chatService.ChannelMessagesReader(server.Context(), uint(request.ChannelId))

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

func (s chatServiceServer) SendChatMessage(ctx context.Context, request *pb.SendChatMessageRequest) (*empty.Empty, error) {
	if err := s.verifyUserOwnsCharacter(ctx, characterTarget{Name: request.ChatMessage.CharacterName}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, s.chatService.SendChannelMessage(ctx, request.ChatMessage.CharacterName, request.ChatMessage.Message, uint(request.ChannelId))
}

func (s chatServiceServer) SendDirectMessage(ctx context.Context, request *pb.SendDirectMessageRequest) (*empty.Empty, error) {
	// Only the owning user can send messages. Do not allow even super admins to spoof.
	if err := s.verifyUserOwnsCharacter(ctx, characterTarget{Name: request.ChatMessage.CharacterName}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, s.chatService.SendDirectMessage(ctx, request.ChatMessage.CharacterName, request.ChatMessage.Message, request.CharacterName)
}

func (s chatServiceServer) ConnectDirectMessage(msg *pb.CharacterName, srv pb.ChatService_ConnectDirectMessageServer) error {
	if !interceptor.AuthorizedForOther(srv.Context()) {
		if err := s.verifyUserOwnsCharacter(srv.Context(), characterTarget{Name: msg.CharacterName}); err != nil {
			return err
		}
	}

	r := s.chatService.DirectMessagesReader(srv.Context(), msg.CharacterName)
	for {
		msg, err := r.ReadMessage(srv.Context())
		if err != nil {
			_ = r.Close()
			return err
		}

		err = srv.Send(&pb.ChatMessage{
			CharacterName: string(msg.Key),
			Message:       string(msg.Value),
		})

		if err != nil {
			_ = r.Close()
			return err
		}
	}
}

func (s chatServiceServer) CreateChannel(ctx context.Context, msg *pb.CreateChannelMessage) (*empty.Empty, error) {
	_, err := s.chatService.CreateChannel(ctx, &model.ChatChannel{
		Name:   msg.Name,
		Public: msg.Public,
	})

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) DeleteChannel(ctx context.Context, msg *pb.ChannelIdMessage) (*empty.Empty, error) {
	err := s.chatService.DeleteChannel(ctx, &model.ChatChannel{
		Model: gorm.Model{ID: uint(msg.ChannelId)},
	})

	return &emptypb.Empty{}, err
}

func (s chatServiceServer) AllChatChannels(ctx context.Context, _ *empty.Empty) (*pb.ChatChannels, error) {
	channels, err := s.chatService.AllChannels(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.ChatChannels{
		Channels: make([]*pb.ChatChannel, len(channels)),
	}

	for i, c := range channels {
		resp.Channels[i] = c.ToPb()
	}

	return resp, nil
}

func (s chatServiceServer) GetChannel(ctx context.Context, msg *pb.ChannelIdMessage) (*pb.ChatChannel, error) {
	channel, err := s.chatService.GetChannel(ctx, uint(msg.ChannelId))
	if err != nil {
		return nil, err
	}

	return channel.ToPb(), nil
}

func (s chatServiceServer) EditChannel(ctx context.Context, msg *pb.UpdateChatChannelRequest) (*empty.Empty, error) {
	return &emptypb.Empty{}, s.chatService.UpdateChannel(ctx, msg)
}

func (s chatServiceServer) GetAuthorizedChatChannels(ctx context.Context, msg *pb.RequestAuthorizedChatChannels) (*pb.ChatChannels, error) {
	if !interceptor.AuthorizedForOther(ctx) {
		if err := s.verifyUserOwnsCharacter(ctx, characterTarget{Id: msg.CharacterId}); err != nil {
			return nil, err
		}
	}

	channels, err := s.chatService.AuthorizedChannelsForCharacter(ctx, msg.CharacterId)
	if err != nil {
		return nil, err
	}

	resp := &pb.ChatChannels{
		Channels: make([]*pb.ChatChannel, len(channels)),
	}

	for i, c := range channels {
		resp.Channels[i] = c.ToPb()
	}

	return resp, nil
}

func (s chatServiceServer) AuthorizeUserForChatChannel(ctx context.Context, msg *pb.RequestChatChannelAuthChange) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.chatService.ChangeAuthorizationForCharacter(ctx, msg.CharacterId, msg.Ids, true)
}

func (s chatServiceServer) DeauthorizeUserForChatChannel(ctx context.Context, msg *pb.RequestChatChannelAuthChange) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.chatService.ChangeAuthorizationForCharacter(ctx, msg.CharacterId, msg.Ids, false)
}

// verifyUserOwnsCharacter returns an error if the user doesn't own the given characterName or characterId
func (s chatServiceServer) verifyUserOwnsCharacter(ctx context.Context, target characterTarget) error {
	userId, err := interceptor.ExtractSubject(ctx, s.jwtService)
	if err != nil {
		return status.Error(codes.Internal, "Unable to parse auth token")
	}

	characters, err := s.charactersServiceClient.GetAllCharactersForUser(serverAuthContext(ctx, s.jwtService, "sro.com/chat"), &pb.UserTarget{UserId: uint64(userId)})
	if err != nil {
		return status.Error(codes.Internal, "Unable to parse auth token")
	}

	found := false
	for _, character := range characters.Characters {
		if character.Name.Value == target.Name || character.Id == target.Id {
			found = true
			break
		}
	}

	if !found {
		return status.Error(codes.Unauthenticated, "Unauthorized for this character")
	}

	return nil
}

func NewChatServiceServer(chatService service.ChatService, jwtService service.JWTService, csc pb.CharactersServiceClient) pb.ChatServiceServer {
	return chatServiceServer{
		chatService:             chatService,
		jwtService:              jwtService,
		charactersServiceClient: csc,
	}
}

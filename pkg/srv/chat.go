package srv

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	chat "github.com/ShatteredRealms/go-backend/cmd/chat/app"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type chatServiceServer struct {
	pb.UnimplementedChatServiceServer
	server *chat.ChatServerContext
}

var (
	ChatRoles = make([]*gocloak.Role, 0)

	RoleChat = registerChatRole(&gocloak.Role{
		Name:        gocloak.StringP("chat"),
		Description: gocloak.StringP("Allows getting and listening to chat channels and direct messages to  self as well as sending messages on chat channels and to other users. Does not give permissions to a particular channel."),
	})

	RoleChatChannelManage = registerChatRole(&gocloak.Role{
		Name:        gocloak.StringP("chat_manage"),
		Description: gocloak.StringP("Allows viewin, creation, editing and deletion of chat channels."),
	})
)

func registerChatRole(role *gocloak.Role) *gocloak.Role {
	ChatRoles = append(ChatRoles, role)
	return role
}

func (s chatServiceServer) ConnectChannel(target *pb.ChannelTarget, server pb.ChatService_ConnectChannelServer) error {
	claims, err := helpers.ExtractClaims(server.Context())
	if err != nil {
		return ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChat, model.ChatClientId) {
		return ErrNotAuthorized
	}

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
	claims, err := helpers.ExtractClaims(server.Context())
	if err != nil {
		return ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChat, model.ChatClientId) {
		return ErrNotAuthorized
	}

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
	claims, err := helpers.ExtractClaims(server.Context())
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChat, model.ChatClientId) {
		return nil, ErrNotAuthorized
	}

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
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChat, model.ChatClientId) {
		return nil, ErrNotAuthorized
	}

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
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, ErrNotAuthorized
	}

	c, err := s.server.ChatService.GetChannel(ctx, uint(target.ChannelId))
	if err != nil {
		log.WithContext(ctx).Errorf("get chat channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to get chat channel")
	}

	return c.ToPb(), nil
}

func (s chatServiceServer) CreateChannel(ctx context.Context, message *pb.CreateChannelMessage) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, ErrNotAuthorized
	}

	_, err = s.server.ChatService.CreateChannel(
		ctx,
		&model.ChatChannel{
			Name:      message.Name,
			Dimension: message.Dimension,
		},
	)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Error(codes.InvalidArgument, "name is already taken in this dimension")
		}

		log.WithContext(ctx).Errorf("create chat channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to create chat channel")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) DeleteChannel(ctx context.Context, target *pb.ChannelTarget) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, ErrNotAuthorized
	}

	err = s.server.ChatService.DeleteChannel(ctx, &model.ChatChannel{
		Model: gorm.Model{ID: uint(target.ChannelId)},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDoesNotExist
		}

		log.WithContext(ctx).Errorf("delete channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to delete channel")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) EditChannel(ctx context.Context, request *pb.UpdateChatChannelRequest) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, ErrNotAuthorized
	}

	err = s.server.ChatService.UpdateChannel(ctx, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDoesNotExist
		}

		log.WithContext(ctx).Errorf("edit channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to edit channel")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) AllChatChannels(ctx context.Context, empty *emptypb.Empty) (*pb.ChatChannels, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, ErrNotAuthorized
	}

	channels, err := s.server.ChatService.AllChannels()
	if err != nil {
		log.WithContext(ctx).Errorf("edit channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to get channels")
	}

	return channels.ToPb(), nil
}

func (s chatServiceServer) GetAuthorizedChatChannels(ctx context.Context, request *pb.RequestAuthorizedChatChannels) (*pb.ChatChannels, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChat, model.ChatClientId) {
		return nil, ErrNotAuthorized
	}

	channels, err := s.server.ChatService.AuthorizedChannelsForCharacter(ctx, request.Character)
	if err != nil {
		log.WithContext(ctx).Errorf("get authorized channels: %v", err)
		return nil, status.Error(codes.Internal, "unable to get channels")
	}

	return channels.ToPb(), nil
}

func (s chatServiceServer) AuthorizeUserForChatChannel(ctx context.Context, request *pb.RequestChatChannelAuthChange) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, ErrNotAuthorized
	}

	err = s.server.ChatService.ChangeAuthorizationForCharacter(
		ctx,
		request.Character,
		*helpers.ArrayOfUint64ToUint(&request.Ids),
		true,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Error(codes.InvalidArgument, "existing chat authorization conflict")
		}

		log.WithContext(ctx).Errorf("get authorized channels: %v", err)
		return nil, status.Error(codes.Internal, "unable to change authorization")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) DeauthorizeUserForChatChannel(ctx context.Context, request *pb.RequestChatChannelAuthChange) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, ErrNotAuthorized
	}

	err = s.server.ChatService.ChangeAuthorizationForCharacter(
		ctx,
		request.Character,
		*helpers.ArrayOfUint64ToUint(&request.Ids),
		true,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDoesNotExist
		}

		log.WithContext(ctx).Errorf("get authorized channels: %v", err)
		return nil, status.Error(codes.Internal, "unable to change authorization")
	}

	return &emptypb.Empty{}, nil
}

func NewChatServiceServer(
	ctx context.Context,
	server *chat.ChatServerContext,
) (pb.ChatServiceServer, error) {
	token, err := server.KeycloakClient.LoginClient(
		ctx,
		server.GlobalConfig.Chat.Keycloak.ClientId,
		server.GlobalConfig.Chat.Keycloak.ClientSecret,
		server.GlobalConfig.Chat.Keycloak.Realm,
	)
	if err != nil {
		return nil, fmt.Errorf("login keycloak: %v", err)
	}

	err = createRoles(ctx,
		token.AccessToken,
		server.GlobalConfig.Chat.Keycloak.Realm,
		server.GlobalConfig.Chat.Keycloak.Id,
		&ChatRoles,
	)
	if err != nil {
		return nil, err
	}

	return &chatServiceServer{
		server: server,
	}, nil
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

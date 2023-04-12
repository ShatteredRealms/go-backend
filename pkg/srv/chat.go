package srv

import (
	"context"
	"errors"
	"fmt"
	"reflect"

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
		Description: gocloak.StringP("Allows viewing, creation, editing and deletion of chat channels."),
	})
)

func registerChatRole(role *gocloak.Role) *gocloak.Role {
	ChatRoles = append(ChatRoles, role)
	return role
}

func (s chatServiceServer) ConnectChannel(
	request *pb.ChatChannelTarget,
	server pb.ChatService_ConnectChannelServer,
) error {
	claims, err := helpers.ExtractClaims(server.Context())
	if err != nil {
		log.WithContext(server.Context()).Infof("extract claims failed: %v", err)
		return model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChat, model.ChatClientId) {
		return model.ErrUnauthorized
	}

	r := s.server.ChatService.ChannelMessagesReader(server.Context(), uint(request.Id))

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

func (s chatServiceServer) ConnectDirectMessage(
	request *pb.CharacterTarget,
	server pb.ChatService_ConnectDirectMessageServer,
) error {
	claims, err := helpers.ExtractClaims(server.Context())
	if err != nil {
		log.WithContext(server.Context()).Infof("extract claims failed: %v", err)
		return model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChat, model.ChatClientId) {
		return model.ErrUnauthorized
	}

	char, err := s.verifyUserOwnsCharacter(server.Context(), request)
	if err != nil {
		return model.ErrUnauthorized
	}

	r := s.server.ChatService.DirectMessagesReader(server.Context(), char.Name)
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

func (s chatServiceServer) SendChatMessage(
	ctx context.Context,
	request *pb.SendChatMessageRequest,
) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Infof("extract claims failed: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChat, model.ChatClientId) {
		log.WithContext(ctx).Infof("unauthorized request")
		return nil, model.ErrUnauthorized
	}

	if _, err = s.verifyUserOwnsCharacter(
		ctx,
		&pb.CharacterTarget{Target: &pb.CharacterTarget_Name{Name: request.ChatMessage.CharacterName}},
	); err != nil {
		log.WithContext(ctx).Infof("verify owns character failed: %v", err)
		return nil, err
	}

	if err := s.server.ChatService.SendChannelMessage(
		ctx,
		request.ChatMessage.CharacterName,
		request.ChatMessage.Message,
		uint(request.ChannelId),
	); err != nil {
		log.WithContext(ctx).Errorf("send channel chat message: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to send message")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) SendDirectMessage(
	ctx context.Context,
	request *pb.SendDirectMessageRequest,
) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Infof("extract claims failed: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChat, model.ChatClientId) {
		return nil, model.ErrUnauthorized
	}

	if _, err = s.verifyUserOwnsCharacter(
		ctx,
		&pb.CharacterTarget{Target: &pb.CharacterTarget_Name{Name: request.ChatMessage.CharacterName}},
	); err != nil {
		return nil, err
	}

	srvCtx, err := s.serverContext(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("create server context: %v", err)
		return nil, model.ErrHandleRequest
	}
	targetCharacterName, err := helpers.GetCharacterNameFromTarget(srvCtx, s.server.CharacterService, request.Target)
	if err != nil {
		return nil, err
	}

	if err := s.server.ChatService.SendDirectMessage(
		ctx,
		request.ChatMessage.CharacterName,
		request.ChatMessage.Message,
		targetCharacterName,
	); err != nil {
		log.WithContext(ctx).Errorf("send direct chat message: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to send message")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) GetChannel(
	ctx context.Context,
	request *pb.ChatChannelTarget,
) (*pb.ChatChannel, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Infof("extract claims failed: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, model.ErrUnauthorized
	}

	c, err := s.server.ChatService.GetChannel(ctx, uint(request.Id))
	if err != nil {
		log.WithContext(ctx).Errorf("get chat channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to get chat channel")
	}

	return c.ToPb(), nil
}

func (s chatServiceServer) CreateChannel(
	ctx context.Context,
	request *pb.CreateChannelMessage,
) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Infof("extract claims failed: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, model.ErrUnauthorized
	}

	_, err = s.server.ChatService.CreateChannel(
		ctx,
		&model.ChatChannel{
			Name:      request.Name,
			Dimension: request.Dimension,
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

func (s chatServiceServer) DeleteChannel(
	ctx context.Context,
	request *pb.ChatChannelTarget,
) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Infof("extract claims failed: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, model.ErrUnauthorized
	}

	err = s.server.ChatService.DeleteChannel(ctx, &model.ChatChannel{
		Model: gorm.Model{ID: uint(request.Id)},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrDoesNotExist
		}

		log.WithContext(ctx).Errorf("delete channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to delete channel")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) EditChannel(
	ctx context.Context,
	request *pb.UpdateChatChannelRequest,
) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Infof("extract claims failed: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, model.ErrUnauthorized
	}

	err = s.server.ChatService.UpdateChannel(ctx, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrDoesNotExist
		}

		log.WithContext(ctx).Errorf("edit channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to edit channel")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) AllChatChannels(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.ChatChannels, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Infof("extract claims failed: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, model.ErrUnauthorized
	}

	channels, err := s.server.ChatService.AllChannels(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("edit channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to get channels")
	}

	return channels.ToPb(), nil
}

func (s chatServiceServer) GetAuthorizedChatChannels(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*pb.ChatChannels, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Infof("extract claims failed: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChat, model.ChatClientId) {
		return nil, model.ErrUnauthorized
	}

	srvCtx, err := s.serverContext(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("create server context: %v", err)
		return nil, model.ErrHandleRequest
	}
	targetCharacterId, err := helpers.GetCharacterIdFromTarget(srvCtx, s.server.CharacterService, request)
	if err != nil {
		return nil, err
	}

	channels, err := s.server.ChatService.AuthorizedChannelsForCharacter(ctx, targetCharacterId)
	if err != nil {
		log.WithContext(ctx).Errorf("get authorized channels: %v", err)
		return nil, status.Error(codes.Internal, "unable to get channels")
	}

	return channels.ToPb(), nil
}

func (s chatServiceServer) UpdateUserChatChannelAuthorizations(
	ctx context.Context,
	request *pb.RequestChatChannelAuthChange,
) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Infof("extract claims failed: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, model.ChatClientId) {
		return nil, model.ErrUnauthorized
	}

	targetCharacterId := uint(0)
	switch target := request.Character.Target.(type) {
	case *pb.CharacterTarget_Name:
		srvCtx, err := s.serverContext(ctx)
		if err != nil {
			log.WithContext(ctx).Errorf("create server context: %v", err)
			return nil, model.ErrHandleRequest
		}
		targetChar, err := s.server.CharacterService.GetCharacter(srvCtx, request.Character)
		if err != nil {
			return nil, err
		}
		targetCharacterId = uint(targetChar.Id)

	case *pb.CharacterTarget_Id:
		targetCharacterId = uint(target.Id)

	default:
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(request.Character).Name())
		return nil, model.ErrHandleRequest
	}

	err = s.server.ChatService.ChangeAuthorizationForCharacter(
		ctx,
		targetCharacterId,
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
		server.KeycloakClient,
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

func (s chatServiceServer) serverContext(ctx context.Context) (context.Context, error) {
	token, err := s.server.KeycloakClient.LoginClient(
		ctx,
		s.server.GlobalConfig.Chat.Keycloak.ClientId,
		s.server.GlobalConfig.Chat.Keycloak.ClientSecret,
		s.server.GlobalConfig.Chat.Keycloak.Realm,
	)
	if err != nil {
		log.WithContext(ctx).Errorf("keycloak login failure: %v", err)
		return nil, err
	}

	return helpers.ContextAddClientToken(
		ctx,
		token.AccessToken,
	), nil
}

func (s chatServiceServer) verifyUserOwnsCharacter(ctx context.Context, request *pb.CharacterTarget) (*pb.CharacterResponse, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Infof("extract claims failed: %v", err)
		return nil, model.ErrUnauthorized
	}

	srvCtx, err := s.serverContext(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("create server context: %v", err)
		return nil, model.ErrHandleRequest
	}

	chars, err := s.server.CharacterService.GetAllCharactersForUser(srvCtx, &pb.UserTarget{
		Target: &pb.UserTarget_Id{Id: claims.Subject},
	})
	if err != nil {
		log.WithContext(ctx).Errorf("chat character service get for user: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to verify character")
	}

	switch target := request.Target.(type) {
	case *pb.CharacterTarget_Id:
		for _, c := range chars.Characters {
			if c.Id == target.Id {
				return c, nil
			}
		}

	case *pb.CharacterTarget_Name:
		for _, c := range chars.Characters {
			if c.Name == target.Name {
				return c, nil
			}
		}

	default:
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(request.Target).Name())
		return nil, model.ErrHandleRequest
	}

	return nil, status.Errorf(codes.Unauthenticated, "character not found")
}

package srv

import (
	"context"
	"errors"

	chatApp "github.com/ShatteredRealms/go-backend/cmd/chat/app"
	"github.com/ShatteredRealms/go-backend/pkg/auth"
	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model/character"
	"github.com/ShatteredRealms/go-backend/pkg/model/chat"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/gocloak/v13"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type chatServiceServer struct {
	pb.UnimplementedChatServiceServer
	server *chatApp.ChatServerContext
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
	claims, ok := auth.RetrieveClaims(server.Context())
	if !ok {
		return common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChat, auth.ChatClientId) {
		return common.ErrUnauthorized.Err()
	}

	// Validate requester has chat channel permissions
	if !claims.HasResourceRole(RoleChatChannelManage, auth.ChatClientId) {
		err := s.checkUserChannelAuth(server.Context(), claims.ID, uint(request.Id))
		if err != nil {
			log.Logger.WithContext(server.Context()).Infof("verify auth failed for %s on channel %d: %v", claims.ID, request.Id, err)
			return err
		}
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
	claims, ok := auth.RetrieveClaims(server.Context())
	if !ok {
		return common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChat, auth.ChatClientId) {
		return common.ErrUnauthorized.Err()
	}

	char, err := s.verifyUserOwnsCharacter(server.Context(), request)
	if err == common.ErrNotOwner {
		if !claims.HasResourceRole(RoleChatChannelManage, auth.ChatClientId) {
			return common.ErrUnauthorized.Err()
		}
	} else if err != nil {
		return common.ErrUnauthorized.Err()
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
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChat, auth.ChatClientId) {
		log.Logger.WithContext(ctx).Infof("unauthorized request")
		return nil, common.ErrUnauthorized.Err()
	}
	character, err := s.verifyUserOwnsCharacter(
		ctx,
		&pb.CharacterTarget{Type: &pb.CharacterTarget_Name{Name: request.ChatMessage.CharacterName}},
	)
	if err != nil {
		log.Logger.WithContext(ctx).Infof("verify owns character failed: %v", err)
		return nil, err
	}

	// Validate requester has chat channel permissions
	// @TODO: Optimize by only checking specific channel (cache results)
	if !claims.HasResourceRole(RoleChatChannelManage, auth.ChatClientId) {
		serverCtx, err := s.serverContext(ctx)
		if err != nil {
			log.Logger.WithContext(ctx).Infof("creating server context: %v", err)
			return nil, common.ErrHandleRequest.Err()
		}

		channels, err := s.server.ChatService.AuthorizedChannelsForCharacter(serverCtx, uint(character.Id))
		if err != nil {
			log.Logger.WithContext(ctx).Infof("verify auth failed for %s on channel %d: %v", claims.ID, request.ChannelId, err)
			return nil, err
		}

		canSend := false
		for _, channel := range channels {
			if channel.ID == uint(request.ChannelId) {
				canSend = true
				break
			}
		}
		if !canSend {
			log.Logger.WithContext(ctx).Infof("%s attempted sending message to chat channel %d without permission", character.Name, request.ChannelId)
			return nil, common.ErrUnauthorized.Err()
		}
	}

	err = s.server.ChatService.SendChannelMessage(
		ctx,
		request.ChatMessage.CharacterName,
		request.ChatMessage.Message,
		uint(request.ChannelId),
	)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("send channel chat message: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to send message")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) SendDirectMessage(
	ctx context.Context,
	request *pb.SendDirectMessageRequest,
) (*emptypb.Empty, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChat, auth.ChatClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	if _, err := s.verifyUserOwnsCharacter(
		ctx,
		&pb.CharacterTarget{Type: &pb.CharacterTarget_Name{Name: request.ChatMessage.CharacterName}},
	); err != nil {
		return nil, err
	}

	srvCtx, err := s.serverContext(ctx)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("create server context: %v", err)
		return nil, common.ErrHandleRequest.Err()
	}
	targetCharacterName, err := character.GetCharacterNameFromTarget(srvCtx, s.server.CharacterService, request.Target)
	if err != nil {
		return nil, err
	}

	if err := s.server.ChatService.SendDirectMessage(
		ctx,
		request.ChatMessage.CharacterName,
		request.ChatMessage.Message,
		targetCharacterName,
	); err != nil {
		log.Logger.WithContext(ctx).Errorf("send direct chat message: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to send message")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) GetChannel(
	ctx context.Context,
	request *pb.ChatChannelTarget,
) (*pb.ChatChannel, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, auth.ChatClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	c, err := s.server.ChatService.GetChannel(ctx, uint(request.Id))
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("get chat channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to get chat channel")
	}

	if c == nil {
		return nil, common.ErrDoesNotExist.Err()
	}

	return c.ToPb(), nil
}

func (s chatServiceServer) CreateChannel(
	ctx context.Context,
	request *pb.CreateChannelMessage,
) (*emptypb.Empty, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, auth.ChatClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	_, err := s.server.ChatService.CreateChannel(
		ctx,
		&chat.ChatChannel{
			Name:      request.Name,
			Dimension: request.Dimension,
		},
	)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Error(codes.InvalidArgument, "name is already taken in this dimension")
		}

		log.Logger.WithContext(ctx).Errorf("create chat channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to create chat channel")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) DeleteChannel(
	ctx context.Context,
	request *pb.ChatChannelTarget,
) (*emptypb.Empty, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, auth.ChatClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	err := s.server.ChatService.DeleteChannel(ctx, &chat.ChatChannel{
		Model: gorm.Model{ID: uint(request.Id)},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrDoesNotExist.Err()
		}

		log.Logger.WithContext(ctx).Errorf("delete channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to delete channel")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) EditChannel(
	ctx context.Context,
	request *pb.UpdateChatChannelRequest,
) (*emptypb.Empty, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, auth.ChatClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	_, err := s.server.ChatService.UpdateChannel(ctx, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrDoesNotExist.Err()
		}

		log.Logger.WithContext(ctx).Errorf("edit channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to edit channel")
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) AllChatChannels(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.ChatChannels, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, auth.ChatClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	channels, err := s.server.ChatService.AllChannels(ctx)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("edit channel: %v", err)
		return nil, status.Error(codes.Internal, "unable to get channels")
	}

	return channels.ToPb(), nil
}

func (s chatServiceServer) GetAuthorizedChatChannels(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*pb.ChatChannels, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChat, auth.ChatClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	character, err := s.verifyUserOwnsCharacter(ctx, request)
	if err == common.ErrNotOwner {
		if !claims.HasResourceRole(RoleChatChannelManage, auth.ChatClientId) {
			return nil, status.Error(codes.PermissionDenied, common.ErrNotOwner.Error())
		}
	} else if err != nil || character == nil {
		log.Logger.WithContext(ctx).Infof("verify owns character failed: %v", err)
		return nil, err
	}

	channels, err := s.server.ChatService.AuthorizedChannelsForCharacter(ctx, uint(character.Id))
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("get authorized channels: %v", err)
		return nil, status.Error(codes.Internal, "unable to get channels")
	}

	return channels.ToPb(), nil
}

func (s chatServiceServer) UpdateUserChatChannelAuthorizations(
	ctx context.Context,
	request *pb.RequestChatChannelAuthChange,
) (*emptypb.Empty, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleChatChannelManage, auth.ChatClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	srvCtx, err := s.serverContext(ctx)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("create server context: %v", err)
		return nil, common.ErrHandleRequest.Err()
	}
	targetCharacterId, err := character.GetCharacterIdFromTarget(srvCtx, s.server.CharacterService, request.Character)
	if err != nil {
		return nil, err
	}

	err = s.server.ChatService.ChangeAuthorizationForCharacter(
		ctx,
		targetCharacterId,
		*helpers.ArrayOfUint64ToUint(&request.Ids),
		request.Add,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Error(codes.InvalidArgument, "existing chat authorization conflict")
		}

		log.Logger.WithContext(ctx).Errorf("get authorized channels: %v", err)
		return nil, status.Error(codes.Internal, "unable to change authorization")
	}

	return &emptypb.Empty{}, nil
}

func NewChatServiceServer(
	ctx context.Context,
	server *chatApp.ChatServerContext,
) (pb.ChatServiceServer, error) {
	err := createRoles(ctx,
		server.ServerContext,
		&ChatRoles,
	)
	if err != nil {
		return nil, err
	}

	return &chatServiceServer{
		server: server,
	}, nil
}

// @TODO: Cache the token
func (s chatServiceServer) serverContext(ctx context.Context) (context.Context, error) {
	token, err := s.server.KeycloakClient.LoginClient(
		ctx,
		s.server.GlobalConfig.Chat.Keycloak.ClientId,
		s.server.GlobalConfig.Chat.Keycloak.ClientSecret,
		s.server.GlobalConfig.Keycloak.Realm,
	)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("keycloak login failure: %v", err)
		return nil, err
	}

	return auth.AddOutgoingToken(
		ctx,
		token.AccessToken,
	), nil
}

func (s chatServiceServer) verifyUserOwnsCharacter(ctx context.Context, request *pb.CharacterTarget) (*pb.CharacterDetails, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	srvCtx, err := s.serverContext(ctx)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("create server context: %v", err)
		return nil, common.ErrHandleRequest.Err()
	}

	character, err := s.server.CharacterService.GetCharacter(srvCtx, request)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("chat character service get for user: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to verify character")
	}

	if character == nil {
		return nil, status.Errorf(codes.Internal, "character does not exist")
	}

	if character.Owner != claims.Subject {
		return character, common.ErrNotOwner
	}

	return character, nil
}

func (s chatServiceServer) checkUserChannelAuth(ctx context.Context, userId string, channelId uint) error {
	serverAuthCtx, err := s.serverContext(ctx)
	if err != nil {
		return common.ErrHandleRequest.Err()
	}

	characters, err := s.server.CharacterService.GetAllCharactersForUser(serverAuthCtx, &pb.UserTarget{Target: &pb.UserTarget_Id{Id: userId}})

	if err != nil {
		log.Logger.WithContext(ctx).Errorf("get characters: %v", err)
		return common.ErrHandleRequest.Err()
	}

	for _, character := range characters.Characters {
		channels, err := s.server.ChatService.AuthorizedChannelsForCharacter(serverAuthCtx, uint(character.Id))
		if err != nil {
			log.Logger.WithContext(ctx).Errorf("getting authorized channels: %v", err)
			return common.ErrHandleRequest.Err()
		}

		for _, channel := range channels {
			if channel.ID == channelId {
				return nil
			}
		}
	}

	return common.ErrUnauthorized.Err()
}

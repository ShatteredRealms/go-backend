package srv

import (
	context "context"
	"fmt"
	"reflect"

	"github.com/Nerzal/gocloak/v13"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	characters "github.com/ShatteredRealms/go-backend/cmd/characters/app"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type charactersServiceServer struct {
	pb.UnimplementedCharactersServiceServer
	server *characters.CharactersServerContext
}

var (
	CharacterRoles = make([]*gocloak.Role, 0)

	RoleAddCharacterPlayTime = registerCharacterRole(&gocloak.Role{
		Name:        gocloak.StringP("playtime"),
		Description: gocloak.StringP("Allows adding playtime to any character"),
	})

	RoleCharacterManagement = registerCharacterRole(&gocloak.Role{
		Name:        gocloak.StringP("manage"),
		Description: gocloak.StringP("Allows creating, reading and deleting of own characters"),
	})

	RoleCharacterManagementOther = registerCharacterRole(&gocloak.Role{
		Name:        gocloak.StringP("manage_other"),
		Description: gocloak.StringP("Allows creating, reading and deleting of any characters"),
	})
)

func registerCharacterRole(role *gocloak.Role) *gocloak.Role {
	CharacterRoles = append(CharacterRoles, role)
	return role
}

// AddCharacterPlayTime implements pb.CharactersServiceServer
func (s *charactersServiceServer) AddCharacterPlayTime(
	ctx context.Context,
	request *pb.AddPlayTimeRequest,
) (*pb.PlayTimeResponse, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("extract claims: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleAddCharacterPlayTime, model.CharactersClientId) {
		return nil, model.ErrUnauthorized
	}

	characterId, err := s.getCharacterTargetId(ctx, request.Character)
	if err != nil {
		return nil, err
	}

	// Add playtime
	time, err := s.server.Service.AddPlayTime(ctx, characterId, request.Time)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not update playtime")
	}

	return &pb.PlayTimeResponse{Time: time}, nil
}

// CreateCharacter implements pb.CharactersServiceServer
func (s *charactersServiceServer) CreateCharacter(
	ctx context.Context,
	request *pb.CreateCharacterRequest,
) (*pb.CharacterResponse, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("extract claims: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, model.CharactersClientId) {
		return nil, model.ErrUnauthorized
	}

	ownerId, err := s.getUserIdFromTarget(ctx, request.Owner)
	if err != nil {
		return nil, err
	}

	// If not requesting to create character for self, verify requester has permission for other
	if claims.Subject != ownerId && !claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, model.ErrUnauthorized
	}

	// Create new character
	char, err := s.server.Service.Create(ctx, ownerId, request.Name, request.Gender, request.Realm)
	if err != nil || char == nil {
		log.WithContext(ctx).Errorf("create char: %v", err)
		return nil, status.Error(codes.Internal, "unable to create character")
	}

	return &pb.CharacterResponse{
		Id:       uint64(char.ID),
		Owner:    char.OwnerId,
		Name:     char.Name,
		Gender:   char.Gender,
		Realm:    char.Realm,
		PlayTime: char.PlayTime,
		Location: nil,
	}, nil
}

// DeleteCharacter implements pb.CharactersServiceServer
func (s *charactersServiceServer) DeleteCharacter(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("extract claims: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, model.CharactersClientId) {
		return nil, model.ErrUnauthorized
	}

	character, err := s.getCharacterFromTarget(ctx, request)
	if err != nil {
		return nil, err
	}

	// If not requesting to delete requester own character, verify it has permission to delete others
	if claims.Subject != character.OwnerId && !claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, model.ErrUnauthorized
	}

	err = s.server.Service.Delete(ctx, character.ID)
	if err != nil {
		log.WithContext(ctx).Errorf("delete character %d: %v", character.ID, err)
		return nil, status.Error(codes.Internal, "unable to delete character")
	}

	return &emptypb.Empty{}, nil
}

// EditCharacter implements pb.CharactersServiceServer
func (s *charactersServiceServer) EditCharacter(
	ctx context.Context,
	request *pb.EditCharacterRequest,
) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("extract claims: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, model.CharactersClientId) {
		return nil, model.ErrUnauthorized
	}

	character, err := s.getCharacterFromTarget(ctx, request.Target)
	if err != nil {
		return nil, err
	}

	if character.OwnerId != claims.Subject && claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission to change playtime otherwise don't allow changing
	if !claims.HasResourceRole(RoleAddCharacterPlayTime, model.CharactersClientId) {
		request.PlayTime = nil
	}

	_, err = s.server.Service.Edit(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to edit user")
	}

	return &emptypb.Empty{}, nil
}

// GetCharacter implements pb.CharactersServiceServer
func (s *charactersServiceServer) GetCharacter(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*pb.CharacterResponse, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Debugf("ctx auth: %v", metautils.ExtractIncoming(ctx).Get("authorization"))
		log.WithContext(ctx).Errorf("extract claims: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, model.CharactersClientId) {
		log.WithContext(ctx).Error("no permission")
		return nil, model.ErrUnauthorized
	}

	character, err := s.getCharacterFromTarget(ctx, request)
	if err != nil {
		log.WithContext(ctx).Errorf("get character from target: %v", err)
		return nil, err
	}

	if character.OwnerId != claims.Subject && claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, model.ErrUnauthorized
	}

	return character.ToPb(), nil
}

// GetAllCharactersForUser implements pb.CharactersServiceServer
func (s *charactersServiceServer) GetAllCharactersForUser(
	ctx context.Context,
	request *pb.UserTarget,
) (*pb.CharactersResponse, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("extract claims: %v", err)
		return nil, model.ErrUnauthorized
	}

	id, err := s.getUserIdFromTarget(ctx, request)
	if err != nil {
		log.WithContext(ctx).Errorf("get user id from target: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if id != claims.Subject &&
		!claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, model.ErrUnauthorized
	}

	chars, err := s.server.Service.FindAllByOwner(ctx, id)
	if err != nil {
		log.WithContext(ctx).Errorf("find by owner %s: %v", claims.Subject, chars)
		return nil, status.Error(codes.Internal, "unable to find chars")
	}

	return chars.ToPb(), nil
}

// GetCharacters implements pb.CharactersServiceServer
func (s *charactersServiceServer) GetCharacters(
	ctx context.Context,
	msg *emptypb.Empty,
) (*pb.CharactersResponse, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("extract claims: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, model.ErrUnauthorized
	}

	chars, err := s.server.Service.FindAll(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("find all characters: %v", err)
		return nil, status.Error(codes.Internal, "unable to find chars")
	}

	return chars.ToPb(), nil
}

// GetGenders implements pb.CharactersServiceServer
func (s *charactersServiceServer) GetGenders(
	context.Context,
	*emptypb.Empty,
) (*pb.Genders, error) {
	return &pb.Genders{Genders: model.GetGenders()}, nil
}

// GetRealms implements pb.CharactersServiceServer
func (s *charactersServiceServer) GetRealms(
	context.Context,
	*emptypb.Empty,
) (*pb.Realms, error) {
	return &pb.Realms{Realms: model.GetRealms()}, nil
}

func NewCharactersServiceServer(
	ctx context.Context,
	server *characters.CharactersServerContext,
) (pb.CharactersServiceServer, error) {
	token, err := server.KeycloakClient.LoginClient(
		ctx,
		server.GlobalConfig.Characters.Keycloak.ClientId,
		server.GlobalConfig.Characters.Keycloak.ClientSecret,
		server.GlobalConfig.Characters.Keycloak.Realm,
	)
	if err != nil {
		return nil, fmt.Errorf("login keycloak: %v", err)
	}

	err = createRoles(ctx,
		server.KeycloakClient,
		token.AccessToken,
		server.GlobalConfig.Characters.Keycloak.Realm,
		server.GlobalConfig.Characters.Keycloak.Id,
		&CharacterRoles,
	)
	if err != nil {
		return nil, err
	}

	return &charactersServiceServer{
		server: server,
	}, nil
}

func (s charactersServiceServer) serverContext(ctx context.Context) (context.Context, error) {
	token, err := s.server.KeycloakClient.LoginClient(
		ctx,
		s.server.GlobalConfig.Chat.Keycloak.ClientId,
		s.server.GlobalConfig.Chat.Keycloak.ClientSecret,
		s.server.GlobalConfig.Chat.Keycloak.Realm,
	)
	if err != nil {
		return nil, err
	}

	return helpers.ContextAddClientToken(
		ctx,
		token.AccessToken,
	), nil
}

func (s charactersServiceServer) getCharacterTargetId(
	ctx context.Context,
	request *pb.CharacterTarget,
) (uint, error) {
	var characterId uint
	switch target := request.Target.(type) {
	case *pb.CharacterTarget_Id:
		characterId = uint(target.Id)
	case *pb.CharacterTarget_Name:
		char, err := s.server.Service.FindByName(ctx, target.Name)
		if err != nil {
			log.WithContext(ctx).Errorf("find character %s: %v", target.Name, err)
			return 0, status.Error(codes.Internal, "unable to find character")
		}
		if char == nil {
			return 0, status.Error(codes.InvalidArgument, "character does not exist")
		}

		characterId = char.ID
	default:
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return 0, model.ErrHandleRequest
	}

	return characterId, nil
}

func (s charactersServiceServer) getCharacterFromTarget(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*model.Character, error) {
	var character *model.Character
	var err error

	switch target := request.Target.(type) {
	case *pb.CharacterTarget_Id:
		character, err = s.server.Service.FindById(ctx, uint(target.Id))

	case *pb.CharacterTarget_Name:
		character, err = s.server.Service.FindByName(ctx, target.Name)

	default:
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		log.WithContext(ctx).Debugf("err: %v", err)
		return nil, model.ErrHandleRequest
	}

	if character == nil {
		log.WithContext(ctx).Debugf("character not found")
		return nil, model.ErrDoesNotExist
	}

	return character, nil
}

func (s charactersServiceServer) getUserIdFromTarget(
	ctx context.Context,
	request *pb.UserTarget,
) (string, error) {
	token, err := s.server.KeycloakClient.LoginClient(
		ctx,
		s.server.GlobalConfig.Characters.Keycloak.ClientId,
		s.server.GlobalConfig.Characters.Keycloak.ClientSecret,
		s.server.GlobalConfig.Characters.Keycloak.Realm,
	)
	if err != nil {
		log.WithContext(ctx).Errorf("login keycloak: %v", err)
		return "", model.ErrHandleRequest
	}

	ownerId := request.GetId()
	if val, ok := request.Target.(*pb.UserTarget_Username); ok {
		resp, err := s.server.KeycloakClient.GetUsers(
			ctx,
			token.AccessToken,
			s.server.GlobalConfig.Characters.Keycloak.Realm,
			gocloak.GetUsersParams{
				Exact:    gocloak.BoolP(true),
				Username: gocloak.StringP(val.Username),
			},
		)
		if err != nil {
			log.WithContext(ctx).Errorf("keycloak get users: %v", err)
			return "", model.ErrHandleRequest
		}
		if len(resp) == 0 || len(resp) > 1 {
			return "", model.ErrDoesNotExist
		}

		ownerId = *resp[0].ID
	}

	return ownerId, nil
}

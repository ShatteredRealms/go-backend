package srv

import (
	context "context"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
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
	msg *pb.CharacterTarget,
) (*pb.PlayTimeResponse, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleAddCharacterPlayTime, model.CharactersClientId) {
		return nil, ErrNotAuthorized
	}

	// Add playtime
	time, err := s.server.Service.AddPlayTime(ctx, msg.CharacterId, 1)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not update playtime")
	}

	return &pb.PlayTimeResponse{Time: time}, nil
}

// CreateCharacter implements pb.CharactersServiceServer
func (s *charactersServiceServer) CreateCharacter(
	ctx context.Context,
	msg *pb.CreateCharacterRequest,
) (*pb.CharacterResponse, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleCharacterManagement, model.CharactersClientId) {
		return nil, ErrNotAuthorized
	}

	// If not requesting to create character for self, verify requester has permission for other
	if claims.Subject != msg.UserId && !claims.HasRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, ErrNotAuthorized
	}

	// Create new character
	char, err := s.server.Service.Create(ctx, msg.UserId, msg.Name, msg.Gender, msg.Realm)
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
	msg *pb.CharacterTarget,
) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleCharacterManagement, model.CharactersClientId) {
		return nil, ErrNotAuthorized
	}

	character, err := s.server.Service.FindById(ctx, msg.CharacterId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "character %d does not exist", msg.CharacterId)
	}

	// If not requesting to delete requester own character, verify it has permission to delete others
	if claims.Subject != character.OwnerId && !claims.HasRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, ErrNotAuthorized
	}

	err = s.server.Service.Delete(ctx, msg.CharacterId)
	if err != nil {
		log.WithContext(ctx).Errorf("delete character %d: %v", msg.CharacterId, err)
		return nil, status.Error(codes.Internal, "unable to delete character")
	}

	return &emptypb.Empty{}, nil
}

// EditCharacter implements pb.CharactersServiceServer
func (s *charactersServiceServer) EditCharacter(
	ctx context.Context,
	msg *pb.EditCharacterRequest,
) (*emptypb.Empty, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleCharacterManagement, model.CharactersClientId) {
		return nil, ErrNotAuthorized
	}

	char, err := s.server.Service.FindById(ctx, msg.Id)
	if err != nil {
		log.WithContext(ctx).Errorf("find character %d: %v", msg.Id, err)
		return nil, status.Error(codes.Internal, "unable to find character")
	}

	if char == nil {
		return nil, status.Error(codes.InvalidArgument, "character does not exist")
	}

	if char.OwnerId != claims.Subject && claims.HasRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, ErrNotAuthorized
	}

	_, err = s.server.Service.Edit(ctx, msg)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to edit user")
	}

	return &emptypb.Empty{}, nil
}

// GetCharacter implements pb.CharactersServiceServer
func (s *charactersServiceServer) GetCharacter(
	ctx context.Context,
	msg *pb.CharacterTarget,
) (*pb.CharacterResponse, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	// Validate requester has correct permission
	if !claims.HasRole(RoleCharacterManagement, model.CharactersClientId) {
		return nil, ErrNotAuthorized
	}

	char, err := s.server.Service.FindById(ctx, msg.CharacterId)
	if err != nil {
		log.WithContext(ctx).Errorf("find character %d: %v", msg.CharacterId, err)
		return nil, status.Error(codes.Internal, "unable to find character")
	}

	if char == nil {
		return nil, status.Error(codes.InvalidArgument, "character does not exist")
	}

	if char.OwnerId != claims.Subject && claims.HasRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, ErrNotAuthorized
	}

	return char.ToPb(), nil
}

// GetCharacters implements pb.CharactersServiceServer
func (s *charactersServiceServer) GetCharacters(
	ctx context.Context,
	msg *emptypb.Empty,
) (*pb.CharactersResponse, error) {
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		return nil, ErrNotAuthorized
	}

	chars, err := s.server.Service.FindAllByOwner(ctx, claims.Subject)
	if err != nil {
		log.WithContext(ctx).Errorf("find by owner %s: %v", claims.Subject, chars)
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

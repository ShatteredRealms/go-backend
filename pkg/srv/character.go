package srv

import (
	context "context"
	"fmt"
	"reflect"

	"github.com/Nerzal/gocloak/v13"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	characters "github.com/ShatteredRealms/go-backend/cmd/character/app"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type charactersServiceServer struct {
	pb.UnimplementedCharacterServiceServer
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
		Description: gocloak.StringP("Allows creating, reading, editing and deleting of any characters"),
	})

	RoleInventoryManagement = registerCharacterRole(&gocloak.Role{
		Name:        gocloak.StringP("inventory_manage"),
		Description: gocloak.StringP("Allows getting and updating character inventories"),
	})
)

func registerCharacterRole(role *gocloak.Role) *gocloak.Role {
	CharacterRoles = append(CharacterRoles, role)
	return role
}

// AddCharacterPlayTime implements pb.CharacterServiceServer
func (s *charactersServiceServer) AddCharacterPlayTime(
	ctx context.Context,
	request *pb.AddPlayTimeRequest,
) (*pb.PlayTimeResponse, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleAddCharacterPlayTime, model.CharactersClientId) {
		return nil, model.ErrUnauthorized.Err()
	}

	characterId, err := s.getCharacterTargetId(ctx, request.Character)
	if err != nil {
		return nil, err
	}

	// Add playtime
	character, err := s.server.CharacterService.AddPlayTime(ctx, characterId, request.Time)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not update playtime")
	}

	return &pb.PlayTimeResponse{Time: character.PlayTime}, nil
}

// CreateCharacter implements pb.CharacterServiceServer
func (s *charactersServiceServer) CreateCharacter(
	ctx context.Context,
	request *pb.CreateCharacterRequest,
) (*pb.CharacterDetails, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, model.CharactersClientId) {
		return nil, model.ErrUnauthorized.Err()
	}

	ownerId, err := s.getUserIdFromTarget(ctx, request.Owner)
	if err != nil {
		return nil, err
	}

	// If not requesting to create character for self, verify requester has permission for other
	if claims.Subject != ownerId && !claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, model.ErrUnauthorized.Err()
	}

	// Create new character
	char, err := s.server.CharacterService.Create(ctx, ownerId, request.Name, request.Gender, request.Realm)
	if err != nil || char == nil {
		log.Logger.WithContext(ctx).Errorf("create char: %v", err)
		return nil, status.Error(codes.Internal, "unable to create character")
	}

	return &pb.CharacterDetails{
		Id:       uint64(char.ID),
		Owner:    char.OwnerId,
		Name:     char.Name,
		Gender:   char.Gender,
		Realm:    char.Realm,
		PlayTime: char.PlayTime,
		Location: nil,
	}, nil
}

// DeleteCharacter implements pb.CharacterServiceServer
func (s *charactersServiceServer) DeleteCharacter(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*emptypb.Empty, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, model.CharactersClientId) {
		return nil, model.ErrUnauthorized.Err()
	}

	character, err := s.getCharacterFromTarget(ctx, request)
	if err != nil {
		return nil, err
	}

	if character == nil {
		return nil, model.ErrDoesNotExist.Err()
	}

	// If not requesting to delete requester own character, verify it has permission to delete others
	if claims.Subject != character.OwnerId && !claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, model.ErrUnauthorized.Err()
	}

	err = s.server.CharacterService.Delete(ctx, character.ID)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("delete character %d: %v", character.ID, err)
		return nil, status.Error(codes.Internal, "unable to delete character")
	}

	return &emptypb.Empty{}, nil
}

// EditCharacter implements pb.CharacterServiceServer
func (s *charactersServiceServer) EditCharacter(
	ctx context.Context,
	request *pb.EditCharacterRequest,
) (*emptypb.Empty, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission to change playtime otherwise don't allow changing
	if !claims.HasResourceRole(RoleAddCharacterPlayTime, model.CharactersClientId) {
		request.OptionalPlayTime = nil
	}

	_, err = s.server.CharacterService.Edit(ctx, request)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("edit character: %v", err)
		return nil, status.Error(codes.Internal, "unable to character user")
	}

	return &emptypb.Empty{}, nil
}

// GetCharacter implements pb.CharacterServiceServer
func (s *charactersServiceServer) GetCharacter(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*pb.CharacterDetails, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, model.CharactersClientId) {
		log.Logger.WithContext(ctx).Error("no permission")
		return nil, model.ErrUnauthorized.Err()
	}

	character, err := s.getCharacterFromTarget(ctx, request)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("get character from target: %v", err)
		return nil, err
	}

	if character == nil {
		return nil, model.ErrDoesNotExist.Err()
	}

	if character.OwnerId != claims.Subject && !claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		log.Logger.WithContext(ctx).Infof("user %s requested character %s without %s", claims.Subject, character.Name, *RoleCharacterManagementOther.Name)
		return nil, model.ErrUnauthorized.Err()
	}

	return character.ToPb(), nil
}

// GetAllCharactersForUser implements pb.CharacterServiceServer
func (s *charactersServiceServer) GetAllCharactersForUser(
	ctx context.Context,
	request *pb.UserTarget,
) (*pb.CharactersDetails, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, model.CharactersClientId) {
		return nil, model.ErrUnauthorized.Err()
	}

	id, err := s.getUserIdFromTarget(ctx, request)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("get user id from target: %v", err)
		return nil, model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if id != claims.Subject &&
		!claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, model.ErrUnauthorized.Err()
	}

	chars, err := s.server.CharacterService.FindAllByOwner(ctx, id)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("find by owner %s: %v", claims.Subject, chars)
		return nil, status.Error(codes.Internal, "unable to find chars")
	}

	return chars.ToPb(), nil
}

// GetCharacters implements pb.CharacterServiceServer
func (s *charactersServiceServer) GetCharacters(
	ctx context.Context,
	msg *emptypb.Empty,
) (*pb.CharactersDetails, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagementOther, model.CharactersClientId) {
		return nil, model.ErrUnauthorized.Err()
	}

	chars, err := s.server.CharacterService.FindAll(ctx)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("find all characters: %v", err)
		return nil, status.Error(codes.Internal, "unable to find chars")
	}

	return chars.ToPb(), nil
}

// GetInventory implements pb.CharacterServiceServer.
func (s *charactersServiceServer) GetInventory(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*pb.Inventory, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleInventoryManagement, model.CharactersClientId) {
		return nil, model.ErrUnauthorized.Err()
	}

	character, err := s.getCharacterFromTarget(ctx, request)
	if err != nil {
		return nil, err
	}

	if character == nil {
		return nil, model.ErrDoesNotExist.Err()
	}

	inv, err := s.server.InventoryService.GetInventory(ctx, character.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.Inventory{}, nil
		}

		return nil, status.Errorf(codes.Internal, "get inventory for %s: %s", character.Name, err.Error())
	}

	return inv.ToPb(), nil
}

// SetInventory implements pb.CharacterServiceServer.
func (s *charactersServiceServer) SetInventory(
	ctx context.Context,
	request *pb.UpdateInventoryRequest,
) (*emptypb.Empty, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleInventoryManagement, model.CharactersClientId) {
		return nil, model.ErrUnauthorized.Err()
	}

	character, err := s.getCharacterFromTarget(ctx, request.Target)
	if err != nil {
		return nil, err
	}

	if character == nil {
		return nil, model.ErrDoesNotExist.Err()
	}

	newInv := &model.CharacterInventory{
		CharacterId: character.ID,
		Inventory:   model.InventoryItemsFromPb(request.InventoryItems),
		Bank:        model.InventoryItemsFromPb(request.BankItems),
	}
	err = s.server.InventoryService.UpdateInventory(ctx, newInv)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &emptypb.Empty{}, nil
		}

		return nil, status.Errorf(codes.Internal, "set inventory for %s: %s", character.Name, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func NewCharacterServiceServer(
	ctx context.Context,
	server *characters.CharactersServerContext,
) (pb.CharacterServiceServer, error) {
	token, err := server.KeycloakClient.LoginClient(
		ctx,
		server.GlobalConfig.Character.Keycloak.ClientId,
		server.GlobalConfig.Character.Keycloak.ClientSecret,
		server.GlobalConfig.Keycloak.Realm,
	)
	if err != nil {
		return nil, fmt.Errorf("login keycloak: %v", err)
	}

	err = createRoles(ctx,
		server.KeycloakClient,
		token.AccessToken,
		server.GlobalConfig.Keycloak.Realm,
		server.GlobalConfig.Character.Keycloak.Id,
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
		s.server.GlobalConfig.Character.Keycloak.ClientId,
		s.server.GlobalConfig.Character.Keycloak.ClientSecret,
		s.server.GlobalConfig.Keycloak.Realm,
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
	switch target := request.Type.(type) {
	case *pb.CharacterTarget_Id:
		characterId = uint(target.Id)
	case *pb.CharacterTarget_Name:
		char, err := s.server.CharacterService.FindByName(ctx, target.Name)
		if err != nil {
			log.Logger.WithContext(ctx).Errorf("find character %s: %v", target.Name, err)
			return 0, status.Error(codes.Internal, "unable to find character")
		}
		if char == nil {
			return 0, status.Error(codes.InvalidArgument, "character does not exist")
		}

		characterId = char.ID
	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return 0, model.ErrHandleRequest.Err()
	}

	return characterId, nil
}

func (s charactersServiceServer) getCharacterFromTarget(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*model.Character, error) {
	var character *model.Character
	var err error

	if request == nil || request.Type == nil {
		return nil, status.Errorf(codes.InvalidArgument, "target is empty")
	}

	switch target := request.Type.(type) {
	case *pb.CharacterTarget_Id:
		character, err = s.server.CharacterService.FindById(ctx, uint(target.Id))

	case *pb.CharacterTarget_Name:
		character, err = s.server.CharacterService.FindByName(ctx, target.Name)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest.Err()
	}

	if err != nil {
		log.Logger.WithContext(ctx).Debugf("err: %v", err)
		return nil, model.ErrHandleRequest.Err()
	}

	if character == nil {
		log.Logger.WithContext(ctx).Debugf("character not found")
		return nil, model.ErrDoesNotExist.Err()
	}

	return character, nil
}

func (s charactersServiceServer) getUserIdFromTarget(
	ctx context.Context,
	request *pb.UserTarget,
) (string, error) {
	token, err := s.server.KeycloakClient.LoginClient(
		ctx,
		s.server.GlobalConfig.Character.Keycloak.ClientId,
		s.server.GlobalConfig.Character.Keycloak.ClientSecret,
		s.server.GlobalConfig.Keycloak.Realm,
	)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("login keycloak: %v", err)
		return "", model.ErrHandleRequest.Err()
	}

	ownerId := request.GetId()
	if val, ok := request.Target.(*pb.UserTarget_Username); ok {
		resp, err := s.server.KeycloakClient.GetUsers(
			ctx,
			token.AccessToken,
			s.server.GlobalConfig.Keycloak.Realm,
			gocloak.GetUsersParams{
				Exact:    gocloak.BoolP(true),
				Username: gocloak.StringP(val.Username),
			},
		)
		if err != nil {
			log.Logger.WithContext(ctx).Errorf("keycloak get users: %v", err)
			return "", model.ErrHandleRequest.Err()
		}
		if len(resp) == 0 || len(resp) > 1 {
			return "", model.ErrDoesNotExist.Err()
		}

		ownerId = *resp[0].ID
	}

	return ownerId, nil
}

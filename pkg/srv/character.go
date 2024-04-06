package srv

import (
	context "context"
	"errors"
	"reflect"

	"github.com/ShatteredRealms/go-backend/pkg/auth"
	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/WilSimpson/gocloak/v13"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	characterApp "github.com/ShatteredRealms/go-backend/cmd/character/app"
	"github.com/ShatteredRealms/go-backend/pkg/model/character"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type charactersServiceServer struct {
	pb.UnimplementedCharacterServiceServer
	server *characterApp.CharacterServerContext
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
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleAddCharacterPlayTime, auth.CharacterClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	characterId, err := s.getCharacterTargetId(ctx, request.Character)
	if err != nil {
		return nil, err
	}

	// Add playtime
	chara, err := s.server.CharacterService.AddPlayTime(ctx, characterId, request.Time)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not update playtime")
	}

	return &pb.PlayTimeResponse{Time: chara.PlayTime}, nil
}

// CreateCharacter implements pb.CharacterServiceServer
func (s *charactersServiceServer) CreateCharacter(
	ctx context.Context,
	request *pb.CreateCharacterRequest,
) (*pb.CharacterDetails, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, auth.CharacterClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	ownerId, err := s.getUserIdFromTarget(ctx, request.Owner)
	if err != nil {
		return nil, err
	}

	// If not requesting to create character for self, verify requester has permission for other
	if claims.Subject != ownerId && !claims.HasResourceRole(RoleCharacterManagementOther, auth.CharacterClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	srvManagerClient, err := s.server.GetServerManagerClient()
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("get server manager client: %v", err)
		return nil, ErrInternalCreateCharacter
	}

	authCtx, err := s.server.OutgoingClientAuth(ctx)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("outgoing client auth: %v", err)
		return nil, ErrInternalCreateCharacter
	}

	_, err = srvManagerClient.GetDimension(authCtx, &pb.DimensionTarget{FindBy: &pb.DimensionTarget_Name{Name: request.Dimension}})
	if err != nil {
		if errors.Is(err, common.ErrDoesNotExist.Err()) {
			log.Logger.WithContext(ctx).Errorf("invalid dimension requested: %v", err)
			return nil, ErrInvalidDimension
		}
		log.Logger.WithContext(ctx).Errorf("get dimension: %v", err)
		return nil, ErrInternalCreateCharacter
	}

	// Create new character
	char, err := s.server.CharacterService.Create(ctx, ownerId, request.Name, request.Gender, request.Realm, request.Dimension)
	if err != nil || char == nil {
		log.Logger.WithContext(ctx).Errorf("create char: %v", err)
		return nil, ErrInternalCreateCharacter
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
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, auth.CharacterClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	char, err := s.getCharacterFromTarget(ctx, request)
	if err != nil {
		return nil, err
	}

	if char == nil {
		return nil, common.ErrDoesNotExist.Err()
	}

	// If not requesting to delete requester own character, verify it has permission to delete others
	if claims.Subject != char.OwnerId && !claims.HasResourceRole(RoleCharacterManagementOther, auth.CharacterClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	err = s.server.CharacterService.Delete(ctx, char.ID)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("delete character %d: %v", char.ID, err)
		return nil, status.Error(codes.Internal, "unable to delete character")
	}

	return &emptypb.Empty{}, nil
}

// EditCharacter implements pb.CharacterServiceServer
func (s *charactersServiceServer) EditCharacter(
	ctx context.Context,
	request *pb.EditCharacterRequest,
) (*emptypb.Empty, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagementOther, auth.CharacterClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission to change playtime otherwise don't allow changing
	if !claims.HasResourceRole(RoleAddCharacterPlayTime, auth.CharacterClientId) {
		request.OptionalPlayTime = nil
	}

	_, err := s.server.CharacterService.Edit(ctx, request)
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
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, auth.CharacterClientId) {
		log.Logger.WithContext(ctx).Error("no permission")
		return nil, common.ErrUnauthorized.Err()
	}

	char, err := s.getCharacterFromTarget(ctx, request)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("get character from target: %v", err)
		return nil, err
	}

	if char == nil {
		return nil, common.ErrDoesNotExist.Err()
	}

	if char.OwnerId != claims.Subject && !claims.HasResourceRole(RoleCharacterManagementOther, auth.CharacterClientId) {
		log.Logger.WithContext(ctx).Infof("user %s requested character %s without %s", claims.Subject, char.Name, *RoleCharacterManagementOther.Name)
		return nil, common.ErrUnauthorized.Err()
	}

	return char.ToPb(), nil
}

// GetAllCharactersForUser implements pb.CharacterServiceServer
func (s *charactersServiceServer) GetAllCharactersForUser(
	ctx context.Context,
	request *pb.UserTarget,
) (*pb.CharactersDetails, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagement, auth.CharacterClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	id, err := s.getUserIdFromTarget(ctx, request)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("get user id from target: %v", err)
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if id != claims.Subject &&
		!claims.HasResourceRole(RoleCharacterManagementOther, auth.CharacterClientId) {
		return nil, common.ErrUnauthorized.Err()
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
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleCharacterManagementOther, auth.CharacterClientId) {
		return nil, common.ErrUnauthorized.Err()
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
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleInventoryManagement, auth.CharacterClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	char, err := s.getCharacterFromTarget(ctx, request)
	if err != nil {
		return nil, err
	}

	if char == nil {
		return nil, common.ErrDoesNotExist.Err()
	}

	inv, err := s.server.InventoryService.GetInventory(ctx, char.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.Inventory{}, nil
		}

		return nil, status.Errorf(codes.Internal, "get inventory for %s: %s", char.Name, err.Error())
	}

	return inv.ToPb(), nil
}

// SetInventory implements pb.CharacterServiceServer.
func (s *charactersServiceServer) SetInventory(
	ctx context.Context,
	request *pb.UpdateInventoryRequest,
) (*emptypb.Empty, error) {
	claims, ok := auth.RetrieveClaims(ctx)
	if !ok {
		return nil, common.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleInventoryManagement, auth.CharacterClientId) {
		return nil, common.ErrUnauthorized.Err()
	}

	char, err := s.getCharacterFromTarget(ctx, request.Target)
	if err != nil {
		return nil, err
	}

	if char == nil {
		return nil, common.ErrDoesNotExist.Err()
	}

	newInv := &character.Inventory{
		CharacterId: char.ID,
		Inventory:   character.InventoryItemsFromPb(request.InventoryItems),
		Bank:        character.InventoryItemsFromPb(request.BankItems),
	}
	err = s.server.InventoryService.UpdateInventory(ctx, newInv)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &emptypb.Empty{}, nil
		}

		return nil, status.Errorf(codes.Internal, "set inventory for %s: %s", char.Name, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func NewCharacterServiceServer(
	ctx context.Context,
	server *characterApp.CharacterServerContext,
) (pb.CharacterServiceServer, error) {
	err := createRoles(ctx,
		server.ServerContext,
		&CharacterRoles,
	)
	if err != nil {
		return nil, err
	}

	return &charactersServiceServer{
		server: server,
	}, nil
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
		return 0, common.ErrHandleRequest.Err()
	}

	return characterId, nil
}

func (s charactersServiceServer) getCharacterFromTarget(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*character.Character, error) {
	var char *character.Character
	var err error

	if request == nil || request.Type == nil {
		return nil, status.Errorf(codes.InvalidArgument, "target is empty")
	}

	switch target := request.Type.(type) {
	case *pb.CharacterTarget_Id:
		char, err = s.server.CharacterService.FindById(ctx, uint(target.Id))

	case *pb.CharacterTarget_Name:
		char, err = s.server.CharacterService.FindByName(ctx, target.Name)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, common.ErrHandleRequest.Err()
	}

	if err != nil {
		log.Logger.WithContext(ctx).Debugf("err: %v", err)
		return nil, common.ErrHandleRequest.Err()
	}

	if char == nil {
		log.Logger.WithContext(ctx).Debugf("character not found")
		return nil, common.ErrDoesNotExist.Err()
	}

	return char, nil
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
		return "", common.ErrHandleRequest.Err()
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
			return "", common.ErrHandleRequest.Err()
		}
		if len(resp) == 0 || len(resp) > 1 {
			return "", common.ErrDoesNotExist.Err()
		}

		ownerId = *resp[0].ID
	}

	return ownerId, nil
}

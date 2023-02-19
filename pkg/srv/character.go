package srv

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	utilService "github.com/ShatteredRealms/go-backend/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type characterServiceServer struct {
	pb.UnimplementedCharactersServiceServer
	characterService service.CharacterService
	jwtService       utilService.JWTService
}

func NewCharacterServiceServer(
	characterService service.CharacterService,
	jwtService utilService.JWTService,
) pb.CharactersServiceServer {
	return &characterServiceServer{
		characterService: characterService,
		jwtService:       jwtService,
	}
}

func (s *characterServiceServer) GetAllGenders(context.Context, *emptypb.Empty) (*pb.Genders, error) {
	return model.Genders, nil
}

func (s *characterServiceServer) GetAllRealms(context.Context, *emptypb.Empty) (*pb.Realms, error) {
	return model.Realms, nil
}

func (s *characterServiceServer) GetAllCharacters(ctx context.Context, message *emptypb.Empty) (*pb.Characters, error) {
	characters, err := s.characterService.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return modelCharactersToPb(characters), nil
}

func (s *characterServiceServer) GetAllCharactersForUser(ctx context.Context, message *pb.UserTarget) (*pb.Characters, error) {
	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, message.Username)
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	characters, err := s.characterService.FindAllByOwner(ctx, message.Username)
	if err != nil {
		return nil, err
	}
	return modelCharactersToPb(characters), nil
}

func (s *characterServiceServer) GetCharacter(ctx context.Context, message *pb.CharacterTarget) (*pb.Character, error) {
	character, err := s.characterService.FindById(ctx, message.CharacterId)
	if err != nil {
		return nil, err
	}

	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, character.Owner)
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	return modelCharacterToPb(character), nil
}

func (s *characterServiceServer) CreateCharacter(ctx context.Context, message *pb.CreateCharacterRequest) (*pb.Character, error) {
	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, message.Username)
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	character, err := s.characterService.Create(ctx, message.Username, message.Name, message.Gender, message.Realm)
	if err != nil {
		return nil, err
	}

	return modelCharacterToPb(character), nil
}

func (s *characterServiceServer) DeleteCharacter(ctx context.Context, message *pb.Character) (*emptypb.Empty, error) {
	character, err := s.characterService.FindById(ctx, message.Id)
	if err != nil {
		return nil, err
	}

	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, character.Owner)
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	return &emptypb.Empty{}, s.characterService.Delete(ctx, message.Id)
}

func (s *characterServiceServer) EditCharacter(ctx context.Context, message *pb.Character) (*pb.Character, error) {
	character, err := s.characterService.FindById(ctx, message.Id)
	if err != nil {
		return nil, err
	}

	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, character.Owner)
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	character, err = s.characterService.Edit(ctx, message)
	if err != nil {
		return nil, err
	}

	return modelCharacterToPb(character), nil
}

func (s *characterServiceServer) AddCharacterPlayTime(ctx context.Context, message *pb.PlayTimeMessage) (*pb.PlayTimeMessage, error) {
	character, err := s.characterService.FindById(ctx, message.CharacterId)
	if err != nil {
		return nil, err
	}

	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, character.Owner)
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	newTime, err := s.characterService.AddPlayTime(ctx, message.CharacterId, message.Time)
	if err != nil {
		return nil, err
	}

	return &pb.PlayTimeMessage{Time: newTime}, nil
}

func modelCharacterToPb(character *model.Character) *pb.Character {
	return &pb.Character{
		Id:       uint64(character.ID),
		Owner:    wrapperspb.String(character.Owner),
		Name:     wrapperspb.String(character.Name),
		Gender:   wrapperspb.UInt64(character.GenderId),
		Realm:    wrapperspb.UInt64(character.RealmId),
		PlayTime: wrapperspb.UInt64(character.PlayTime),
		Location: &pb.Location{
			World: character.Location.World,
			X:     character.Location.X,
			Y:     character.Location.Y,
			Z:     character.Location.Z,
		},
	}
}

func modelCharactersToPb(characters []*model.Character) *pb.Characters {
	out := make([]*pb.Character, len(characters))

	for i, character := range characters {
		out[i] = modelCharacterToPb(character)
	}

	return &pb.Characters{Characters: out}
}

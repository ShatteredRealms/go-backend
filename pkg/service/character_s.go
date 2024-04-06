package service

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model/character"
	"github.com/ShatteredRealms/go-backend/pkg/model/game"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
)

type CharacterService interface {
	Create(ctx context.Context, ownerId string, name string, gender string, realm string, dimension string) (*character.Character, error)
	Edit(ctx context.Context, char *pb.EditCharacterRequest) (*character.Character, error)
	Delete(ctx context.Context, id uint) error

	FindById(ctx context.Context, id uint) (*character.Character, error)
	FindByName(ctx context.Context, name string) (*character.Character, error)

	FindAllByOwner(ctx context.Context, ownerId string) (character.Characters, error)

	FindAll(context.Context) (character.Characters, error)

	AddPlayTime(ctx context.Context, characterId uint, amount uint64) (*character.Character, error)
}

type characterService struct {
	repo repository.CharacterRepository
}

func (s characterService) FindByName(ctx context.Context, name string) (*character.Character, error) {
	return s.repo.FindByName(ctx, name)
}

func NewCharacterService(
	ctx context.Context,
	r repository.CharacterRepository,
) (CharacterService, error) {

	err := r.Migrate(ctx)
	if err != nil {
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	return characterService{
		repo: r,
	}, nil
}

func (s characterService) Create(ctx context.Context, ownerId string, name string, gender string, realm string, dimension string) (*character.Character, error) {
	character := character.Character{
		OwnerId:   ownerId,
		Name:      name,
		Gender:    gender,
		Realm:     realm,
		Dimension: dimension,
		PlayTime:  0,
	}

	if err := character.Validate(); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, &character)
}

func (s characterService) FindById(ctx context.Context, id uint) (*character.Character, error) {
	return s.repo.FindById(ctx, id)
}

func (s characterService) Edit(ctx context.Context, char *pb.EditCharacterRequest) (*character.Character, error) {
	var currentCharacter *character.Character
	var err error

	switch target := char.Target.Type.(type) {
	case *pb.CharacterTarget_Id:
		currentCharacter, err = s.FindById(ctx, uint(target.Id))
	case *pb.CharacterTarget_Name:
		currentCharacter, err = s.FindByName(ctx, target.Name)
	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return nil, common.ErrHandleRequest.Err()

	}

	if err != nil {
		return nil, err
	}

	if char.OptionalNewName != nil &&
		char.GetNewName() != "" {
		currentCharacter.Name = char.GetNewName()
	}
	if char.OptionalOwnerId != nil &&
		char.GetOwnerId() != "" {
		currentCharacter.OwnerId = char.GetOwnerId()
	}

	if char.OptionalPlayTime != nil {
		currentCharacter.PlayTime = char.GetPlayTime()
	}

	if char.OptionalGender != nil &&
		char.GetGender() != "" {
		currentCharacter.Gender = char.GetGender()
	}

	if char.OptionalRealm != nil &&
		char.GetRealm() != "" {
		currentCharacter.Realm = char.GetRealm()
	}

	if char.OptionalLocation != nil &&
		char.GetLocation().World != "" {
		currentCharacter.Location = *game.LocationFromPb(char.GetLocation())
	}

	err = currentCharacter.Validate()
	if err != nil {
		return nil, err
	}

	return s.repo.Save(ctx, currentCharacter)
}

func (s characterService) Delete(ctx context.Context, id uint) error {
	character, err := s.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("find by id: %v", err)
	}

	if character == nil {
		return fmt.Errorf("character id %d not found", id)
	}

	return s.repo.Delete(ctx, character)
}

func (s characterService) FindAll(ctx context.Context) (character.Characters, error) {
	return s.repo.FindAll(ctx)
}

func (s characterService) FindAllByOwner(ctx context.Context, owner string) (character.Characters, error) {
	return s.repo.FindAllByOwner(ctx, owner)
}

func (s characterService) AddPlayTime(ctx context.Context, characterId uint, amount uint64) (*character.Character, error) {
	character, err := s.FindById(ctx, characterId)
	if err != nil {
		return nil, err
	}

	character.PlayTime += amount
	character, err = s.repo.Save(ctx, character)

	return character, err
}

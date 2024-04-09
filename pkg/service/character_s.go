package service

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model/character"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
)

type CharacterService interface {
	Create(ctx context.Context, ownerId string, name string, gender string, realm string, dimension string) (*character.Character, error)
	Save(ctx context.Context, char *character.Character) (*character.Character, error)
	Delete(ctx context.Context, id uint) error

	FindById(ctx context.Context, id uint) (*character.Character, error)
	FindByName(ctx context.Context, name string) (*character.Character, error)
	FindByTarget(ctx context.Context, target *pb.CharacterTarget) (*character.Character, error)

	FindAllByOwner(ctx context.Context, ownerId string) (character.Characters, error)

	FindAll(context.Context) (character.Characters, error)

	AddPlayTime(ctx context.Context, characterId uint, amount uint64) (*character.Character, error)
}

type characterService struct {
	repo repository.CharacterRepository
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

// Save implements CharacterService.
func (s characterService) Save(ctx context.Context, char *character.Character) (*character.Character, error) {
	err := char.Validate()
	if err != nil {
		return nil, err
	}

	return s.repo.Save(ctx, char)
}

// FindByTarget implements CharacterService.
func (s characterService) FindByTarget(ctx context.Context, target *pb.CharacterTarget) (*character.Character, error) {
	var char *character.Character
	var err error

	switch target := target.Type.(type) {
	case *pb.CharacterTarget_Id:
		char, err = s.FindById(ctx, uint(target.Id))
	case *pb.CharacterTarget_Name:
		char, err = s.FindByName(ctx, target.Name)
	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return nil, common.ErrHandleRequest.Err()

	}

	if err != nil {
		return nil, err
	}

	return char, nil
}

func (s characterService) FindByName(ctx context.Context, name string) (*character.Character, error) {
	return s.repo.FindByName(ctx, name)
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

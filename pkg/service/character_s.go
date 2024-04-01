package service

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
)

type CharacterService interface {
	Create(ctx context.Context, ownerId string, name string, gender string, realm string) (*model.Character, error)
	Edit(ctx context.Context, character *pb.EditCharacterRequest) (*model.Character, error)
	Delete(ctx context.Context, id uint) error

	FindById(ctx context.Context, id uint) (*model.Character, error)
	FindByName(ctx context.Context, name string) (*model.Character, error)

	FindAllByOwner(ctx context.Context, ownerId string) (model.Characters, error)

	FindAll(context.Context) (model.Characters, error)

	AddPlayTime(ctx context.Context, characterId uint, amount uint64) (*model.Character, error)
}

type characterService struct {
	repo repository.CharacterRepository
}

func (s characterService) FindByName(ctx context.Context, name string) (*model.Character, error) {
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

func (s characterService) Create(ctx context.Context, ownerId string, name string, gender string, realm string) (*model.Character, error) {
	character := model.Character{
		OwnerId:  ownerId,
		Name:     name,
		Gender:   gender,
		Realm:    realm,
		PlayTime: 0,
	}

	if err := character.Validate(); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, &character)
}

func (s characterService) FindById(ctx context.Context, id uint) (*model.Character, error) {
	return s.repo.FindById(ctx, id)
}

func (s characterService) Edit(ctx context.Context, character *pb.EditCharacterRequest) (*model.Character, error) {
	var currentCharacter *model.Character
	var err error

	switch target := character.Target.Type.(type) {
	case *pb.CharacterTarget_Id:
		currentCharacter, err = s.FindById(ctx, uint(target.Id))
	case *pb.CharacterTarget_Name:
		currentCharacter, err = s.FindByName(ctx, target.Name)
	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return nil, model.ErrHandleRequest.Err()

	}

	if err != nil {
		return nil, err
	}

	if character.OptionalNewName != nil &&
		character.GetNewName() != "" {
		currentCharacter.Name = character.GetNewName()
	}
	if character.OptionalOwnerId != nil &&
		character.GetOwnerId() != "" {
		currentCharacter.OwnerId = character.GetOwnerId()
	}

	if character.OptionalPlayTime != nil {
		currentCharacter.PlayTime = character.GetPlayTime()
	}

	if character.OptionalGender != nil &&
		character.GetGender() != "" {
		currentCharacter.Gender = character.GetGender()
	}

	if character.OptionalRealm != nil &&
		character.GetRealm() != "" {
		currentCharacter.Realm = character.GetRealm()
	}

	if character.OptionalLocation != nil &&
		character.GetLocation().World != "" {
		currentCharacter.Location = *model.LocationFromPb(character.GetLocation())
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

func (s characterService) FindAll(ctx context.Context) (model.Characters, error) {
	return s.repo.FindAll(ctx)
}

func (s characterService) FindAllByOwner(ctx context.Context, owner string) (model.Characters, error) {
	return s.repo.FindAllByOwner(ctx, owner)
}

func (s characterService) AddPlayTime(ctx context.Context, characterId uint, amount uint64) (*model.Character, error) {
	character, err := s.FindById(ctx, characterId)
	if err != nil {
		return nil, err
	}

	character.PlayTime += amount
	character, err = s.repo.Save(ctx, character)

	return character, err
}

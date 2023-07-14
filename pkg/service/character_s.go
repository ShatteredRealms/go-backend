package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	log "github.com/sirupsen/logrus"
)

type CharacterService interface {
	Create(ctx context.Context, ownerId string, name string, gender string, realm string) (*model.Character, error)
	Edit(ctx context.Context, character *pb.EditCharacterRequest) (*model.Character, error)
	Delete(ctx context.Context, id uint) error

	FindById(ctx context.Context, id uint) (*model.Character, error)
	FindByName(ctx context.Context, name string) (*model.Character, error)

	FindAllByOwner(ctx context.Context, ownerId string) (model.Characters, error)

	FindAll(context.Context) (model.Characters, error)

	AddPlayTime(ctx context.Context, characterId uint, amount uint64) (uint64, error)
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
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, err
	}

	if character.OptionalNewName != nil {
		currentCharacter.Name = character.OptionalNewName.(*pb.EditCharacterRequest_NewName).NewName
	}
	if character.OptionalOwnerId != nil {
		currentCharacter.OwnerId = character.OptionalOwnerId.(*pb.EditCharacterRequest_OwnerId).OwnerId
	}

	if character.OptionalPlayTime != nil {
		currentCharacter.PlayTime = character.OptionalPlayTime.(*pb.EditCharacterRequest_PlayTime).PlayTime
	}

	if character.OptionalGender != nil {
		currentCharacter.Gender = character.OptionalGender.(*pb.EditCharacterRequest_Gender).Gender
	}

	if character.OptionalRealm != nil {
		currentCharacter.Realm = character.OptionalRealm.(*pb.EditCharacterRequest_Realm).Realm
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

func (s characterService) AddPlayTime(ctx context.Context, characterId uint, amount uint64) (uint64, error) {
	character, err := s.FindById(ctx, characterId)
	if err != nil {
		return 0, err
	}

	character.PlayTime += amount
	_, err = s.repo.Save(ctx, character)

	return character.PlayTime, err
}

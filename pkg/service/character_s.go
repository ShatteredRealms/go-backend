package service

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/repository"
)

type CharacterService interface {
	Create(ctx context.Context, ownerId uint64, name string, genderId uint64, realmId uint64) (*model.Character, error)
	FindById(ctx context.Context, id uint64) (*model.Character, error)
	Edit(ctx context.Context, character *pb.Character) (*model.Character, error)
	Delete(ctx context.Context, id uint64) error
	FindAll(context.Context) ([]*model.Character, error)
	FindAllByOwner(ctx context.Context, ownerId uint64) ([]*model.Character, error)
	AddPlayTime(ctx context.Context, characterId uint64, amount uint64) (uint64, error)
}

type characterService struct {
	repo repository.CharacterRepository
}

func NewCharacterService(r repository.CharacterRepository) CharacterService {
	return characterService{
		repo: r,
	}
}

func (s characterService) Create(ctx context.Context, ownerId uint64, name string, genderId uint64, realmId uint64) (*model.Character, error) {
	character := model.Character{
		OwnerId:  ownerId,
		Name:     name,
		GenderId: genderId,
		RealmId:  realmId,
		PlayTime: 0,
	}

	if err := character.Validate(); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, &character)
}

func (s characterService) FindById(ctx context.Context, id uint64) (*model.Character, error) {
	return s.repo.FindById(ctx, id)
}

func (s characterService) Edit(ctx context.Context, character *pb.Character) (*model.Character, error) {
	currentCharacter, err := s.FindById(ctx, character.Id)
	if err != nil {
		return nil, err
	}

	if character.Name != nil {
		currentCharacter.Name = character.Name.Value
	}

	if character.Owner != nil {
		currentCharacter.OwnerId = character.Owner.Value
	}

	if character.PlayTime != nil {
		currentCharacter.PlayTime = character.PlayTime.Value
	}

	if character.Gender != nil {
		currentCharacter.GenderId = character.Gender.Value
	}

	if character.Realm != nil {
		currentCharacter.RealmId = character.Realm.Value
	}

	return s.repo.Save(ctx, currentCharacter)
}

func (s characterService) Delete(ctx context.Context, id uint64) error {
	character, err := s.FindById(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, character)
}

func (s characterService) FindAll(ctx context.Context) ([]*model.Character, error) {
	return s.repo.FindAll(ctx)
}

func (s characterService) FindAllByOwner(ctx context.Context, ownerId uint64) ([]*model.Character, error) {
	return s.repo.FindAllByOwner(ctx, ownerId)
}

func (s characterService) AddPlayTime(ctx context.Context, characterId uint64, amount uint64) (uint64, error) {
	character, err := s.FindById(ctx, characterId)
	if err != nil {
		return 0, err
	}

	character.PlayTime += amount
	_, err = s.repo.Save(ctx, character)

	return character.PlayTime, err
}

package repository

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CharacterRepository interface {
	Create(ctx context.Context, character *model.Character) (*model.Character, error)
	Save(ctx context.Context, character *model.Character) (*model.Character, error)
	Delete(ctx context.Context, character *model.Character) error

	FindById(ctx context.Context, id uint) (*model.Character, error)
	FindByName(ctx context.Context, name string) (*model.Character, error)

	FindAllByOwner(ctx context.Context, owner string) (model.Characters, error)

	FindAll(context.Context) ([]*model.Character, error)

	WithTrx(trx *gorm.DB) CharacterRepository
	Migrate(ctx context.Context) error
}

type characterRepository struct {
	DB *gorm.DB
}

func (r characterRepository) FindByName(ctx context.Context, name string) (*model.Character, error) {
	var character *model.Character = nil
	result := r.DB.WithContext(ctx).Where("name = ?", name).Find(&character)
	if result.Error != nil {
		log.WithContext(ctx).Debugf("find by name err: %v", result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		log.WithContext(ctx).Debugf("find by name: no rows affected. character: %+v", character)
		return nil, nil
	}

	log.WithContext(ctx).Debugf("character name %s found", name)
	return character, nil
}

func NewCharacterRepository(db *gorm.DB) CharacterRepository {
	return characterRepository{
		DB: db,
	}
}

func (r characterRepository) Create(ctx context.Context, character *model.Character) (*model.Character, error) {
	// Set the ID to zero so the database can generate the value
	character.ID = 0

	err := r.DB.WithContext(ctx).Create(&character).Error
	if err != nil {
		return nil, err
	}

	return character, nil
}

func (r characterRepository) Save(ctx context.Context, character *model.Character) (*model.Character, error) {
	err := r.DB.WithContext(ctx).Save(&character).Error
	if err != nil {
		return nil, err
	}

	return character, nil
}

func (r characterRepository) Delete(ctx context.Context, character *model.Character) error {
	if character == nil {
		return fmt.Errorf("character is nil")
	}
	return r.DB.WithContext(ctx).Where("id = ?", character.ID).Delete(&character).Error
}

func (r characterRepository) FindById(ctx context.Context, id uint) (*model.Character, error) {
	var character *model.Character
	result := r.DB.WithContext(ctx).Where("id = ?", id).Find(&character)
	if result.Error != nil {
		log.WithContext(ctx).Debugf("find by id err: %v", result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		log.WithContext(ctx).Debugf("find by id: no rows affected. character: %+v", character)
		return nil, nil
	}

	log.WithContext(ctx).Debugf("character id %d found", id)
	return character, nil
}

func (r characterRepository) FindAll(ctx context.Context) ([]*model.Character, error) {
	var characters []*model.Character
	return characters, r.DB.WithContext(ctx).Find(&characters).Error
}

func (r characterRepository) FindAllByOwner(ctx context.Context, owner string) (model.Characters, error) {
	var characters model.Characters
	return characters, r.DB.WithContext(ctx).Where("owner_id = ?", owner).Find(&characters).Error
}

func (r characterRepository) WithTrx(trx *gorm.DB) CharacterRepository {
	if trx == nil {
		return r
	}

	r.DB = trx
	return r
}

func (r characterRepository) Migrate(ctx context.Context) error {
	return r.DB.WithContext(ctx).AutoMigrate(&model.Character{})
}

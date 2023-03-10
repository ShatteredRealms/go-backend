package repository

import (
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"gorm.io/gorm"
)

type CharacterRepository interface {
	Create(ctx context.Context, character *model.Character) (*model.Character, error)
	Save(ctx context.Context, character *model.Character) (*model.Character, error)
	Delete(ctx context.Context, character *model.Character) error

	FindById(ctx context.Context, id uint64) (*model.Character, error)
	FindByName(ctx context.Context, name string) (*model.Character, error)

	FindAllByOwner(ctx context.Context, owner string) (model.Characters, error)

	FindAll(context.Context) ([]*model.Character, error)

	WithTrx(trx *gorm.DB) CharacterRepository
	Migrate() error
}

type characterRepository struct {
	DB *gorm.DB
}

func (r characterRepository) FindByName(ctx context.Context, name string) (*model.Character, error) {
	var character *model.Character = nil
	result := r.DB.WithContext(ctx).Where("name = ?", name).Find(&character)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

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
	return r.DB.WithContext(ctx).Delete(&character).Error
}

func (r characterRepository) FindById(ctx context.Context, id uint64) (*model.Character, error) {
	var character *model.Character = nil
	result := r.DB.WithContext(ctx).First(&character, id)
	if result.Error != nil {
		return nil, result.Error
	}

	if r.DB.RowsAffected == 0 {
		return nil, nil
	}

	return character, nil
}

func (r characterRepository) FindAll(ctx context.Context) ([]*model.Character, error) {
	var characters []*model.Character
	return characters, r.DB.WithContext(ctx).Find(&characters).Error
}

func (r characterRepository) FindAllByOwner(ctx context.Context, owner string) (model.Characters, error) {
	var characters model.Characters
	return characters, r.DB.WithContext(ctx).Where("owner = ?", owner).Find(&characters).Error
}

func (r characterRepository) WithTrx(trx *gorm.DB) CharacterRepository {
	if trx == nil {
		return r
	}

	r.DB = trx
	return r
}

func (r characterRepository) Migrate() error {
	return r.DB.AutoMigrate(&model.Character{})
}

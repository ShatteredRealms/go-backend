package repository

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/srospan"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

var (
	tracer    = otel.Tracer("character-repo")
	meter     = otel.Meter("character-repo")
	createCnt metric.Int64Counter
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

func NewCharacterRepository(db *gorm.DB) (repo CharacterRepository, err error) {
	repo = characterRepository{
		DB: db,
	}
	createCnt, err = meter.Int64Counter("sro.character.create",
		metric.WithDescription("The number of created characters by creator"),
		metric.WithUnit("{ownerId}"))
	return
}

func (r characterRepository) FindByName(ctx context.Context, name string) (*model.Character, error) {
	var character *model.Character
	result := r.DB.WithContext(ctx).Where("name = ?", name).Find(&character)
	if result.Error != nil {
		log.Logger.WithContext(ctx).Debugf("find by name err: %v", result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		log.Logger.WithContext(ctx).Debugf("find by name: no rows affected. character: %+v", character)
		return nil, nil
	}

	log.Logger.WithContext(ctx).Debugf("character name %s found", name)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(int(character.ID)),
		srospan.TargetCharacterName(character.Name),
	)
	return character, nil
}

func (r characterRepository) Create(ctx context.Context, character *model.Character) (*model.Character, error) {
	// Set the ID to zero so the database can generate the value
	character.ID = 0

	err := r.DB.WithContext(ctx).Create(&character).Error
	if err != nil {
		log.Logger.WithContext(ctx).Debugf("create error: %+v", character)
		return nil, err
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(int(character.ID)),
		srospan.TargetCharacterName(character.Name),
	)
	return character, nil
}

func (r characterRepository) Save(ctx context.Context, character *model.Character) (*model.Character, error) {
	err := r.DB.WithContext(ctx).Save(&character).Error
	if err != nil {
		return nil, err
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(int(character.ID)),
		srospan.TargetCharacterName(character.Name),
	)
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
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(int(character.ID)),
		srospan.TargetCharacterName(character.Name),
	)
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

package repository

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model/character"
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
	Create(ctx context.Context, char *character.Character) (*character.Character, error)
	Save(ctx context.Context, char *character.Character) (*character.Character, error)
	Delete(ctx context.Context, char *character.Character) error

	FindById(ctx context.Context, id uint) (*character.Character, error)
	FindByName(ctx context.Context, name string) (*character.Character, error)

	FindAllByOwner(ctx context.Context, owner string) (character.Characters, error)

	FindAll(context.Context) ([]*character.Character, error)

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

func (r characterRepository) FindByName(ctx context.Context, name string) (*character.Character, error) {
	var char *character.Character
	result := r.DB.WithContext(ctx).Where("name = ?", name).Find(&char)
	if result.Error != nil {
		log.Logger.WithContext(ctx).Debugf("find by name err: %v", result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		log.Logger.WithContext(ctx).Debugf("find by name: no rows affected. character: %+v", char)
		return nil, nil
	}

	log.Logger.WithContext(ctx).Debugf("character name %s found", name)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(int(char.ID)),
		srospan.TargetCharacterName(char.Name),
	)
	return char, nil
}

func (r characterRepository) Create(ctx context.Context, char *character.Character) (*character.Character, error) {
	// Set the ID to zero so the database can generate the value
	char.ID = 0

	err := r.DB.WithContext(ctx).Create(&char).Error
	if err != nil {
		log.Logger.WithContext(ctx).Debugf("create error: %+v", char)
		return nil, err
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(int(char.ID)),
		srospan.TargetCharacterName(char.Name),
	)
	return char, nil
}

func (r characterRepository) Save(ctx context.Context, char *character.Character) (*character.Character, error) {
	err := r.DB.WithContext(ctx).Save(&char).Error
	if err != nil {
		return nil, err
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(int(char.ID)),
		srospan.TargetCharacterName(char.Name),
	)
	return char, nil
}

func (r characterRepository) Delete(ctx context.Context, char *character.Character) error {
	if char == nil {
		return fmt.Errorf("character is nil")
	}
	return r.DB.WithContext(ctx).Where("id = ?", char.ID).Delete(&char).Error
}

func (r characterRepository) FindById(ctx context.Context, id uint) (*character.Character, error) {
	var char *character.Character
	result := r.DB.WithContext(ctx).Where("id = ?", id).Find(&char)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(int(char.ID)),
		srospan.TargetCharacterName(char.Name),
	)
	return char, nil
}

func (r characterRepository) FindAll(ctx context.Context) ([]*character.Character, error) {
	var chars []*character.Character
	return chars, r.DB.WithContext(ctx).Find(&chars).Error
}

func (r characterRepository) FindAllByOwner(ctx context.Context, owner string) (character.Characters, error) {
	var chars character.Characters
	return chars, r.DB.WithContext(ctx).Where("owner_id = ?", owner).Find(&chars).Error
}

func (r characterRepository) WithTrx(trx *gorm.DB) CharacterRepository {
	if trx == nil {
		return r
	}

	r.DB = trx
	return r
}

func (r characterRepository) Migrate(ctx context.Context) error {
	return r.DB.WithContext(ctx).AutoMigrate(&character.Character{})
}

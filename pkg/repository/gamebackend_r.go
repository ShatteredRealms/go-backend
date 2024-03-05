package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gamebackendRepository struct {
	DB *gorm.DB
}

type dimensionRepository interface {
	CreateDimension(
		ctx context.Context,
		name string,
		location string,
		version string,
		mapIds []*uuid.UUID,
	) (*model.Dimension, error)
	DuplicateDimension(ctx context.Context, refId *uuid.UUID, name string) (*model.Dimension, error)
	FindDimensionByName(ctx context.Context, name string) (*model.Dimension, error)
	FindDimensionById(ctx context.Context, id *uuid.UUID) (*model.Dimension, error)
	FindDimensionsByNames(ctx context.Context, names []string) (model.Dimensions, error)
	FindDimensionsByIds(ctx context.Context, ids []*uuid.UUID) (model.Dimensions, error)
	FindDimensionsWithMapIds(ctx context.Context, ids []*uuid.UUID) (model.Dimensions, error)
	FindAllDimensions(ctx context.Context) (model.Dimensions, error)
	SaveDimension(ctx context.Context, dimension *model.Dimension) (*model.Dimension, error)
	DeleteDimensionByName(ctx context.Context, name string) error
	DeleteDimensionById(ctx context.Context, id *uuid.UUID) error
}

type mapRepository interface {
	CreateMap(
		ctx context.Context,
		name string,
		path string,
		maxPlayers uint64,
		instanced bool,
	) (*model.Map, error)
	FindMapByName(ctx context.Context, name string) (*model.Map, error)
	FindMapById(ctx context.Context, id *uuid.UUID) (*model.Map, error)
	FindMapsByNames(ctx context.Context, names []string) (model.Maps, error)
	FindMapsByIds(ctx context.Context, ids []*uuid.UUID) (model.Maps, error)
	FindAllMaps(ctx context.Context) (model.Maps, error)
	SaveMap(ctx context.Context, m *model.Map) (*model.Map, error)
	DeleteMapByName(ctx context.Context, name string) error
	DeleteMapById(ctx context.Context, id *uuid.UUID) error
}

type GamebackendRepository interface {
	CreatePendingConnection(ctx context.Context, character string, serverName string) (*model.PendingConnection, error)
	DeletePendingConnection(ctx context.Context, id *uuid.UUID) error
	FindPendingConnection(ctx context.Context, id *uuid.UUID) *model.PendingConnection

	dimensionRepository
	mapRepository

	WithTrx(trx *gorm.DB) GamebackendRepository
	Migrate(ctx context.Context) error
}

func NewGamebackendRepository(db *gorm.DB) GamebackendRepository {
	return &gamebackendRepository{
		DB: db,
	}
}

// CreatePendingConnection implements GamebackendRepository.
func (r *gamebackendRepository) CreatePendingConnection(
	ctx context.Context,
	character string,
	serverName string,
) (*model.PendingConnection, error) {
	if character == "" {
		return nil, errors.New("no character given")
	}

	pc := &model.PendingConnection{
		Character:  character,
		ServerName: serverName,
	}
	err := r.DB.WithContext(ctx).Create(&pc).Error
	if err != nil {
		return nil, err
	}

	return pc, nil
}

// DeletePendingConnection implements GamebackendRepository.
func (r *gamebackendRepository) DeletePendingConnection(ctx context.Context, id *uuid.UUID) error {
	if id == nil {
		return fmt.Errorf("character is nil")
	}

	return r.DB.WithContext(ctx).Delete(&model.PendingConnection{}, id).Error
}

// FindPendingConnection implements GamebackendRepository.
func (r *gamebackendRepository) FindPendingConnection(ctx context.Context, id *uuid.UUID) *model.PendingConnection {
	if id == nil {
		return nil
	}

	var pendingConnection *model.PendingConnection
	result := r.DB.WithContext(ctx).Where("id = ?", id).Find(&pendingConnection)
	if result.Error != nil {
		log.Logger.WithContext(ctx).Errorf("find by id err: %v", result.Error)
		return nil
	}

	if result.RowsAffected == 0 {
		log.Logger.WithContext(ctx).Debugf("find by id not found: %s", id.String())
		return nil
	}

	log.Logger.WithContext(ctx).Debugf("found pending connection id %s", id.String())
	return pendingConnection
}

// CreateDimension implements GamebackendRepository.
func (r *gamebackendRepository) CreateDimension(
	ctx context.Context,
	name string,
	location string,
	version string,
	mapIds []*uuid.UUID,
) (*model.Dimension, error) {
	maps, err := r.FindMapsByIds(ctx, mapIds)
	if err != nil {
		return nil, err
	}

	dimension := &model.Dimension{
		Name:     name,
		Location: location,
		Version:  version,
		Maps:     maps,
	}

	if err := r.DB.WithContext(ctx).Preload(clause.Associations).Create(&dimension).Error; err != nil {
		return nil, err
	}

	return dimension, nil
}

// CreateMap implements GamebackendRepository.
func (r *gamebackendRepository) CreateMap(
	ctx context.Context,
	name string,
	path string,
	maxPlayers uint64,
	instanced bool,
) (*model.Map, error) {
	newMap := &model.Map{
		Name:       name,
		Path:       path,
		MaxPlayers: maxPlayers,
		Instanced:  instanced,
	}

	if err := r.DB.WithContext(ctx).Create(&newMap).Error; err != nil {
		return nil, err
	}

	return newMap, nil
}

// DeleteDimensionById implements GamebackendRepository.
func (r *gamebackendRepository) DeleteDimensionById(ctx context.Context, id *uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&model.Dimension{}, id).Error
}

// DeleteDimensionByName implements GamebackendRepository.
func (r *gamebackendRepository) DeleteDimensionByName(ctx context.Context, name string) error {
	return r.DB.WithContext(ctx).Delete(&model.Dimension{}, "name = ?", name).Error
}

// DeleteMapById implements GamebackendRepository.
func (r *gamebackendRepository) DeleteMapById(ctx context.Context, id *uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&model.Map{}, id).Error
}

// DeleteMapByName implements GamebackendRepository.
func (r *gamebackendRepository) DeleteMapByName(ctx context.Context, name string) error {
	return r.DB.WithContext(ctx).Delete(&model.Map{}, "name = ?", name).Error
}

// DuplicateDimension implements GamebackendRepository.
func (r *gamebackendRepository) DuplicateDimension(
	ctx context.Context,
	refId *uuid.UUID,
	name string,
) (*model.Dimension, error) {
	dimension, err := r.FindDimensionById(ctx, refId)
	if err != nil {
		return nil, err
	}

	if dimension == nil {
		return nil, model.ErrDoesNotExist
	}

	dimension.Id = nil
	dimension.Name = name
	if err = r.DB.WithContext(ctx).Preload(clause.Associations).Create(&dimension).Error; err != nil {
		return nil, err
	}

	return dimension, nil
}

// SaveDimension implements GamebackendRepository.
func (r *gamebackendRepository) SaveDimension(
	ctx context.Context,
	dimension *model.Dimension,
) (*model.Dimension, error) {
	if dimension == nil {
		return nil, fmt.Errorf("dimension nil")
	}
	log.Logger.WithContext(ctx).Infof("dimension maps: %+v", dimension.Maps)
	err := r.DB.WithContext(ctx).Save(&dimension).Error
	if err != nil {
		return nil, err
	}

	err = r.DB.WithContext(ctx).Model(&dimension).Association("Maps").Replace(dimension.Maps)
	if err != nil {
		return nil, err
	}

	return dimension, nil
}

// SaveMap implements GamebackendRepository.
func (r *gamebackendRepository) SaveMap(ctx context.Context, m *model.Map) (*model.Map, error) {
	err := r.DB.WithContext(ctx).Save(&m).Error
	if err != nil {
		return nil, err
	}

	return m, nil
}

// FindAllDimensions implements GamebackendRepository.
func (r *gamebackendRepository) FindAllDimensions(ctx context.Context) (model.Dimensions, error) {
	var dimensions model.Dimensions
	return dimensions, r.DB.WithContext(ctx).Preload(clause.Associations).Find(&dimensions).Error
}

// FindAllMaps implements GamebackendRepository.
func (r *gamebackendRepository) FindAllMaps(ctx context.Context) (model.Maps, error) {
	var maps model.Maps
	return maps, r.DB.WithContext(ctx).Find(&maps).Error
}

// FindDimensionById implements GamebackendRepository.
func (r *gamebackendRepository) FindDimensionById(ctx context.Context, id *uuid.UUID) (*model.Dimension, error) {
	var dimension *model.Dimension
	result := r.DB.WithContext(ctx).Preload(clause.Associations).Find(&dimension, id)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return dimension, nil
}

// FindDimensionByName implements GamebackendRepository.
func (r *gamebackendRepository) FindDimensionByName(ctx context.Context, name string) (*model.Dimension, error) {
	var dimension *model.Dimension
	result := r.DB.WithContext(ctx).Preload(clause.Associations).Find(&dimension, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return dimension, nil
}

// FindMapById implements GamebackendRepository.
func (r *gamebackendRepository) FindMapById(ctx context.Context, id *uuid.UUID) (*model.Map, error) {
	if id == nil {
		return nil, fmt.Errorf("error nil: id")
	}
	var m *model.Map
	result := r.DB.WithContext(ctx).Find(&m, id)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return m, nil
}

// FindMapByName implements GamebackendRepository.
func (r *gamebackendRepository) FindMapByName(ctx context.Context, name string) (*model.Map, error) {
	var m *model.Map
	result := r.DB.WithContext(ctx).Find(&m, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return m, nil
}

// FindDimensionsByIds implements GamebackendRepository.
func (r *gamebackendRepository) FindDimensionsByIds(ctx context.Context, ids []*uuid.UUID) (model.Dimensions, error) {
	var found model.Dimensions
	return found, r.DB.WithContext(ctx).Preload(clause.Associations).Find(&found, "id IN ?", ids).Error
}

// FindDimensionsByNames implements GamebackendRepository.
func (r *gamebackendRepository) FindDimensionsByNames(ctx context.Context, names []string) (model.Dimensions, error) {
	var found model.Dimensions
	return found, r.DB.WithContext(ctx).Preload(clause.Associations).Find(&found, "name IN ?", names).Error
}

// FindMapsByIds implements GamebackendRepository.
func (r *gamebackendRepository) FindMapsByIds(ctx context.Context, ids []*uuid.UUID) (model.Maps, error) {
	var found model.Maps
	return found, r.DB.WithContext(ctx).Find(&found, "id IN ?", ids).Error
}

// FindMapsByNames implements GamebackendRepository.
func (r *gamebackendRepository) FindMapsByNames(ctx context.Context, names []string) (model.Maps, error) {
	var found model.Maps
	return found, r.DB.WithContext(ctx).Find(&found, "name IN ?", names).Error
}

// FindDimensionsWithMapIds implements GamebackendRepository.
func (r *gamebackendRepository) FindDimensionsWithMapIds(ctx context.Context, ids []*uuid.UUID) (model.Dimensions, error) {
	var dimensions model.Dimensions
	return dimensions, r.DB.WithContext(ctx).
		Model(&model.Dimension{}).
		Joins("JOIN dimension_maps ON dimensions.id = dimension_maps.dimension_id").
		Where("dimension_maps.map_id IN ?", ids).
		Find(&dimensions).Error
	// Model(&model.Map{}).
	// Where("id IN ?", ids).
	// Association("Dimensions").
	// Find(&dimensions)
}

// WithTrx implmeents GamebackendRepository.
func (r *gamebackendRepository) WithTrx(trx *gorm.DB) GamebackendRepository {
	if trx == nil {
		return r
	}

	r.DB = trx
	return r
}

// Migrate implements GamebackendRepository.
func (r *gamebackendRepository) Migrate(ctx context.Context) error {
	return r.DB.WithContext(ctx).AutoMigrate(
		&model.PendingConnection{},
		&model.Dimension{},
		&model.Map{},
	)
}

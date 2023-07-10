package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
		chatTemplateIds []*uuid.UUID,
	) (*model.Dimension, error)
	DuplicateDimension(ctx context.Context, refId *uuid.UUID, name string) (*model.Dimension, error)
	FindDimensionByName(ctx context.Context, name string) (*model.Dimension, error)
	FindDimensionById(ctx context.Context, id *uuid.UUID) (*model.Dimension, error)
	FindDimensionsByNames(ctx context.Context, names []string) (model.Dimensions, error)
	FindDimensionsByIds(ctx context.Context, ids []*uuid.UUID) (model.Dimensions, error)
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

type chatTemplateRepository interface {
	CreateChatTemplate(
		ctx context.Context,
		name string,
	) (*model.ChatTemplate, error)
	FindChatTemplateByName(ctx context.Context, name string) (*model.ChatTemplate, error)
	FindChatTemplateById(ctx context.Context, id *uuid.UUID) (*model.ChatTemplate, error)
	FindChatTemplatesByNames(ctx context.Context, names []string) (model.ChatTemplates, error)
	FindChatTemplatesByIds(ctx context.Context, ids []*uuid.UUID) (model.ChatTemplates, error)
	FindAllChatTemplates(ctx context.Context) (model.ChatTemplates, error)
	SaveChatTemplate(ctx context.Context, chatTemplate *model.ChatTemplate) (*model.ChatTemplate, error)
	DeleteChatTemplateByName(ctx context.Context, name string) error
	DeleteChatTemplateById(ctx context.Context, id *uuid.UUID) error
}

type GamebackendRepository interface {
	CreatePendingConnection(ctx context.Context, character string, serverName string) (*model.PendingConnection, error)
	DeletePendingConnection(ctx context.Context, id *uuid.UUID) error
	FindPendingConnection(ctx context.Context, id *uuid.UUID) *model.PendingConnection

	dimensionRepository
	mapRepository
	chatTemplateRepository

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
	var pendingConnection *model.PendingConnection
	result := r.DB.WithContext(ctx).Where("id = ?", id).Find(&pendingConnection)
	if result.Error != nil {
		log.WithContext(ctx).Errorf("find by id err: %v", result.Error)
		return nil
	}

	if result.RowsAffected == 0 {
		log.WithContext(ctx).Debugf("find by id not found: %s", pendingConnection.Id.String())
		return nil
	}

	log.WithContext(ctx).Debugf("found pending connection id %s", id.String())
	return pendingConnection
}

// CreateChatTemplate implements GamebackendRepository.
func (r *gamebackendRepository) CreateChatTemplate(ctx context.Context, name string) (*model.ChatTemplate, error) {
	chatTemplate := &model.ChatTemplate{
		Name: name,
	}

	if err := r.DB.WithContext(ctx).Create(&chatTemplate).Error; err != nil {
		return nil, err
	}

	return chatTemplate, nil
}

// CreateDimension implements GamebackendRepository.
func (r *gamebackendRepository) CreateDimension(
	ctx context.Context,
	name string,
	location string,
	version string,
	mapIds []*uuid.UUID,
	chatTemplateIds []*uuid.UUID,
) (*model.Dimension, error) {
	maps, err := r.FindMapsByIds(ctx, mapIds)
	if err != nil {
		return nil, err
	}

	chatTemplates, err := r.FindChatTemplatesByIds(ctx, chatTemplateIds)
	if err != nil {
		return nil, err
	}

	dimension := &model.Dimension{
		Name:           name,
		ServerLocation: location,
		Version:        version,
		Maps:           maps,
		ChatTemplates:  chatTemplates,
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

// DeleteChatTemplateById implements GamebackendRepository.
func (r *gamebackendRepository) DeleteChatTemplateById(ctx context.Context, id *uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&model.ChatTemplate{}, id).Error
}

// DeleteChatTemplateByName implements GamebackendRepository.
func (r *gamebackendRepository) DeleteChatTemplateByName(ctx context.Context, name string) error {
	return r.DB.WithContext(ctx).Delete(&model.ChatTemplate{}, "name = ?", name).Error
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

// SaveChatTemplate implements GamebackendRepository.
func (r *gamebackendRepository) SaveChatTemplate(
	ctx context.Context,
	chatTemplate *model.ChatTemplate,
) (*model.ChatTemplate, error) {
	err := r.DB.WithContext(ctx).Save(&chatTemplate).Error
	if err != nil {
		return nil, err
	}

	return chatTemplate, nil
}

// SaveDimension implements GamebackendRepository.
func (r *gamebackendRepository) SaveDimension(
	ctx context.Context,
	dimension *model.Dimension,
) (*model.Dimension, error) {
	err := r.DB.WithContext(ctx).Preload(clause.Associations).Save(&dimension).Error
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

// FindAllChatTemplates implements GamebackendRepository.
func (r *gamebackendRepository) FindAllChatTemplates(ctx context.Context) (model.ChatTemplates, error) {
	var chatTemplates model.ChatTemplates
	return chatTemplates, r.DB.WithContext(ctx).Find(&chatTemplates).Error
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

// FindChatTemplateById implements GamebackendRepository.
func (r *gamebackendRepository) FindChatTemplateById(ctx context.Context, id *uuid.UUID) (*model.ChatTemplate, error) {
	var chatTemplate *model.ChatTemplate
	result := r.DB.WithContext(ctx).Find(&chatTemplate, id)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return chatTemplate, nil
}

// FindChatTemplateByName implements GamebackendRepository.
func (r *gamebackendRepository) FindChatTemplateByName(ctx context.Context, name string) (*model.ChatTemplate, error) {
	var chatTemplate *model.ChatTemplate
	result := r.DB.WithContext(ctx).Find(&chatTemplate, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return chatTemplate, nil
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

// FindChatTemplatesByIds implements GamebackendRepository.
func (r *gamebackendRepository) FindChatTemplatesByIds(
	ctx context.Context,
	ids []*uuid.UUID,
) (model.ChatTemplates, error) {
	var found model.ChatTemplates
	return found, r.DB.WithContext(ctx).Find(&found, ids).Error
}

// FindChatTemplatesByNames implements GamebackendRepository.
func (r *gamebackendRepository) FindChatTemplatesByNames(
	ctx context.Context,
	names []string,
) (model.ChatTemplates, error) {
	var found model.ChatTemplates
	return found, r.DB.WithContext(ctx).Find(&found, "name IN ?", names).Error
}

// FindDimensionsByIds implements GamebackendRepository.
func (r *gamebackendRepository) FindDimensionsByIds(ctx context.Context, ids []*uuid.UUID) (model.Dimensions, error) {
	var found model.Dimensions
	return found, r.DB.WithContext(ctx).Preload(clause.Associations).Find(&found, ids).Error
}

// FindDimensionsByNames implements GamebackendRepository.
func (r *gamebackendRepository) FindDimensionsByNames(ctx context.Context, names []string) (model.Dimensions, error) {
	var found model.Dimensions
	return found, r.DB.WithContext(ctx).Preload(clause.Associations).Find(&found, "name IN ?", names).Error
}

// FindMapsByIds implements GamebackendRepository.
func (r *gamebackendRepository) FindMapsByIds(ctx context.Context, ids []*uuid.UUID) (model.Maps, error) {
	var found model.Maps
	return found, r.DB.WithContext(ctx).Find(&found, ids).Error
}

// FindMapsByNames implements GamebackendRepository.
func (r *gamebackendRepository) FindMapsByNames(ctx context.Context, names []string) (model.Maps, error) {
	var found model.Maps
	return found, r.DB.WithContext(ctx).Find(&found, "name IN ?", names).Error
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
		&model.ChatTemplate{},
		&model.Map{},
	)
}

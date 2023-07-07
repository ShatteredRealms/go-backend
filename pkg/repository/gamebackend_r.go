package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type gamebackendRepository struct {
	DB *gorm.DB
}

type GamebackendRepository interface {
	CreatePendingConnection(ctx context.Context, character string, serverName string) (*model.PendingConnection, error)
	DeletePendingConnection(ctx context.Context, id *uuid.UUID) error
	FindPendingConnection(ctx context.Context, id *uuid.UUID) *model.PendingConnection

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
	return r.DB.WithContext(ctx).AutoMigrate(&model.PendingConnection{})
}

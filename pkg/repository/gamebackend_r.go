package repository

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"gorm.io/gorm"
)

type gamebackendRepository struct {
	DB *gorm.DB
}

type GamebackendRepository interface {
	RegisterUser(ctx context.Context)
	Migrate(ctx context.Context) error
}

func (r gamebackendRepository) Migrate(ctx context.Context) error {
	return r.DB.WithContext(ctx).AutoMigrate(&model.NewCharacterConnection{})
}

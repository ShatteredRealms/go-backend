package repository

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"gorm.io/gorm"
)

type ChatRepository interface {
	AllChannels(ctx context.Context) ([]*model.ChatChannel, error)
	GetChannel(ctx context.Context, id uint) (*model.ChatChannel, error)
	CreateChannel(ctx context.Context, channel *model.ChatChannel) (*model.ChatChannel, error)
	DeleteChannel(ctx context.Context, channel *model.ChatChannel) error
	UpdateChannel(ctx context.Context, channel *model.ChatChannel) error

	FindDeletedWithName(ctx context.Context, name string) (*model.ChatChannel, error)

	Migrate(ctx context.Context) error
}

type chatRepository struct {
	DB *gorm.DB
}

func (r chatRepository) UpdateChannel(ctx context.Context, channel *model.ChatChannel) error {
	return r.DB.WithContext(ctx).Save(&channel).Error
}

func (r chatRepository) AllChannels(ctx context.Context) ([]*model.ChatChannel, error) {
	var channels []*model.ChatChannel
	r.DB.WithContext(ctx).Find(&channels)
	return channels, r.DB.Error
}

func (r chatRepository) CreateChannel(ctx context.Context, channel *model.ChatChannel) (*model.ChatChannel, error) {
	return channel, r.DB.WithContext(ctx).Create(&channel).Error
}

func (r chatRepository) DeleteChannel(ctx context.Context, channel *model.ChatChannel) error {
	return r.DB.WithContext(ctx).Unscoped().Delete(channel).Error
}

func (r chatRepository) Migrate(ctx context.Context) error {
	return r.DB.WithContext(ctx).AutoMigrate(&model.ChatChannel{})
}

func (r chatRepository) GetChannel(ctx context.Context, id uint) (*model.ChatChannel, error) {
	var channel *model.ChatChannel
	r.DB.WithContext(ctx).Where("id = ?", id).Find(&channel)
	return channel, r.DB.Error
}

func (r chatRepository) FindDeletedWithName(ctx context.Context, name string) (*model.ChatChannel, error) {
	var channel *model.ChatChannel
	r.DB.WithContext(ctx).Unscoped().Where("name = ?", name).Find(&channel)
	return channel, r.DB.Error
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return chatRepository{DB: db}
}

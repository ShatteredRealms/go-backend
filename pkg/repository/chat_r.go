package repository

import (
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"gorm.io/gorm"
)

type ChatRepository interface {
	AllChannels(ctx context.Context) (model.ChatChannels, error)
	FindChannelById(ctx context.Context, id uint) (*model.ChatChannel, error)
	CreateChannel(ctx context.Context, channel *model.ChatChannel) (*model.ChatChannel, error)
	DeleteChannel(ctx context.Context, channel *model.ChatChannel) error
	FullDeleteChannel(ctx context.Context, channel *model.ChatChannel) error
	UpdateChannel(ctx context.Context, channel *model.ChatChannel) (*model.ChatChannel, error)

	FindDeletedWithName(ctx context.Context, name string) (*model.ChatChannel, error)

	AuthorizedChannelsForCharacter(ctx context.Context, characterId uint) (model.ChatChannels, error)
	ChangeAuthorizationForCharacter(ctx context.Context, characterId uint, channelIds []uint, addAuth bool) error

	Migrate(ctx context.Context) error
}

type chatRepository struct {
	DB *gorm.DB
}

func (r chatRepository) ChangeAuthorizationForCharacter(ctx context.Context, characterId uint, channelIds []uint, addAuth bool) error {
	if addAuth {
		tx := r.DB.Begin()
		for _, id := range channelIds {
			if err := tx.Create(&model.ChatChannelPermission{
				ChannelId:   id,
				CharacterId: characterId,
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		tx.Commit()
		return nil
	}

	return r.DB.Delete(&model.ChatChannelPermission{}, "character_id = ? AND channel_id IN ?", characterId, channelIds).Error
}

func (r chatRepository) AuthorizedChannelsForCharacter(ctx context.Context, characterId uint) (model.ChatChannels, error) {
	var channels model.ChatChannels
	r.DB.WithContext(ctx).
		Model(&model.ChatChannel{}).
		Joins("JOIN chat_channel_permissions ON chat_channels.id = chat_channel_permissions.channel_id").
		Where("chat_channel_permissions.character_id = ?", characterId).
		Find(&channels)
	return channels, r.DB.Error
}

func (r chatRepository) UpdateChannel(ctx context.Context, channel *model.ChatChannel) (*model.ChatChannel, error) {
	return channel, r.DB.WithContext(ctx).Save(&channel).Error
}

func (r chatRepository) AllChannels(ctx context.Context) (model.ChatChannels, error) {
	var channels model.ChatChannels
	r.DB.WithContext(ctx).Find(&channels)
	return channels, r.DB.Error
}

func (r chatRepository) CreateChannel(ctx context.Context, channel *model.ChatChannel) (*model.ChatChannel, error) {
	return channel, r.DB.WithContext(ctx).Create(&channel).Error
}

func (r chatRepository) DeleteChannel(ctx context.Context, channel *model.ChatChannel) error {
	return r.DB.WithContext(ctx).Delete(channel).Error
}

func (r chatRepository) FullDeleteChannel(ctx context.Context, channel *model.ChatChannel) error {
	return r.DB.WithContext(ctx).Unscoped().Delete(channel).Error
}

func (r chatRepository) Migrate(ctx context.Context) error {
	err := r.DB.WithContext(ctx).AutoMigrate(&model.ChatChannel{})
	if err != nil {
		return err
	}

	return r.DB.WithContext(ctx).AutoMigrate(&model.ChatChannelPermission{})
}

func (r chatRepository) FindChannelById(ctx context.Context, id uint) (*model.ChatChannel, error) {
	var channel *model.ChatChannel
	result := r.DB.WithContext(ctx).Where("id = ?", id).Find(&channel)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return channel, nil
}

func (r chatRepository) FindDeletedWithName(ctx context.Context, name string) (*model.ChatChannel, error) {
	var channel *model.ChatChannel
	result := r.DB.WithContext(ctx).Unscoped().Where("name = ?", name).Find(&channel)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return channel, nil
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return chatRepository{DB: db}
}

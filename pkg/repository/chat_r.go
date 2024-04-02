package repository

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/model/chat"
	"github.com/ShatteredRealms/go-backend/pkg/srospan"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type ChatRepository interface {
	AllChannels(ctx context.Context) (chat.ChatChannels, error)
	FindChannelById(ctx context.Context, id uint) (*chat.ChatChannel, error)
	CreateChannel(ctx context.Context, channel *chat.ChatChannel) (*chat.ChatChannel, error)
	DeleteChannel(ctx context.Context, channel *chat.ChatChannel) error
	FullDeleteChannel(ctx context.Context, channel *chat.ChatChannel) error
	UpdateChannel(ctx context.Context, channel *chat.ChatChannel) (*chat.ChatChannel, error)

	FindDeletedWithName(ctx context.Context, name string) (*chat.ChatChannel, error)

	AuthorizedChannelsForCharacter(ctx context.Context, characterId uint) (chat.ChatChannels, error)
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
			if err := tx.Create(&chat.ChatChannelPermission{
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

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(int(characterId)),
	)

	return r.DB.Delete(&chat.ChatChannelPermission{}, "character_id = ? AND channel_id IN ?", characterId, channelIds).Error
}

func (r chatRepository) AuthorizedChannelsForCharacter(ctx context.Context, characterId uint) (chat.ChatChannels, error) {
	var channels chat.ChatChannels
	r.DB.WithContext(ctx).
		Model(&chat.ChatChannel{}).
		Joins("JOIN chat_channel_permissions ON chat_channels.id = chat_channel_permissions.channel_id").
		Where("chat_channel_permissions.character_id = ?", characterId).
		Find(&channels)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.TargetCharacterId(int(characterId)),
	)

	return channels, r.DB.Error
}

func (r chatRepository) UpdateChannel(ctx context.Context, channel *chat.ChatChannel) (*chat.ChatChannel, error) {
	trace.SpanFromContext(ctx).SetAttributes(srospan.ChatChannelAttributes(channel)...)
	return channel, r.DB.WithContext(ctx).Save(&channel).Error
}

func (r chatRepository) AllChannels(ctx context.Context) (chat.ChatChannels, error) {
	var channels chat.ChatChannels
	r.DB.WithContext(ctx).Find(&channels)
	return channels, r.DB.Error
}

func (r chatRepository) CreateChannel(ctx context.Context, channel *chat.ChatChannel) (*chat.ChatChannel, error) {
	trace.SpanFromContext(ctx).SetAttributes(srospan.ChatChannelAttributes(channel)...)
	return channel, r.DB.WithContext(ctx).Create(&channel).Error
}

func (r chatRepository) DeleteChannel(ctx context.Context, channel *chat.ChatChannel) error {
	trace.SpanFromContext(ctx).SetAttributes(srospan.ChatChannelAttributes(channel)...)
	return r.DB.WithContext(ctx).Delete(channel).Error
}

func (r chatRepository) FullDeleteChannel(ctx context.Context, channel *chat.ChatChannel) error {
	trace.SpanFromContext(ctx).SetAttributes(srospan.ChatChannelAttributes(channel)...)
	return r.DB.WithContext(ctx).Unscoped().Delete(channel).Error
}

func (r chatRepository) Migrate(ctx context.Context) error {
	return r.DB.WithContext(ctx).AutoMigrate(&chat.ChatChannel{}, &chat.ChatChannelPermission{})
}

func (r chatRepository) FindChannelById(ctx context.Context, id uint) (*chat.ChatChannel, error) {
	var channel *chat.ChatChannel
	result := r.DB.WithContext(ctx).Where("id = ?", id).Find(&channel)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	trace.SpanFromContext(ctx).SetAttributes(srospan.ChatChannelAttributes(channel)...)
	return channel, nil
}

func (r chatRepository) FindDeletedWithName(ctx context.Context, name string) (*chat.ChatChannel, error) {
	var channel *chat.ChatChannel
	result := r.DB.WithContext(ctx).Unscoped().Where("name = ?", name).Find(&channel)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	trace.SpanFromContext(ctx).SetAttributes(srospan.ChatChannelAttributes(channel)...)
	return channel, nil
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return chatRepository{DB: db}
}

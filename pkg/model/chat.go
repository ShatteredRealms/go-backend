package model

import (
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"gorm.io/gorm"
)

type ChatChannel struct {
	gorm.Model
	Name      string `gorm:"unique_index:idx_channel" json:"name"`
	Dimension string `gorm:"unique_index:idx_channel" json:"dimension"`
	Public    bool   `json:"public"`
}
type ChatChannels []*ChatChannel

type ChatChannelPermission struct {
	ChannelId     uint   `gorm:"unique_index:idx_permission" json:"channelId"`
	CharacterName string `gorm:"unique_index:idx_permission" json:"characterId"`
}
type ChatChannelPermissions []*ChatChannelPermission

func (c *ChatChannel) ToPb() *pb.ChatChannel {
	return &pb.ChatChannel{
		Name:   c.Name,
		Public: c.Public,
	}
}

func (c ChatChannels) ToPb() *pb.ChatChannels {
	resp := &pb.ChatChannels{Channels: make([]*pb.ChatChannel, len(c))}
	for idx, channel := range c {
		resp.Channels[idx] = channel.ToPb()
	}

	return resp
}

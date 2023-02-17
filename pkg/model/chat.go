package model

import (
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"gorm.io/gorm"
)

type ChatChannel struct {
	gorm.Model
	Name      string `gorm:"unique_index:idx_channel" json:"name"`
	Dimension string `gorm:"unique_index:idx_channel" json:"dimension"`
	Public    bool   `json:"public"`
}

type ChatChannelPermission struct {
	ChannelId     uint   `gorm:"unique_index:idx_permission" json:"channelId"`
	CharacterName string `gorm:"unique_index:idx_permission" json:"characterId"`
}

func (c ChatChannel) ToPb() *pb.ChatChannel {
	return &pb.ChatChannel{
		Name:   c.Name,
		Public: c.Public,
	}
}

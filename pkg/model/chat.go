package model

import (
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"gorm.io/gorm"
)

type ChatChannel struct {
	gorm.Model
	Name   string `gorm:"unique" json:"name"`
	Public bool   `json:"public"`
}

type ChatChannelPermission struct {
	ChannelId   uint `gorm:"unique_index:idx_permission" json:"channelId"`
	CharacterId uint `gorm:"unique_index:idx_permission" json:"characterId"`
}

func (c ChatChannel) ToPb() *pb.ChatChannel {
	return &pb.ChatChannel{
		Id:     uint64(c.ID),
		Name:   c.Name,
		Public: c.Public,
	}
}

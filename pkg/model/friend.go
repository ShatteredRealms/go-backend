package model

import "time"

type UserFriend struct {
	Username       string    `gorm:"primarykey" json:"userId"`
	FriendUsername string    `gorm:"primarykey" json:"friendUserId"`
	CreatedAt      time.Time `json:"createdAt"`
}

type UserBlock struct {
	Username        string    `gorm:"primarykey" json:"userId"`
	BlockedUsername string    `gorm:"primarykey" json:"blockedUserId"`
	CreatedAt       time.Time `json:"createdAt"`
}

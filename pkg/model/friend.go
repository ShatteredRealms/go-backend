package model

import "time"

type UserFriend struct {
	UserId       uint      `gorm:"unique_index:idx_friend" json:"userId"`
	FriendUserId uint      `gorm:"unique_index:idx_friend" json:"friendUserId"`
	CreatedAt    time.Time `json:"createdAt"`
}

type UserBlock struct {
	UserId        uint      `gorm:"unique_index:idx_block" json:"userId"`
	BlockedUserId uint      `gorm:"unique_index:idx_block" json:"blockedUserId"`
	CreatedAt     time.Time `json:"createdAt"`
}

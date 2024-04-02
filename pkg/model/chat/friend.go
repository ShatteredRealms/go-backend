package chat

import "time"

type UserFriend struct {
	OwningId  string    `gorm:"primarykey" json:"user"`
	FriendId  string    `gorm:"primarykey" json:"friendUser"`
	CreatedAt time.Time `json:"createdAt"`
}
type UserFriends []*UserFriend

type UserBlock struct {
	OwningId  string    `gorm:"primarykey" json:"user"`
	BlockedId string    `gorm:"primarykey" json:"blockedUser"`
	CreatedAt time.Time `json:"createdAt"`
}
type UserBlocks []*UserBlock

package model

import (
	"time"

	"github.com/google/uuid"
)

type PendingConnection struct {
	// Id secret used by a server to lookup a pending connection
	Id *uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`

	// Character name of the character that owns it
	Character string

	// ServerName the name of the server the character is assigned to
	ServerName string

	// CreatedAt when the pending connection was created
	CreatedAt time.Time
}

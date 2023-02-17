package model

import (
	"gorm.io/gorm"
	"time"
)

type Dimension struct {
	Name     string `gorm:"primarykey" json:"name"`
	Location string `json:"location"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

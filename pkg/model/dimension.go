package model

import (
	"errors"
	"fmt"
	goaway "github.com/TwiN/go-away"
	"gorm.io/gorm"
	"regexp"
)

const (
	MinDimensionNameLength = 1
	MaxDimensionNameLength = 64
)

var (
	DimensionNameRegex, _ = regexp.Compile("^[a-zA-Z0-9]+$")

	// ErrDimensionNameToShort ErrDimensionNameToLong thrown when a character name is too short
	ErrDimensionNameToShort = errors.New(fmt.Sprintf("name must be at least %d characters", MinCharacterNameLength))

	// ErrDimensionNameToLong thrown when a character name is too long
	ErrDimensionNameToLong = errors.New(fmt.Sprintf("name can be at most %d characters", MaxCharacterNameLength))
)

// Dimension Also known as a "server". A cluster of servers that a player can choose to play in. Each dimension is
// separate and cannot interact with other dimensions. This is recursive, meaning all entities that are tied to a
// dimension cannot interact with other dimensions.
type Dimension struct {
	gorm.Model
	Name           string `gorm:"unique" json:"name"`
	ServerLocation string `gorm:"not null" json:"serverLocation"`
}

func (c *Dimension) ValidateLocation() error {
	if _, ok := ServerLocations[c.ServerLocation]; ok {
		return nil
	}

	return ErrInvalidServerLocation
}

func (c *Dimension) ValidateName() error {
	if len(c.Name) < MinDimensionNameLength {
		return ErrDimensionNameToShort
	}

	if len(c.Name) > MaxDimensionNameLength {
		return ErrDimensionNameToLong
	}

	if !DimensionNameRegex.MatchString(c.Name) {
		return ErrInvalidName
	}

	if goaway.IsProfane(c.Name) {
		return ErrNameProfane
	}

	return nil
}

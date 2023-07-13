package model

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/pb"
	goaway "github.com/TwiN/go-away"
	"gorm.io/plugin/soft_delete"
)

const (
	MinCharacterNameLength = 3
	MaxCharacterNameLength = 20
)

var (
	CharacterNameRegex, _ = regexp.Compile("^[a-zA-Z0-9]+$")

	// ErrCharacterNameToShort thrown when a character name is too short
	ErrCharacterNameToShort = errors.New(fmt.Sprintf("name must be at least %d characters", MinCharacterNameLength))

	// ErrCharacterNameToLong thrown when a character name is too long
	ErrCharacterNameToLong = errors.New(fmt.Sprintf("name can be at most %d characters", MaxCharacterNameLength))
)

type Character struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:udx_name"`

	// Owner The username/account that owns the character
	OwnerId string `gorm:"not null" json:"owner"`
	Name    string `gorm:"not null;uniqueIndex:udx_name" json:"name"`
	Gender  string `gorm:"not null" json:"gender"`
	Realm   string `gorm:"not null" json:"realm"`

	// PlayTime Time in minutes the character has played
	PlayTime uint64 `gorm:"not null" json:"play_time"`

	// Location last location recorded for the character
	Location Location `gorm:"type:bytes;serializer:gob" json:"location"`
}
type Characters []*Character

func (c *Character) Validate() error {
	if err := c.ValidateGender(); err != nil {
		return err
	}

	if err := c.ValidateRealm(); err != nil {
		return err
	}

	return c.ValidateName()
}

func (c *Character) ValidateName() error {
	if len(c.Name) < MinCharacterNameLength {
		return ErrCharacterNameToShort
	}

	if len(c.Name) > MaxCharacterNameLength {
		return ErrCharacterNameToLong
	}

	if !CharacterNameRegex.MatchString(c.Name) {
		return ErrInvalidName
	}

	if goaway.IsProfane(c.Name) {
		return ErrNameProfane
	}

	return nil
}

func (c *Character) ValidateGender() error {
	if c.Gender == "Male" ||
		c.Gender == "Female" {
		return nil
	}

	return ErrInvalidGender
}

func (c *Character) ValidateRealm() error {
	if c.Realm == "Human" ||
		c.Realm == "Cyborg" {
		return nil
	}

	return ErrInvalidRealm
}

func (c *Character) ToPb() *pb.CharacterDetails {
	return &pb.CharacterDetails{
		Id:       uint64(c.ID),
		Owner:    c.OwnerId,
		Name:     c.Name,
		Gender:   c.Gender,
		Realm:    c.Realm,
		PlayTime: c.PlayTime,
		Location: c.Location.ToPb(),
	}
}

func (c Characters) ToPb() *pb.CharactersDetails {
	resp := &pb.CharactersDetails{Characters: make([]*pb.CharacterDetails, len(c))}
	for idx, character := range c {
		resp.Characters[idx] = character.ToPb()
	}

	return resp
}

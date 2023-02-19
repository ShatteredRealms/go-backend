package model

import (
	"errors"
	"fmt"
	goaway "github.com/TwiN/go-away"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"gorm.io/gorm"
	"regexp"
)

const (
	MinNameLength = 3
	MaxNameLength = 20

	MaxRealmId  = 2
	MaxGenderId = 3
)

var (
	Genders = &pb.Genders{
		Genders: []*pb.Gender{
			{
				Name: "Male",
				Id:   1,
			},
			{
				Name: "Female",
				Id:   2,
			},
			{
				Name: "None",
				Id:   3,
			},
		},
	}

	Realms = &pb.Realms{
		Realms: []*pb.Realm{
			{
				Name: "Human",
				Id:   1,
			},
			{
				Name: "Cyborg",
				Id:   2,
			},
		},
	}

	NameRegex, _ = regexp.Compile("^[a-zA-Z0-9]+$")

	ErrInvalidNameCharacter = errors.New("name contains invalid character(s)")

	// ErrNameProfane thrown when a name is profane
	ErrNameProfane = errors.New("name unavailable")

	// ErrInvalidRealm thrown when a character belongs to an unknown realm
	ErrInvalidRealm = errors.New("invalid realm")

	// ErrInvalidGender thrown when a character belongs to an unknown gender
	ErrInvalidGender = errors.New("invalid gender")

	// ErrNameToShort thrown when a character name is too short
	ErrNameToShort = errors.New(fmt.Sprintf("name must be at least %d characters", MinNameLength))

	// ErrNameToLong thrown when a character name is too long
	ErrNameToLong = errors.New(fmt.Sprintf("name can be at most %d characters", MaxNameLength))
)

type Character struct {
	gorm.Model

	// Owner The username/account that owns the character
	Owner    string `gorm:"not null" json:"owner"`
	Name     string `gorm:"not null;unique" json:"name"`
	GenderId uint64 `gorm:"not null" json:"gender_id"`
	RealmId  uint64 `gorm:"not null" json:"realm_id"`

	// PlayTime Time in minutes the character has played
	PlayTime uint64 `gorm:"not null" json:"play_time"`

	// Location last location recorded for the character
	Location Location `gorm:"type:bytes;serializer:gob" json:"location"`
}

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
	if len(c.Name) < MinNameLength {
		return ErrNameToShort
	}

	if len(c.Name) > MaxNameLength {
		return ErrNameToLong
	}

	if !NameRegex.MatchString(c.Name) {
		return ErrInvalidNameCharacter
	}

	if goaway.IsProfane(c.Name) {
		return ErrNameProfane
	}

	return nil
}

func (c *Character) ValidateGender() error {
	for _, v := range Genders.Genders {
		if c.GenderId == v.Id {
			return nil
		}
	}

	return ErrInvalidGender
}

func (c *Character) ValidateRealm() error {
	for _, v := range Realms.Realms {
		if c.RealmId == v.Id {
			return nil
		}
	}

	return ErrInvalidRealm
}

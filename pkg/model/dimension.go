package model

import (
	"errors"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

const (
	MinDimensionNameLength = 1
	MaxDimensionNameLength = 64
)

var (
	// ErrDimensionNameToShort ErrDimensionNameToLong thrown when a character name is too short
	ErrDimensionNameToShort = errors.New(fmt.Sprintf("name must be at least %d characters", MinCharacterNameLength))

	// ErrDimensionNameToLong thrown when a character name is too long
	ErrDimensionNameToLong = errors.New(fmt.Sprintf("name can be at most %d characters", MaxCharacterNameLength))
)

// Dimension Also known as a "server". A cluster of servers that a player can choose to play in. Each dimension is
// separate and cannot interact with other dimensions. This is recursive, meaning all entities that are tied to a
// dimension cannot interact with other dimensions.
type Dimension struct {
	Model
	Name     string `gorm:"unique"`
	Location string `gorm:"not null"`
	Version  string `gorm:"not null"`
	Maps     Maps   `gorm:"many2many:dimension_maps"`
}

type Dimensions []*Dimension

func (c *Dimension) ValidateLocation() error {
	if _, ok := ServerLocations[c.Location]; ok {
		return nil
	}

	return ErrInvalidServerLocation
}

func (dimension *Dimension) ValidateName() error {
	if len(dimension.Name) < MinDimensionNameLength {
		return ErrDimensionNameToShort
	}

	if len(dimension.Name) > MaxDimensionNameLength {
		return ErrDimensionNameToLong
	}

	return nil
}

func (dimension *Dimension) ToPb() *pb.Dimension {
	maps := make([]*pb.Map, len(dimension.Maps))
	for idx, m := range dimension.Maps {
		maps[idx] = m.ToPb()
	}

	return &pb.Dimension{
		Id:       dimension.Id.String(),
		Name:     dimension.Name,
		Version:  dimension.Version,
		Maps:     dimension.Maps.ToPb(),
		Location: dimension.Location,
	}
}

func (dimensions Dimensions) ToPb() []*pb.Dimension {
	out := make([]*pb.Dimension, len(dimensions))
	for idx, c := range dimensions {
		out[idx] = c.ToPb()
	}

	return out
}

func (dimension *Dimension) GetImageName() string {
	version := "latest"
	if dimension.Version != "" {
		version = dimension.Version
	}

	return fmt.Sprintf("779965382548.dkr.ecr.us-east-1.amazonaws.com/sro/game:%s", version)
}

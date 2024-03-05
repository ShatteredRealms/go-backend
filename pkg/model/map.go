package model

import "github.com/ShatteredRealms/go-backend/pkg/pb"

type Map struct {
	Model
	Name       string `gorm:"unique"`
	Path       string `gorm:"not null"`
	MaxPlayers uint64
	Instanced  bool
	Dimensions Dimensions `gorm:"many2many:dimension_maps"`
}

type Maps []*Map

func (m *Map) ToPb() *pb.Map {
	return &pb.Map{
		Id:         m.Id.String(),
		Name:       m.Name,
		Path:       m.Path,
		MaxPlayers: m.MaxPlayers,
		Instanced:  m.Instanced,
	}
}

func (maps Maps) ToPb() []*pb.Map {
	out := make([]*pb.Map, len(maps))
	for idx, m := range maps {
		out[idx] = m.ToPb()
	}

	return out
}

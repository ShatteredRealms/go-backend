package model

import "github.com/ShatteredRealms/go-backend/pkg/pb"

type Location struct {
	World string  `json:"world"`
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
	Z     float32 `json:"z"`
	Roll  float32 `json:"roll"`
	Pitch float32 `json:"pitch"`
	Yaw   float32 `json:"yaw"`
}

func (l Location) ToPb() *pb.Location {
	return &pb.Location{
		World: l.World,
		X:     l.X,
		Y:     l.Y,
		Z:     l.Z,
		Roll:  l.Roll,
		Pitch: l.Pitch,
		Yaw:   l.Yaw,
	}
}

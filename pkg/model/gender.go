package model

import "github.com/ShatteredRealms/go-backend/pkg/pb"

var (
	Genders = map[string]struct{}{
		"Male":   {},
		"Female": {},
	}

	genderPbs = make([]*pb.Gender, len(Realms))
)

func GetGenders() []*pb.Gender {
	if len(genderPbs) == 0 {
		idx := 0
		for genderName := range Genders {
			genderPbs[idx] = &pb.Gender{Name: genderName}
			idx++
		}
	}

	return genderPbs
}

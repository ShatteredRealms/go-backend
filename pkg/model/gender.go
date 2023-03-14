package model

import "github.com/ShatteredRealms/go-backend/pkg/pb"

var (
	Genders = map[string]struct{}{
		"Male":   {},
		"Female": {},
	}

	genderPbs          = make([]*pb.Gender, len(Realms))
	gendersInitialized = false
)

func GetGenders() []*pb.Gender {
	if !gendersInitialized {
		idx := 0
		for genderName := range Genders {
			genderPbs[idx] = &pb.Gender{Name: genderName}
			idx++
		}
		gendersInitialized = true
	}

	return genderPbs
}

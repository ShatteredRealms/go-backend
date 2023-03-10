package model

import "github.com/ShatteredRealms/go-backend/pkg/pb"

var (
	Realms = map[string]struct{}{
		"Human":  {},
		"Cyborg": {},
	}

	realmsPb = make([]*pb.Realm, len(Realms))
)

func GetRealms() []*pb.Realm {
	if len(realmsPb) == 0 {
		idx := 0
		for realmName := range Realms {
			realmsPb[idx] = &pb.Realm{Name: realmName}
			idx++
		}
	}

	return realmsPb
}

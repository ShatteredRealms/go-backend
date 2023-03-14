package model

import (
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

var (
	Realms = map[string]struct{}{
		"Human":  {},
		"Cyborg": {},
	}

	realmsPb          = make([]*pb.Realm, len(Realms))
	realmsInitialized = false
)

func GetRealms() []*pb.Realm {
	if !realmsInitialized {
		idx := 0
		for realmName := range Realms {
			fmt.Printf("realm: %v", realmName)
			realmsPb[idx] = &pb.Realm{Name: realmName}
			idx++
		}
		realmsInitialized = true
	}

	return realmsPb
}

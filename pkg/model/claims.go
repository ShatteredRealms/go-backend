package model

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v4"
)

var (
	CharactersClientId  = "sro-characters"
	ChatClientId        = "sro-chat"
	GamebackendClientId = "sro-gamebackend"
)

type ResourceRoles struct {
	Roles []string `json:"roles"`
}

type SROClaims struct {
	jwt.RegisteredClaims
	RealmRoles     ClaimRoles          `json:"realm_roles,omitempty"`
	ResourceAccess ClaimResourceAccess `json:"resource_access"`
	Username       string              `json:"preferred_username,omitempty"`
}

type ClaimRoles struct {
	Roles []string `json:"roles,omitempty"`
}
type ClaimResourceAccess map[string]ClaimRoles

func (s SROClaims) HasResourceRole(role *gocloak.Role, clientId string) bool {
	if resource, ok := s.ResourceAccess[clientId]; ok {
		for _, claimRole := range resource.Roles {
			if *role.Name == claimRole {
				return true
			}
		}
	}

	return false
}

func (s SROClaims) HasRole(role *gocloak.Role) bool {
	for _, claimRole := range s.RealmRoles.Roles {
		if *role.Name == claimRole {
			return true
		}
	}

	return false
}

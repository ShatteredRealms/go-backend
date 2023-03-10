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
	RealmRoles       []string `json:"realm_roles.roles,omitempty"`
	CharacterRoles   []string `json:"resource_access.sro-characters.roles,omitempty"`
	GameBackendRoles []string `json:"resource_access.sro-gamebackend.roles,omitempty"`
	ChatRoles        []string `json:"resource_access.sro-chat.roles,omitempty"`
}

func (s SROClaims) HasRole(role *gocloak.Role, clientId string) bool {
	if clientId == GamebackendClientId {
		return containsKey(s.GameBackendRoles, *role.Name)
	}

	if clientId == CharactersClientId {
		return containsKey(s.CharacterRoles, *role.Name)
	}

	if clientId == ChatClientId {
		return containsKey(s.ChatRoles, *role.Name)
	}

	return containsKey(s.RealmRoles, *role.Name)
}

func containsKey(arr []string, key string) bool {
	for _, currentKey := range arr {
		if key == currentKey {
			return true
		}
	}

	return false
}

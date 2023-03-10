package model

import "github.com/golang-jwt/jwt/v4"

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

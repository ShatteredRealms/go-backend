package srv

import "github.com/ShatteredRealms/go-backend/pkg/model"

func registerRole(roles []*model.RoleRepresentation, newRole *model.RoleRepresentation) *model.RoleRepresentation {
	roles = append(roles, newRole)
	return newRole
}

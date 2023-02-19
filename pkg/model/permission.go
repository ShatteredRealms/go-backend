package model

import (
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// TODO(wil) Refactor Permission/Other inside UserPermission and RolePermission and refactor functions that use these
//           structs.

// TODO(wil) Refactor `Other` to `Global` for readability and conveying meaning?

// UserPermission Database model for customized user and permissions join table
type UserPermission struct {
	// User The User to grant the permission to
	Username string `gorm:"primaryKey" json:"username"`

	// Permission The permission that is assigned to the user
	Permission string `gorm:"primaryKey" json:"permission"`

	// Whether the permission applies to users besides itself. If true, then the permission applies even if
	// the target of the method is not itself
	Other bool `gorm:"not null" json:"other"`
}
type UserPermissions []*UserPermission

// RolePermission Database model for customized role and permissions join table
type RolePermission struct {
	// Role The Role to grant the permission to
	Role string `gorm:"primaryKey" json:"role"`

	// Permission The permission that is assigned to the user
	Permission string `gorm:"primaryKey" json:"permission"`

	// Whether the permission applies to users besides itself. If true, then the permission applies even if
	// the target of the method is not itself
	Other bool `gorm:"not null" json:"other"`
}
type RolePermissions []*RolePermission

func (up *UserPermission) ToPB() *pb.UserPermission {
	return &pb.UserPermission{
		Permission: &wrapperspb.StringValue{Value: up.Permission},
		Other:      up.Other,
	}
}

func (rp *RolePermission) ToPB() *pb.UserPermission {
	return &pb.UserPermission{
		Permission: &wrapperspb.StringValue{Value: rp.Permission},
		Other:      rp.Other,
	}
}

func (ups UserPermissions) ToPB() *pb.UserPermissions {
	permissions := make([]*pb.UserPermission, len(ups))
	for i, permission := range ups {
		permissions[i] = permission.ToPB()
	}

	return &pb.UserPermissions{Permissions: permissions}
}

func (rps RolePermissions) ToPB() *pb.UserPermissions {
	permissions := make([]*pb.UserPermission, len(rps))
	for i, permission := range rps {
		permissions[i] = permission.ToPB()
	}

	return &pb.UserPermissions{Permissions: permissions}
}

package model

import (
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"gorm.io/gorm"
	"time"
)

const (
	MinRoleNameLength = 3
	MaxRoleNameLength = 255
)

type Role struct {
	Name string `gorm:"primarykey" json:"name"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Roles []*Role

func (r *Role) Validate() error {
	err := r.validateName()
	if err != nil {
		return err
	}

	return nil
}

func (r *Role) validateName() error {
	if len(r.Name) < MinRoleNameLength {
		return fmt.Errorf("minimum name length is %d", MinRoleNameLength)
	}

	if len(r.Name) > MaxRoleNameLength {
		return fmt.Errorf("maximum name length is %d", MaxRoleNameLength)
	}

	return nil
}

func (r *Role) ToPB() *pb.UserRole {
	return &pb.UserRole{
		Name: r.Name,
	}
}

func (roles Roles) ToPB() *pb.UserRoles {
	out := make([]*pb.UserRole, len(roles))
	for i, role := range roles {
		out[i] = role.ToPB()
	}

	return &pb.UserRoles{Roles: out}
}

package service

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/repository"
	"gorm.io/gorm"
)

type PermissionService interface {
	AddPermissionForUser(ctx context.Context, permission *model.UserPermission) error
	AddPermissionForRole(ctx context.Context, permission *model.RolePermission) error

	RemPermissionForUser(ctx context.Context, permission *model.UserPermission) error
	RemPermissionForRole(ctx context.Context, permission *model.RolePermission) error

	FindPermissionsForUsername(ctx context.Context, username string) model.UserPermissions
	FindPermissionsForRole(ctx context.Context, role string) model.RolePermissions

	ClearPermissionsForRole(ctx context.Context, role string) error
	ClearPermissionsForUsername(ctx context.Context, username string) error

	ResetPermissionsForRole(ctx context.Context, role string, permissions []*model.RolePermission) error
	ResetPermissionsForUsername(ctx context.Context, username string, permissions []*model.UserPermission) error

	UpdateRoleName(ctx context.Context, oldName string, newName string) error

	WithTrx(*gorm.DB) PermissionService
	Migrate() error
}

type permissionService struct {
	permissionRepository repository.PermissionRepository
}

func (s permissionService) UpdateRoleName(ctx context.Context, oldName string, newName string) error {
	return s.permissionRepository.UpdateRoleName(ctx, oldName, newName)
}

func NewPermissionService(r repository.PermissionRepository) PermissionService {
	return permissionService{
		permissionRepository: r,
	}
}

func (s permissionService) AddPermissionForUser(ctx context.Context, permission *model.UserPermission) error {
	return s.permissionRepository.AddPermissionForUser(ctx, permission)
}

func (s permissionService) AddPermissionForRole(ctx context.Context, permission *model.RolePermission) error {
	return s.permissionRepository.AddPermissionForRole(ctx, permission)
}

func (s permissionService) RemPermissionForUser(ctx context.Context, permission *model.UserPermission) error {
	return s.permissionRepository.RemPermissionForUser(ctx, permission)
}

func (s permissionService) RemPermissionForRole(ctx context.Context, permission *model.RolePermission) error {
	return s.permissionRepository.RemPermissionForRole(ctx, permission)
}

func (s permissionService) WithTrx(db *gorm.DB) PermissionService {
	s.permissionRepository = s.permissionRepository.WithTrx(db)
	return s
}

func (s permissionService) FindPermissionsForUsername(ctx context.Context, username string) model.UserPermissions {
	return s.permissionRepository.FindPermissionsForUsername(ctx, username)
}

func (s permissionService) FindPermissionsForRole(ctx context.Context, role string) model.RolePermissions {
	return s.permissionRepository.FindPermissionsForRole(ctx, role)
}

func (s permissionService) ClearPermissionsForRole(ctx context.Context, role string) error {
	return s.permissionRepository.ClearPermissionsForRole(ctx, role)
}

func (s permissionService) ClearPermissionsForUsername(ctx context.Context, username string) error {
	return s.permissionRepository.ClearPermissionsForUsername(ctx, username)
}

func (s permissionService) ResetPermissionsForRole(ctx context.Context, role string, permissions []*model.RolePermission) error {
	err := s.ClearPermissionsForRole(ctx, role)
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		if permission.Role == role {
			err = s.AddPermissionForRole(ctx, permission)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s permissionService) ResetPermissionsForUsername(ctx context.Context, username string, permissions []*model.UserPermission) error {
	err := s.ClearPermissionsForUsername(ctx, username)
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		if permission.Username == username {
			err = s.AddPermissionForUser(ctx, permission)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s permissionService) Migrate() error {
	return s.permissionRepository.Migrate()
}

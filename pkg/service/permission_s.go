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

    FindPermissionsForUserID(ctx context.Context, id uint) model.UserPermissions
    FindPermissionsForRoleID(ctx context.Context, id uint) model.RolePermissions

    ClearPermissionsForRole(ctx context.Context, id uint) error
    ClearPermissionsForUser(ctx context.Context, id uint) error

    ResetPermissionsForRole(ctx context.Context, id uint, permissions []*model.RolePermission) error
    ResetPermissionsForUser(ctx context.Context, id uint, permissions []*model.UserPermission) error

    WithTrx(*gorm.DB) PermissionService
    Migrate() error
}

type permissionService struct {
    permissionRepository repository.PermissionRepository
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

func (s permissionService) FindPermissionsForUserID(ctx context.Context, id uint) model.UserPermissions {
    return s.permissionRepository.FindPermissionsForUserID(ctx, id)
}

func (s permissionService) FindPermissionsForRoleID(ctx context.Context, id uint) model.RolePermissions {
    return s.permissionRepository.FindPermissionsForRoleID(ctx, id)
}

func (s permissionService) ClearPermissionsForRole(ctx context.Context, id uint) error {
    return s.permissionRepository.ClearPermissionsForRole(ctx, id)
}

func (s permissionService) ClearPermissionsForUser(ctx context.Context, id uint) error {
    return s.permissionRepository.ClearPermissionsForUser(ctx, id)
}

func (s permissionService) ResetPermissionsForRole(ctx context.Context, id uint, permissions []*model.RolePermission) error {
    err := s.ClearPermissionsForRole(ctx, id)
    if err != nil {
        return err
    }

    for _, permission := range permissions {
        if permission.RoleID == id {
            err = s.AddPermissionForRole(ctx, permission)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func (s permissionService) ResetPermissionsForUser(ctx context.Context, id uint, permissions []*model.UserPermission) error {
    err := s.ClearPermissionsForUser(ctx, id)
    if err != nil {
        return err
    }

    for _, permission := range permissions {
        if permission.UserID == 0 || permission.UserID == id {
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

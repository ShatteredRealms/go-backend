package repository

import (
    "context"
    "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
    "gorm.io/gorm"
)

type PermissionRepository interface {
    AddPermissionForUser(ctx context.Context, permission *model.UserPermission) error
    AddPermissionForRole(ctx context.Context, permission *model.RolePermission) error

    RemPermissionForUser(ctx context.Context, permission *model.UserPermission) error
    RemPermissionForRole(ctx context.Context, permission *model.RolePermission) error

    FindPermissionsForUserID(ctx context.Context, id uint) model.UserPermissions
    FindPermissionsForRoleID(ctx context.Context, id uint) model.RolePermissions

    ClearPermissionsForRole(ctx context.Context, id uint) error
    ClearPermissionsForUser(ctx context.Context, id uint) error

    WithTrx(*gorm.DB) PermissionRepository
    Migrate() error
}

type permissionRepository struct {
    DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
    return permissionRepository{
        DB: db,
    }
}

func (r permissionRepository) AddPermissionForUser(ctx context.Context, permission *model.UserPermission) error {
    return r.DB.WithContext(ctx).Create(&permission).Error
}

func (r permissionRepository) AddPermissionForRole(ctx context.Context, permission *model.RolePermission) error {
    return r.DB.WithContext(ctx).Create(&permission).Error
}

func (r permissionRepository) RemPermissionForUser(ctx context.Context, permission *model.UserPermission) error {
    return r.DB.WithContext(ctx).Delete(&permission).Error
}

func (r permissionRepository) RemPermissionForRole(ctx context.Context, permission *model.RolePermission) error {
    return r.DB.WithContext(ctx).Delete(&permission).Error
}

func (r permissionRepository) FindPermissionsForUserID(ctx context.Context, id uint) model.UserPermissions {
    var permissions model.UserPermissions
    r.DB.WithContext(ctx).Where("user_id = ?", id).Find(&permissions)
    return permissions
}

func (r permissionRepository) FindPermissionsForRoleID(ctx context.Context, id uint) model.RolePermissions {
    var permissions model.RolePermissions
    r.DB.WithContext(ctx).Where("role_id = ?", id).Find(&permissions)
    return permissions
}

func (r permissionRepository) ClearPermissionsForRole(ctx context.Context, id uint) error {
    return r.DB.WithContext(ctx).Delete(&model.RolePermission{}, "role_id = ?", id).Error
}

func (r permissionRepository) ClearPermissionsForUser(ctx context.Context, id uint) error {
    return r.DB.WithContext(ctx).Delete(&model.UserPermission{}, "user_id = ?", id).Error
}

func (r permissionRepository) WithTrx(trx *gorm.DB) PermissionRepository {
    if trx == nil {
        return r
    }

    r.DB = trx
    return r
}

func (r permissionRepository) Migrate() error {
    err := r.DB.AutoMigrate(&model.UserPermission{})
    if err != nil {
        return err
    }

    return r.DB.AutoMigrate(&model.RolePermission{})
}

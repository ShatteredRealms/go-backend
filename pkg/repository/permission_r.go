package repository

import (
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	AddPermissionForUser(ctx context.Context, permission *model.UserPermission) error
	AddPermissionForRole(ctx context.Context, permission *model.RolePermission) error

	RemPermissionForUser(ctx context.Context, permission *model.UserPermission) error
	RemPermissionForRole(ctx context.Context, permission *model.RolePermission) error

	FindPermissionsForUsername(ctx context.Context, username string) model.UserPermissions
	FindPermissionsForRole(ctx context.Context, role string) model.RolePermissions

	ClearPermissionsForRole(ctx context.Context, role string) error
	ClearPermissionsForUsername(ctx context.Context, username string) error

	WithTrx(*gorm.DB) PermissionRepository
	Migrate() error
	UpdateRoleName(ctx context.Context, oldName string, newName string) error
}

type permissionRepository struct {
	DB *gorm.DB
}

func (r permissionRepository) UpdateRoleName(ctx context.Context, oldName string, newName string) error {
	return r.DB.WithContext(ctx).Where("name = ?", oldName).Update("name", newName).Error
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

func (r permissionRepository) FindPermissionsForUsername(ctx context.Context, username string) model.UserPermissions {
	var permissions model.UserPermissions
	r.DB.WithContext(ctx).Where("username = ?", username).Find(&permissions)
	return permissions
}

func (r permissionRepository) FindPermissionsForRole(ctx context.Context, role string) model.RolePermissions {
	var permissions model.RolePermissions
	r.DB.WithContext(ctx).Where("role = ?", role).Find(&permissions)
	return permissions
}

func (r permissionRepository) ClearPermissionsForRole(ctx context.Context, role string) error {
	return r.DB.WithContext(ctx).Delete(&model.RolePermission{}, "role = ?", role).Error
}

func (r permissionRepository) ClearPermissionsForUsername(ctx context.Context, username string) error {
	return r.DB.WithContext(ctx).Delete(&model.UserPermission{}, "username = ?", username).Error
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

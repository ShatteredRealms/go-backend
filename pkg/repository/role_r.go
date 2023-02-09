package repository

import (
    "context"
    "fmt"
    "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
    "gorm.io/gorm"
)

type RoleRepository interface {
    Create(context.Context, *model.Role) (*model.Role, error)
    Save(context.Context, *model.Role) (*model.Role, error)
    Delete(context.Context, *model.Role) error
    Update(context.Context, *model.Role) error

    All(context.Context) model.Roles
    FindById(ctx context.Context, id uint) *model.Role
    FindByName(ctx context.Context, name string) *model.Role

    WithTrx(*gorm.DB) RoleRepository
    Migrate() error
}

type roleRepository struct {
    DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
    return roleRepository{
        DB: db,
    }
}

func (r roleRepository) Create(ctx context.Context, role *model.Role) (*model.Role, error) {
    err := role.Validate()
    if err != nil {
        return nil, err
    }

    existingRoleWithName := r.FindByName(ctx, role.Name)
    if existingRoleWithName != nil && existingRoleWithName.ID != 0 {
        return nil, fmt.Errorf("name already exists")
    }

    err = r.DB.WithContext(ctx).Create(&role).Error

    return role, err
}

func (r roleRepository) Save(ctx context.Context, role *model.Role) (*model.Role, error) {
    existingRoleWithName := r.FindByName(ctx, role.Name)
    if existingRoleWithName != nil {
        return nil, fmt.Errorf("name already exists")
    }

    return role, r.DB.WithContext(ctx).Save(&role).Error
}

func (r roleRepository) Delete(ctx context.Context, role *model.Role) error {
    return r.DB.WithContext(ctx).Delete(&role).Error
}

func (r roleRepository) Update(ctx context.Context, role *model.Role) error {
    return r.DB.WithContext(ctx).Model(&role).Update("name", role.Name).Error
}

func (r roleRepository) All(ctx context.Context) model.Roles {
    var roles model.Roles
    r.DB.WithContext(ctx).Find(&roles)
    return roles
}

func (r roleRepository) FindById(ctx context.Context, id uint) *model.Role {
    var role *model.Role
    r.DB.WithContext(ctx).WithContext(ctx).Where("id = ?", id).Find(&role)
    return role
}

func (r roleRepository) FindByName(ctx context.Context, name string) *model.Role {
    var role *model.Role
    r.DB.WithContext(ctx).WithContext(ctx).Where("name = ?", name).Find(&role)
    return role
}

func (r roleRepository) WithTrx(trx *gorm.DB) RoleRepository {
    if trx == nil {
        return r
    }

    r.DB = trx
    return r
}

func (r roleRepository) Migrate() error {
    return r.DB.AutoMigrate(&model.Role{})
}

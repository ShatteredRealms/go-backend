package service

import (
    "context"
    "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
    "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/repository"
    "gorm.io/gorm"
)

type RoleService interface {
    Create(context.Context, *model.Role) (*model.Role, error)
    Save(context.Context, *model.Role) (*model.Role, error)
    Delete(context.Context, *model.Role) error
    Update(context.Context, *model.Role) error

    All(context.Context) []*model.Role
    FindById(ctx context.Context, id uint) *model.Role
    FindByName(ctx context.Context, name string) *model.Role

    WithTrx(*gorm.DB) RoleService

    FindAll(ctx context.Context) model.Roles
}

type roleService struct {
    roleRepository repository.RoleRepository
}

func NewRoleService(r repository.RoleRepository) RoleService {
    return roleService{
        roleRepository: r,
    }
}

func (s roleService) Create(ctx context.Context, role *model.Role) (*model.Role, error) {
    return s.roleRepository.Create(ctx, role)
}

func (s roleService) Save(ctx context.Context, role *model.Role) (*model.Role, error) {
    return s.roleRepository.Save(ctx, role)
}

func (s roleService) Delete(ctx context.Context, role *model.Role) error {
    return s.roleRepository.Delete(ctx, role)
}

func (s roleService) Update(ctx context.Context, role *model.Role) error {
    if len(role.Name) == 0 {
        return nil
    }

    return s.roleRepository.Update(ctx, role)
}

func (s roleService) All(ctx context.Context) []*model.Role {
    return s.roleRepository.All(ctx)
}

func (s roleService) FindById(ctx context.Context, id uint) *model.Role {
    return s.roleRepository.FindById(ctx, id)
}

func (s roleService) FindByName(ctx context.Context, name string) *model.Role {
    return s.roleRepository.FindByName(ctx, name)
}

func (s roleService) WithTrx(db *gorm.DB) RoleService {
    s.roleRepository = s.roleRepository.WithTrx(db)
    return s
}

func (s roleService) FindAll(ctx context.Context) model.Roles {
    return s.roleRepository.All(ctx)
}

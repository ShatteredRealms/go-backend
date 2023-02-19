package service

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"gorm.io/gorm"
)

type UserService interface {
	Create(context.Context, *model.User) (*model.User, error)
	Save(context.Context, *model.User) (*model.User, error)
	AddToRole(context.Context, *model.User, *model.Role) error
	RemFromRole(context.Context, *model.User, *model.Role) error
	WithTrx(*gorm.DB) UserService
	FindByEmail(ctx context.Context, email string) []*model.User
	FindByUsername(ctx context.Context, username string) *model.User
	FindAll(context.Context) []*model.User
	Ban(ctx context.Context, user *model.User) error
	UnBan(ctx context.Context, user *model.User) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return userService{
		userRepository: r,
	}
}

func (u userService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	return u.userRepository.Create(ctx, user)
}

func (u userService) Save(ctx context.Context, user *model.User) (*model.User, error) {
	return u.userRepository.Save(ctx, user)
}

func (u userService) AddToRole(ctx context.Context, user *model.User, role *model.Role) error {
	return u.userRepository.AddToRole(ctx, user, role)
}

func (u userService) RemFromRole(ctx context.Context, user *model.User, role *model.Role) error {
	return u.userRepository.RemFromRole(ctx, user, role)
}

func (u userService) WithTrx(trx *gorm.DB) UserService {
	u.userRepository = u.userRepository.WithTrx(trx)
	return u
}

func (u userService) FindByEmail(ctx context.Context, email string) []*model.User {
	return u.userRepository.FindByEmail(ctx, email)
}

func (u userService) FindByUsername(ctx context.Context, username string) *model.User {
	return u.userRepository.FindByUsername(ctx, username)
}

func (u userService) FindAll(ctx context.Context) []*model.User {
	return u.userRepository.All(ctx)
}

func (u userService) Ban(ctx context.Context, user *model.User) error {
	return u.userRepository.Ban(ctx, user)
}

func (u userService) UnBan(ctx context.Context, user *model.User) error {
	return u.userRepository.UnBan(ctx, user)
}

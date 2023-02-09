package repository

import (
	"context"
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"gopkg.in/nullbio/null.v4"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

type UserRepository interface {
	Create(context.Context, *model.User) (*model.User, error)
	Save(context.Context, *model.User) (*model.User, error)
	AddToRole(context.Context, *model.User, *model.Role) error
	RemFromRole(context.Context, *model.User, *model.Role) error
	WithTrx(*gorm.DB) UserRepository
	FindById(ctx context.Context, id uint) *model.User
	FindByEmail(ctx context.Context, email string) []*model.User
	FindByUsername(ctx context.Context, username string) *model.User
	Migrate() error
	All(context.Context) []*model.User
	Ban(ctx context.Context, user *model.User) error
	UnBan(ctx context.Context, user *model.User) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return userRepository{
		DB: db,
	}
}

func (u userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	err := user.Validate()
	if err != nil {
		return user, err
	}

	conflict := u.FindByUsername(ctx, user.Username)
	if conflict.Exists() {
		return user, fmt.Errorf("username is already taken")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
	if err != nil {
		return user, fmt.Errorf("password: %w", err)
	}

	user.Password = string(hashedPass)
	err = u.DB.WithContext(ctx).Create(&user).Error

	return user, err
}

func (u userRepository) Save(ctx context.Context, user *model.User) (*model.User, error) {
	conflict := u.FindByUsername(ctx, user.Username)
	if conflict.Exists() && user.ID != conflict.ID {
		return user, fmt.Errorf("username is already taken")
	}

	return user, u.DB.WithContext(ctx).Save(&user).Error
}

func (u userRepository) AddToRole(ctx context.Context, user *model.User, role *model.Role) error {
	return u.DB.WithContext(ctx).Model(&user).Association("Roles").Append([]model.Role{*role})
}

func (u userRepository) RemFromRole(ctx context.Context, user *model.User, role *model.Role) error {
	return u.DB.WithContext(ctx).Model(&user).Association("Roles").Delete([]model.Role{*role})
}

func (u userRepository) WithTrx(trx *gorm.DB) UserRepository {
	if trx == nil {
		return u
	}

	u.DB = trx
	return u
}

func (u userRepository) FindById(ctx context.Context, id uint) *model.User {
	var user *model.User
	u.DB.WithContext(ctx).Where("id=?", id).Preload("Roles").Find(&user)
	return user
}

func (u userRepository) FindByEmail(ctx context.Context, email string) []*model.User {
	var users []*model.User
	u.DB.WithContext(ctx).Where("email=?", email).Preload("Roles").Find(&users)
	return users
}

func (u userRepository) FindByUsername(ctx context.Context, username string) *model.User {
	var user *model.User
	u.DB.WithContext(ctx).Where("username=?", username).Preload("Roles").Find(&user)
	return user
}

func (u userRepository) Migrate() error {
	_ = u.DB.AutoMigrate(&model.User{})
	return nil
}

func (u userRepository) All(ctx context.Context) []*model.User {
	var users []*model.User
	u.DB.WithContext(ctx).Preload("Roles").Find(&users)
	return users
}

func (u userRepository) Ban(ctx context.Context, user *model.User) error {
	user.BannedAt = null.TimeFrom(time.Now())
	return u.DB.WithContext(ctx).Save(&user).Error
}

func (u userRepository) UnBan(ctx context.Context, user *model.User) error {
	user.BannedAt = null.Time{}
	return u.DB.WithContext(ctx).Save(&user).Error
}

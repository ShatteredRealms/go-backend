package mocks

import (
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
)

type UserService struct {
	CreateReturn struct {
		user *model.User
		err  error
	}
	SaveReturn struct {
		user *model.User
		err  error
	}
	AddToRoleReturn   error
	RemFromRoleReturn error
	FindByIdReturn    *model.User
	FindByEmailReturn *model.User
	FindAllReturn     []*model.User
	BanReturn         error
	UnBanReturn       error
}

func (t UserService) Create(user *model.User) (*model.User, error) {
	return t.CreateReturn.user, t.CreateReturn.err
}

func (t UserService) Save(user *model.User) (*model.User, error) {
	return t.SaveReturn.user, t.SaveReturn.err
}

func (t UserService) AddToRole(user *model.User, role *model.Role) error {
	return t.AddToRoleReturn
}

func (t UserService) RemFromRole(user *model.User, role *model.Role) error {
	return t.RemFromRoleReturn
}

func (t UserService) FindById(id uint) *model.User {
	return t.FindByIdReturn
}

func (t UserService) FindByEmail(email string) *model.User {
	return t.FindByEmailReturn
}

func (t UserService) FindAll() []*model.User {
	return t.FindAllReturn
}

func (t UserService) Ban(user *model.User) error {
	return t.BanReturn
}

func (t UserService) UnBan(user *model.User) error {
	return t.UnBanReturn
}

package factory

import (
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"github.com/brianvoe/gofakeit/v6"
)

type Factory interface {
	NewCharacter() *model.Character
	NewBaseUser() *model.User
	NewUser() *model.User
	Factory() *gofakeit.Faker
}

type factory struct {
	faker *gofakeit.Faker
}

func NewFactory() Factory {
	return &factory{faker: gofakeit.NewCrypto()}
}

func (f *factory) Factory() *gofakeit.Faker {
	return f.faker
}

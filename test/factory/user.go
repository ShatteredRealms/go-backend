package factory

import (
	"github.com/ShatteredRealms/go-backend/pkg/model"
)

func (f *factory) NewBaseUser() *model.User {
	return &model.User{
		FirstName: f.faker.FirstName(),
		LastName:  f.faker.LastName(),
		Email:     f.faker.Email(),
		Password:  f.faker.Password(true, true, true, true, true, f.faker.IntRange(10, 20)),
		Username:  f.faker.Username(),
	}
}

func (f *factory) NewUser() *model.User {
	user := f.NewBaseUser()
	user.ID = f.faker.UintRange(1, 10000)
	user.CreatedAt = f.faker.Date()
	user.UpdatedAt = f.faker.Date()

	return user
}

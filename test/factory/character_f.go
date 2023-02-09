package factory

import (
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"gorm.io/gorm"
)

// NewCharacter Creates a new character with fake data. The character will never be deleted.
func (f *factory) NewCharacter() *model.Character {
	character := f.newCharacterHelper()
	for character.Validate() != nil {
		character = f.newCharacterHelper()
	}

	return character
}

func (f *factory) newCharacterHelper() *model.Character {
	return &model.Character{
		Model: gorm.Model{
			ID:        f.faker.UintRange(1, 100),
			CreatedAt: f.faker.Date(),
			UpdatedAt: f.faker.Date(),
			DeletedAt: gorm.DeletedAt{},
		},
		OwnerId:  uint64(f.faker.UintRange(1, 100)),
		Name:     f.faker.Username(),
		GenderId: uint64(f.faker.UintRange(1, model.MaxGenderId)),
		RealmId:  uint64(f.faker.UintRange(1, model.MaxRealmId)),
		PlayTime: uint64(f.faker.UintRange(0, 10^6)),
	}
}

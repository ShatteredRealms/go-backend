package model_test

import (
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/test/factory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Character", func() {
	f := factory.NewFactory()
	var character *model.Character

	BeforeEach(func() {
		character = f.NewCharacter()
	})

	Context("Validation", func() {
		It("should pass if valid", func() {
			Expect(character.Validate()).To(BeNil())
		})

		It(fmt.Sprintf("should min length of %d", model.MinNameLength), func() {
			Expect(character.Validate()).To(BeNil())
			character.Name = character.Name[:model.MinNameLength-1]
			Expect(character.Validate()).To(Equal(model.ErrNameToShort))
		})

		It(fmt.Sprintf("should max length of %d", model.MaxNameLength), func() {
			Expect(character.Validate()).To(BeNil())
			for len(character.Name) <= model.MaxNameLength {
				character.Name = character.Name + f.Factory().Letter()
			}
			Expect(character.Validate()).To(Equal(model.ErrNameToLong))
		})
		for _, specialCharacter := range []string{
			"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "=",
			"[", "{", "]", "}", "\\", "|",
			";", ":", "'", "\"",
			"<", ",", ".", ">", "/", "?",
		} {
			It(fmt.Sprintf("should not allow special character %s in name", specialCharacter), func() {
				Expect(character.Validate()).To(BeNil(), fmt.Sprintf("name %s should be valid", character.Name))
				character.Name = character.Name + specialCharacter
				Expect(character.Validate()).To(Equal(model.ErrInvalidNameCharacter))
			})
		}

		It("should not allow profane names", func() {
			character.Name = "fuck"
			Expect(character.Validate()).To(Equal(model.ErrNameProfane))
		})

		It("should require a valid realm", func() {
			character.RealmId = model.MaxRealmId + 1
			Expect(character.Validate()).To(Equal(model.ErrInvalidRealm))
		})

		It("should require a valid gender", func() {
			character.GenderId = model.MaxGenderId + 1
			Expect(character.Validate()).To(Equal(model.ErrInvalidGender))
		})
	})
})

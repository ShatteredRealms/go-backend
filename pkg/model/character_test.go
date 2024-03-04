package model_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

var _ = Describe("Character model", func() {
	var (
		character = &model.Character{}
	)

	BeforeEach(func() {
		Expect(faker.FakeData(&character)).To(Succeed())
		character.Name = "name"
		character.Realm = "Cyborg"
		character.Gender = "Male"
	})

	Describe("Validation", func() {
		Context("issues", func() {
			It("wrong gender should error", func() {
				character.Gender = faker.Email()
				Expect(character.Validate()).To(MatchError(model.ErrInvalidGender))
			})

			It("wrong realm should error", func() {
				character.Realm = faker.Email()
				Expect(character.Validate()).To(MatchError(model.ErrInvalidRealm))
			})

			Context("with name", func() {
				It("should error while under minimum length", func() {
					character.Name = "a"
					Expect(character.Validate()).To(MatchError(model.ErrCharacterNameToShort))
				})
				It("should error while above maximum length", func() {
					character.Name = "aaaaaaaaaaaaaaaaaaaaa"
					Expect(character.Validate()).To(MatchError(model.ErrCharacterNameToLong))
				})
				It("should only allow letters and numbers", func() {
					character.Name = "name@"
					Expect(character.Validate()).To(MatchError(model.ErrInvalidName))
					character.Name = "!name"
					Expect(character.Validate()).To(MatchError(model.ErrInvalidName))
					character.Name = " name"
					Expect(character.Validate()).To(MatchError(model.ErrInvalidName))
					character.Name = "name_"
					Expect(character.Validate()).To(MatchError(model.ErrInvalidName))
				})

				It("shouldn't allow profanity", func() {
					character.Name = "fuck"
					Expect(character.Validate()).To(MatchError(model.ErrNameProfane))
				})
			})
		})
		It("should not error for valid character", func() {
			Expect(character.Validate()).To(Succeed())
		})
	})

	validateCharacter := (func(char *model.Character, pb *pb.CharacterDetails) {
		Expect(pb.Id).To(Equal(uint64(char.ID)))
		Expect(pb.Owner).To(Equal(char.OwnerId))
		Expect(pb.Name).To(Equal(char.Name))
		Expect(pb.Gender).To(Equal(char.Gender))
		Expect(pb.Realm).To(Equal(char.Realm))
		Expect(pb.PlayTime).To(Equal(char.PlayTime))
		Expect(pb.Location).NotTo(BeNil())
	})

	Describe("ToPb", func() {
		It("should convert single character to protobuf and retain all fields", func() {
			out := character.ToPb()
			validateCharacter(character, out)
		})

		It("should convert arrays character to protobuf and retain all fields", func() {
			var characters model.Characters
			characters = make([]*model.Character, 10)
			for idx := range characters {
				characters[idx] = &model.Character{}
				faker.FakeData(characters[idx])
			}
			out := characters.ToPb()
			Expect(out.Characters).To(HaveLen(len(characters)))
			for idx := range characters {
				validateCharacter(characters[idx], out.Characters[idx])
			}
		})
	})
})

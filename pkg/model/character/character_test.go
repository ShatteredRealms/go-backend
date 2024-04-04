package character_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/model/character"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

var _ = Describe("Character model", func() {
	var (
		char = &character.Character{}
	)

	BeforeEach(func() {
		Expect(faker.FakeData(&char)).To(Succeed())
		char.Name = "name"
		char.Realm = "Cyborg"
		char.Gender = "Male"
	})

	Describe("Validation", func() {
		Context("issues", func() {
			It("wrong gender should error", func() {
				char.Gender = faker.Email()
				Expect(char.Validate()).To(MatchError(common.ErrInvalidGender))
			})

			It("wrong realm should error", func() {
				char.Realm = faker.Email()
				Expect(char.Validate()).To(MatchError(common.ErrInvalidRealm))
			})

			Context("with name", func() {
				It("should error while under minimum length", func() {
					char.Name = "a"
					Expect(char.Validate()).To(MatchError(character.ErrCharacterNameToShort))
				})
				It("should error while above maximum length", func() {
					char.Name = "aaaaaaaaaaaaaaaaaaaaa"
					Expect(char.Validate()).To(MatchError(character.ErrCharacterNameToLong))
				})
				It("should only allow letters and numbers", func() {
					char.Name = "name@"
					Expect(char.Validate()).To(MatchError(common.ErrInvalidName))
					char.Name = "!name"
					Expect(char.Validate()).To(MatchError(common.ErrInvalidName))
					char.Name = " name"
					Expect(char.Validate()).To(MatchError(common.ErrInvalidName))
					char.Name = "name_"
					Expect(char.Validate()).To(MatchError(common.ErrInvalidName))
				})

				It("shouldn't allow profanity", func() {
					char.Name = "fuck"
					Expect(char.Validate()).To(MatchError(common.ErrNameProfane))
				})
			})
		})
		It("should not error for valid character", func() {
			Expect(char.Validate()).To(Succeed())
		})
	})

	validateCharacter := (func(char *character.Character, pb *pb.CharacterDetails) {
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
			out := char.ToPb()
			validateCharacter(char, out)
		})

		It("should convert arrays character to protobuf and retain all fields", func() {
			var characters character.Characters
			characters = make([]*character.Character, 10)
			for idx := range characters {
				characters[idx] = &character.Character{}
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

package repository_test

import (
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/kend/pkg/repository"
	"github.com/kend/test/factory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("Character Repository", func() {
	var (
		repo      repository.CharacterRepository
		character *model.Character
	)

	f := factory.NewFactory()

	BeforeEach(func() {
		repo = repository.NewCharacterRepository(DB)
		Expect(repo.Migrate()).To(Succeed())

		character = f.NewCharacter()
		character.ID = 0
	})

	Context("Create", func() {
		It("should be work with valid input", func() {
			dbCharacter, err := repo.Create(character)
			Expect(err).To(Succeed())
			Expect(dbCharacter).NotTo(BeNil())
			Expect(dbCharacter.ID).To(Equal(uint(1)))
		})

		It("should fail with duplicate ids", func() {
			dbCharacter, err := repo.Create(character)
			Expect(err).To(Succeed())
			Expect(dbCharacter).NotTo(BeNil())

			dbCharacter2, err := repo.Create(character)
			Expect(err).NotTo(Succeed())
			Expect(dbCharacter2).To(BeNil())
		})
	})

	Context("Save", func() {
		It("should fail with invalid db connection", func() {
			dbCharacter, err := repo.Create(character)
			Expect(err).To(Succeed())

			DB.Exec(`DROP SCHEMA public CASCADE;`)

			newName := f.NewCharacter().Name
			dbCharacter.Name = newName
			dbCharacter, err = repo.Save(character)
			Expect(err).NotTo(Succeed())
			Expect(dbCharacter).To(BeNil())
			DB.Exec(`CREATE SCHEMA public;`)
		})

		It("should work with valid DB connection", func() {
			dbCharacter, err := repo.Create(character)
			Expect(err).To(Succeed())

			newName := f.NewCharacter().Name
			dbCharacter.Name = newName
			dbCharacter, err = repo.Save(character)
			Expect(err).To(Succeed())
			Expect(dbCharacter.Name).To(Equal(newName))
		})
	})

	It("should be able to find by id", func() {
		dbCharacter, err := repo.Create(character)
		Expect(err).To(Succeed())

		foundCharacter, err := repo.FindById(uint64(dbCharacter.ID))
		Expect(err).To(Succeed())

		Expect(foundCharacter.ID).To(Equal(dbCharacter.ID))
		Expect(foundCharacter.Name).To(Equal(dbCharacter.Name))
	})

	It("should be able to delete by id", func() {
		dbCharacter, err := repo.Create(character)
		Expect(err).To(Succeed())

		err = repo.Delete(&model.Character{Model: gorm.Model{ID: dbCharacter.ID}})
		Expect(err).To(Succeed())

		foundCharacter, err := repo.FindById(uint64(dbCharacter.ID))
		Expect(err).NotTo(Succeed())
		Expect(foundCharacter).To(BeNil())
	})

	It("should be able to find all characters", func() {
		characters, err := repo.FindAll()
		Expect(err).To(Succeed())
		Expect(characters).To(BeEmpty())

		count := f.Factory().UintRange(3, 10)
		for i := uint(0); i < count; i++ {
			repo.Create(f.NewCharacter())
		}

		characters, err = repo.FindAll()
		Expect(err).To(Succeed())
		Expect(uint(len(characters))).To(Equal(count))
	})

	It("should be able to find all characters", func() {
		characters, err := repo.FindAll()
		Expect(err).To(Succeed())

		total := len(characters)
		Expect(total).To(Equal(0))

		characterId1 := uint(100)
		characterId2 := uint(101)

		count := f.Factory().UintRange(3, 10)
		total += int(count)
		for i := uint(0); i < count; i++ {
			character := f.NewCharacter()
			character.OwnerId = uint64(characterId1)
			repo.Create(character)
		}

		characters, err = repo.FindAllByOwner(uint64(characterId1))
		Expect(err).To(Succeed())
		Expect(len(characters)).To(Equal(int(count)))

		characters, err = repo.FindAll()
		Expect(err).To(Succeed())
		Expect(len(characters)).To(Equal(total))

		count = f.Factory().UintRange(3, 10)
		total += int(count)
		for i := uint(0); i < count; i++ {
			character := f.NewCharacter()
			character.OwnerId = uint64(characterId2)
			repo.Create(character)
		}

		characters, err = repo.FindAllByOwner(uint64(characterId2))
		Expect(err).To(Succeed())
		Expect(uint(len(characters))).To(Equal(count))

		characters, err = repo.FindAll()
		Expect(err).To(Succeed())
		Expect(len(characters)).To(Equal(total))
	})

	It("should be able to use a different trx", func() {
		Expect(repo.WithTrx(nil)).NotTo(BeNil())
		_, err := repo.FindAll()
		Expect(err).To(Succeed())

		repo.WithTrx(&gorm.DB{})
		_, err = repo.FindAll()
		Expect(err).To(Succeed())
	})
})

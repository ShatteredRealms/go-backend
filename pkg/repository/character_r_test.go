package repository_test

import (
	"context"
	"fmt"
	"time"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/ShatteredRealms/go-backend/pkg/model"
)

var _ = Describe("Character repository", func() {
	var (
		count = 0

		createCharacter = func() *model.Character {
			count++
			character := &model.Character{
				ID:        0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: 0,
				OwnerId:   faker.Username(),
				Name:      fmt.Sprintf("ownerid%d", count),
				Gender:    "Male",
				Realm:     "Human",
				PlayTime:  100,
				Location: model.Location{
					World: faker.Username(),
					X:     1.1,
					Y:     1.2,
					Z:     1.3,
					Roll:  1.4,
					Pitch: 1.5,
					Yaw:   1.6,
				},
			}

			outCharacter, err := characterRepo.Create(nil, character)
			Expect(err).To(BeNil())
			Expect(outCharacter).NotTo(BeNil())
			Expect(outCharacter.Validate()).To(Succeed())
			Expect(outCharacter.ID).To(BeEquivalentTo(character.ID))
			character = outCharacter
			hook.Reset()

			return character
		}
	)

	Describe("FindByName", func() {
		When("invalid input is given", func() {
			It("should not error", func() {
				ctx := context.Background()
				out, err := characterRepo.FindByName(ctx, "\";\n")
				Expect(out).To(BeNil())
				Expect(err).To(BeNil())

				hook.Reset()
				out, err = characterRepo.FindByName(nil, "")
				Expect(out).To(BeNil())
				Expect(err).To(BeNil())
			})
		})

		When("valid input", func() {
			It("should return nil if not found", func() {
				ctx := context.Background()
				out, err := characterRepo.FindByName(ctx, "name")
				Expect(out).To(BeNil())
				Expect(err).To(BeNil())
			})

			It("should return character if found", func() {
				character := createCharacter()
				findCharacter, err := characterRepo.FindByName(nil, character.Name)
				Expect(err).To(BeNil())
				Expect(findCharacter).NotTo(BeNil())
				Expect(findCharacter.Name).To(Equal(character.Name))
			})
		})
	})

	Describe("Create", func() {
		When("no conflicts exist", func() {
			It("should work", func() {
				newCharacter := &model.Character{
					ID:        0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					DeletedAt: 0,
					OwnerId:   faker.Username(),
					Name:      faker.Username(),
					Gender:    "Male",
					Realm:     "Human",
					PlayTime:  100,
					Location: model.Location{
						World: faker.Username(),
						X:     1.1,
						Y:     1.2,
						Z:     1.3,
						Roll:  1.4,
						Pitch: 1.5,
						Yaw:   1.6,
					},
				}
				outCharacter, err := characterRepo.Create(nil, newCharacter)
				Expect(err).To(BeNil())
				Expect(outCharacter).NotTo(BeNil())
				Expect(outCharacter.Validate()).To(Succeed())
			})
		})

		When("conflicts exist", func() {
			It("should throw errors with non-unique name", func() {
				character := createCharacter()
				outCharacter, err := characterRepo.Create(nil, character)
				Expect(err).NotTo(BeNil())
				Expect(outCharacter).To(BeNil())
			})
		})
	})

	Describe("Save", func() {
		When("no conflicts exist", func() {
			It("should work", func() {
				character := createCharacter()
				character.Name += faker.Currency()
				outCharacter, err := characterRepo.Save(nil, character)
				Expect(err).To(BeNil())
				Expect(outCharacter).NotTo(BeNil())
				Expect(outCharacter.Validate()).To(Succeed())
				Expect(outCharacter.Name).To(Equal(character.Name))
			})
		})

		When("conflicts exist", func() {
			It("should throw errors with non-unique name", func() {
				character := createCharacter()
				character.Name += "a"
				outCharacter, err := characterRepo.Create(nil, character)
				Expect(err).NotTo(HaveOccurred())
				Expect(outCharacter).NotTo(BeNil())
				Expect(outCharacter.Name).To(Equal(character.Name))

				character.Name = character.Name[:len(character.Name)-1]
				outCharacter, err = characterRepo.Save(nil, character)
				Expect(err).To(HaveOccurred())
				Expect(outCharacter).To(BeNil())
			})
		})
	})

	Describe("Delete", func() {
		When("valid input is given", func() {
			checkDelete := (func(ctx context.Context) {
				character := createCharacter()
				character.Name += "a"
				outCharacter, err := characterRepo.Create(nil, character)
				Expect(err).NotTo(HaveOccurred())
				Expect(outCharacter).NotTo(BeNil())
				Expect(outCharacter.Name).To(Equal(character.Name))

				out, err := characterRepo.FindById(context.Background(), outCharacter.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())

				Expect(characterRepo.Delete(ctx, character)).To(Succeed())
				out, err = characterRepo.FindById(context.Background(), outCharacter.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should succeed with invalid context", func() {
				checkDelete(nil)
			})

			It("should succeed with valid context", func() {
				checkDelete(context.Background())
			})
		})

		Context("invalid input", func() {
			It("should error", func() {
				Expect(characterRepo.Delete(nil, nil)).NotTo(Succeed())
				Expect(characterRepo.Delete(context.Background(), nil)).NotTo(Succeed())
			})
		})
	})

	Describe("FindById", func() {
		When("invalid input is given", func() {
			It("should not error", func() {
				out, err := characterRepo.FindById(nil, 0)
				Expect(out).To(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})

			It("should error if invalid id", func() {
				ctx := context.Background()
				out, err := characterRepo.FindById(ctx, 1e19)
				Expect(out).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})

		When("valid input is given", func() {
			It("should return nil if not found", func() {
				ctx := context.Background()
				out, err := characterRepo.FindById(ctx, 0)
				Expect(out).To(BeNil())
				Expect(err).To(BeNil())
			})

			It("should return character if found", func() {
				character := createCharacter()
				findCharacter, err := characterRepo.FindById(nil, character.ID)
				Expect(err).To(BeNil())
				Expect(findCharacter).NotTo(BeNil())
				Expect(findCharacter.Name).To(Equal(character.Name))
			})
		})
	})

	Describe("FindAll", func() {
		findAll := (func(ctx context.Context) {
			all, err := characterRepo.FindAll(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(all) >= 1).To(BeTrue())
		})
		It("should work with invalid context", func() {
			findAll(nil)
		})

		It("should work with valid context", func() {
			findAll(context.Background())
		})
	})

	Describe("FindAllByOwner", func() {
		When("a match exists", func() {
			findAllByOwner := (func(ctx context.Context) {
				newCharacter := &model.Character{
					ID:        0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					DeletedAt: 0,
					OwnerId:   faker.Username(),
					Name:      faker.Username() + "a",
					Gender:    "Male",
					Realm:     "Human",
					PlayTime:  100,
					Location: model.Location{
						World: faker.Username(),
						X:     1.1,
						Y:     1.2,
						Z:     1.3,
						Roll:  1.4,
						Pitch: 1.5,
						Yaw:   1.6,
					},
				}

				outCharacter, err := characterRepo.Create(nil, newCharacter)
				Expect(err).To(BeNil())
				Expect(outCharacter).NotTo(BeNil())
				Expect(outCharacter.Validate()).To(Succeed())
				all, err := characterRepo.FindAllByOwner(ctx, newCharacter.OwnerId)
				Expect(err).NotTo(HaveOccurred())
				Expect(all).To(HaveLen(1))
			})
			It("should work with invalid context", func() {
				findAllByOwner(nil)
			})

			It("should work with valid context", func() {
				findAllByOwner(context.Background())
			})
		})

		When("a match doesn't exists", func() {
			findAllByOwner := (func(ctx context.Context, owner string) {
				all, err := characterRepo.FindAllByOwner(ctx, owner)
				Expect(err).NotTo(HaveOccurred())
				Expect(all).To(HaveLen(0))
			})
			It("should work with invalid context", func() {
				findAllByOwner(nil, "")
			})

			It("should work with valid context", func() {
				findAllByOwner(context.Background(), "")
			})
		})
	})

	Describe("WithTrx", func() {
		When("valid trx is provided", func() {
			It("should pass input trx", func() {
				newTrx := &gorm.DB{}
				characterRepo.WithTrx(newTrx)
				_, err := characterRepo.FindAll(nil)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		When("invalid trx is provided", func() {
			It("should use original input trx", func() {
				characterRepo.WithTrx(nil)
				_, err := characterRepo.FindAll(nil)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	AfterEach(func() {
	})
})

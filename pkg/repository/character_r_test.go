package repository_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"gorm.io/gorm"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
)

var _ = Describe("Character repository", func() {
	var (
		hook *test.Hook

		db          *gorm.DB
		dbCloseFunc func()

		repo repository.CharacterRepository

		character *model.Character
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()

		db, dbCloseFunc = testdb.SetupGormWithDocker()
		Expect(db).NotTo(BeNil())
		repo = repository.NewCharacterRepository(db)
		Expect(repo).NotTo(BeNil())
		Expect(repo.Migrate(context.Background())).To(Succeed())

		character = &model.Character{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: 0,
			OwnerId:   "ownerid",
			Name:      "unreal",
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

		outCharacter, err := repo.Create(nil, character)
		Expect(err).To(BeNil())
		Expect(outCharacter).NotTo(BeNil())
		Expect(outCharacter.Validate()).To(Succeed())

		hook.Reset()
	})

	Describe("FindByName", func() {
		Context("invalid input", func() {
			It("should not error", func() {
				ctx := context.Background()
				out, err := repo.FindByName(ctx, "\";\n")
				Expect(out).To(BeNil())
				Expect(err).To(BeNil())

				hook.Reset()
				out, err = repo.FindByName(nil, "")
				Expect(out).To(BeNil())
				Expect(err).To(BeNil())
			})
		})

		Context("valid input", func() {
			It("should return nil if not found", func() {
				ctx := context.Background()
				out, err := repo.FindByName(ctx, "name")
				Expect(out).To(BeNil())
				Expect(err).To(BeNil())
			})

			It("should return character if found", func() {
				findCharacter, err := repo.FindByName(nil, character.Name)
				Expect(err).To(BeNil())
				Expect(findCharacter).NotTo(BeNil())
				Expect(findCharacter.Name).To(Equal(character.Name))
			})
		})
	})

	Describe("Create", func() {
		Context("no conflicts", func() {
			It("should work", func() {
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
				outCharacter, err := repo.Create(nil, newCharacter)
				Expect(err).To(BeNil())
				Expect(outCharacter).NotTo(BeNil())
				Expect(outCharacter.Validate()).To(Succeed())
			})
		})

		Context("conflicts", func() {
			It("should throw errors with non-unique name", func() {
				outCharacter, err := repo.Create(nil, character)
				Expect(err).NotTo(BeNil())
				Expect(outCharacter).To(BeNil())
			})
		})
	})

	Describe("Save", func() {
		Context("no conflicts", func() {
			It("should work", func() {
				character.Name += "a"
				outCharacter, err := repo.Save(nil, character)
				Expect(err).To(BeNil())
				Expect(outCharacter).NotTo(BeNil())
				Expect(outCharacter.Validate()).To(Succeed())
				Expect(outCharacter.Name).To(Equal(character.Name))
			})
		})

		Context("conflicts", func() {
			It("should throw errors with non-unique name", func() {
				character.Name += "a"
				outCharacter, err := repo.Create(nil, character)
				Expect(err).NotTo(HaveOccurred())
				Expect(outCharacter).NotTo(BeNil())
				Expect(outCharacter.Name).To(Equal(character.Name))

				character.Name = character.Name[:len(character.Name)-1]
				outCharacter, err = repo.Save(nil, character)
				Expect(err).To(HaveOccurred())
				Expect(outCharacter).To(BeNil())
			})
		})
	})

	Describe("Delete", func() {
		Context("valid input", func() {
			checkDelete := (func(ctx context.Context) {
				all, err := repo.FindAll(context.Background())
				Expect(err).NotTo(HaveOccurred())
				Expect(all).To(HaveLen(1))

				Expect(repo.Delete(ctx, character)).To(Succeed())
				all, err = repo.FindAll(context.Background())
				Expect(err).NotTo(HaveOccurred())
				Expect(all).To(HaveLen(0))
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
				Expect(repo.Delete(nil, nil)).NotTo(Succeed())
				Expect(repo.Delete(context.Background(), nil)).NotTo(Succeed())
			})
		})
	})

	Describe("FindById", func() {
		Context("invalid input", func() {
			It("should not error", func() {
				out, err := repo.FindById(nil, 0)
				Expect(out).To(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})

			It("should error if invalid id", func() {
				ctx := context.Background()
				out, err := repo.FindById(ctx, 1e19)
				Expect(out).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})

		Context("valid input", func() {
			It("should return nil if not found", func() {
				ctx := context.Background()
				out, err := repo.FindById(ctx, 0)
				Expect(out).To(BeNil())
				Expect(err).To(BeNil())
			})

			It("should return character if found", func() {
				findCharacter, err := repo.FindById(nil, character.ID)
				Expect(err).To(BeNil())
				Expect(findCharacter).NotTo(BeNil())
				Expect(findCharacter.Name).To(Equal(character.Name))
			})
		})
	})

	Describe("FindById", func() {
		findAll := (func(ctx context.Context) {
			all, err := repo.FindAll(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(all).To(HaveLen(1))
		})
		It("should work with invalid context", func() {
			findAll(nil)
		})

		It("should work with valid context", func() {
			findAll(context.Background())
		})
	})

	Describe("FindAllByOwner", func() {
		Context("valid input", func() {
			Context("match exists", func() {
				findAllByOwner := (func(ctx context.Context, owner string) {
					all, err := repo.FindAllByOwner(ctx, owner)
					Expect(err).NotTo(HaveOccurred())
					Expect(all).To(HaveLen(1))
				})
				It("should work with invalid context", func() {
					findAllByOwner(nil, character.OwnerId)
				})

				It("should work with valid context", func() {
					findAllByOwner(context.Background(), character.OwnerId)
				})
			})

			Context("match doesn't exists", func() {
				findAllByOwner := (func(ctx context.Context, owner string) {
					all, err := repo.FindAllByOwner(ctx, owner)
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
	})

	Describe("WithTrx", func() {
		Context("valid trx", func() {
			It("should pass input trx", func() {
				newTrx := &gorm.DB{}
				repo.WithTrx(newTrx)
				chars, err := repo.FindAll(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(chars).To(HaveLen(1))
			})
		})
		Context("invalid trx", func() {
			It("should use original input trx", func() {
				repo.WithTrx(nil)
				_, err := repo.FindAll(nil)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	AfterEach(func() {
		dbCloseFunc()
	})
})

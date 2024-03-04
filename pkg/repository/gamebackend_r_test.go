package repository_test

import (
	"context"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"gorm.io/gorm"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
)

var _ = Describe("Gamebackend repository", func() {
	var (
		hook *test.Hook

		db          *gorm.DB
		dbCloseFunc func()

		repo repository.GamebackendRepository

		m          *model.Map
		dimension  *model.Dimension
		pendingCon *model.PendingConnection
		err        error
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()

		db, dbCloseFunc = testdb.SetupGormWithDocker()
		Expect(db).NotTo(BeNil())

		repo = repository.NewGamebackendRepository(db)
		Expect(repo).NotTo(BeNil())
		Expect(repo.Migrate(context.Background())).To(Succeed())

		m, err = repo.CreateMap(nil, faker.Username(), faker.Username(), 4, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(m).NotTo(BeNil())
		Expect(m.Id).NotTo(BeNil())

		dimension, err = repo.CreateDimension(nil, faker.Username(), faker.Username(), faker.Username(), []*uuid.UUID{m.Id})
		Expect(err).NotTo(HaveOccurred())
		Expect(dimension).NotTo(BeNil())
		Expect(dimension.Id).NotTo(BeNil())

		pendingCon, err = repo.CreatePendingConnection(nil, faker.Username(), faker.Username())
		Expect(err).NotTo(HaveOccurred())
		Expect(pendingCon).NotTo(BeNil())
		Expect(pendingCon.Id).NotTo(BeNil())

		hook.Reset()
	})

	Describe("CreatePendingConnection", func() {
		Context("valid input", func() {
			It("should work", func() {
				out, err := repo.CreatePendingConnection(nil, faker.Username(), faker.Username())
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).NotTo(BeNil())
			})
		})

		Context("invalid input", func() {
			It("should require valid character", func() {
				out, err := repo.CreatePendingConnection(nil, "", faker.Username())
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("CreateDimension", func() {
		Context("valid input", func() {
			It("should work", func() {
				out, err := repo.CreateDimension(nil, dimension.Name+"a", faker.Username(), faker.Username(), []*uuid.UUID{})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).NotTo(BeNil())
			})

			It("should allow duplicate if name was deleted", func() {
				Expect(repo.DeleteDimensionById(nil, dimension.Id)).To(Succeed())
				out, err := repo.CreateDimension(nil, dimension.Name, faker.Username(), faker.Username(), []*uuid.UUID{})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).NotTo(BeNil())
			})
		})

		Context("invalid input", func() {
			// @TODO: Figure unique composite index with name and deleted
			// It("should not allow duplicate names", func() {
			// 	out, err := repo.CreateDimension(nil, dimension.Name, faker.Username(), faker.Username(), []*uuid.UUID{})
			// 	Expect(err).To(HaveOccurred())
			// 	Expect(out).To(BeNil())
			// })
		})
	})

	Describe("CreateMap", func() {
		Context("valid input", func() {
			It("should work", func() {
				out, err := repo.CreateMap(nil, m.Name+"a", faker.Username(), 4, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).NotTo(BeNil())
			})
		})

		Context("invalid input", func() {
			It("should error on duplicate name", func() {
				out, err := repo.CreateMap(nil, m.Name, faker.Username(), 4, false)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("DeletePendingConnection", func() {
		Context("valid input", func() {
			It("should work", func() {
				Expect(repo.FindPendingConnection(nil, pendingCon.Id)).NotTo(BeNil())
				Expect(repo.DeletePendingConnection(nil, pendingCon.Id)).To(Succeed())
				Expect(repo.FindPendingConnection(nil, pendingCon.Id)).To(BeNil())
			})
		})

		Context("invalid input", func() {
			It("should error on nil id", func() {
				Expect(repo.FindPendingConnection(nil, pendingCon.Id)).NotTo(BeNil())
				Expect(repo.DeletePendingConnection(nil, nil)).NotTo(Succeed())
				Expect(repo.FindPendingConnection(nil, pendingCon.Id)).NotTo(BeNil())
			})
		})
	})

	Describe("FindPendingConnection", func() {
		Context("valid input", func() {
			It("should work", func() {
				Expect(repo.FindPendingConnection(nil, pendingCon.Id)).NotTo(BeNil())
			})
		})

		Context("invalid input", func() {
			It("should return no match", func() {
				Expect(repo.FindPendingConnection(nil, nil)).To(BeNil())
			})
		})
	})

	Describe("DeleteDimensionById", func() {
		Context("valid input", func() {
			It("should work", func() {
				Expect(repo.DeleteDimensionById(nil, dimension.Id)).To(Succeed())
				Expect(repo.DeleteDimensionById(nil, dimension.Id)).To(Succeed())
			})
		})
	})

	Describe("DeleteDimensionByName", func() {
		Context("valid input", func() {
			It("should work", func() {
				Expect(repo.DeleteDimensionByName(nil, dimension.Name)).To(Succeed())
				Expect(repo.DeleteDimensionByName(nil, dimension.Name)).To(Succeed())
			})
		})
	})

	Describe("DeleteMapById", func() {
		Context("valid input", func() {
			It("should work", func() {
				Expect(repo.DeleteMapById(nil, m.Id)).To(Succeed())
				Expect(repo.DeleteMapById(nil, m.Id)).To(Succeed())
			})
		})
	})

	Describe("DeleteMapByName", func() {
		Context("valid input", func() {
			It("should work", func() {
				Expect(repo.DeleteMapByName(nil, m.Name)).To(Succeed())
				Expect(repo.DeleteMapByName(nil, m.Name)).To(Succeed())
			})
		})
	})

	Describe("DuplicateDimesnion", func() {
		Context("valid input", func() {
			It("should work", func() {
				out, err := repo.DuplicateDimension(nil, dimension.Id, dimension.Name+"a")
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).NotTo(BeNil())
				Expect(out.Name).To(Equal(dimension.Name + "a"))
				Expect(repo.FindAllDimensions(nil)).To(HaveLen(2))
			})
		})

		Context("invalid input", func() {
			It("should error when given dimension doesn't exist", func() {
				id, err := uuid.NewRandom()
				Expect(err).NotTo(HaveOccurred())

				out, err := repo.DuplicateDimension(nil, &id, dimension.Name+"a")
				Expect(err).To(MatchError(model.ErrDoesNotExist))
				Expect(out).To(BeNil())
				Expect(repo.FindAllDimensions(nil)).To(HaveLen(1))
			})
		})
	})

	Describe("SaveDimension", func() {
		Context("valid input", func() {
			It("should work", func() {
				Expect(dimension.Maps).To(HaveLen(1))
				dimension.Name += "a"
				out, err := repo.SaveDimension(nil, dimension)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).To(Equal(dimension.Id))
				Expect(out.Maps).To(HaveLen(1))
			})
		})

		Context("invalid input", func() {
			It("", func() {
			})
		})
	})

	AfterEach(func() {
		dbCloseFunc()
	})
})

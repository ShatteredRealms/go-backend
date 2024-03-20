package repository_test

import (
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/ShatteredRealms/go-backend/pkg/model"
)

var _ = Describe("Gamebackend repository", func() {
	var (
		createModels = func() (*model.Map, *model.Dimension, *model.PendingConnection) {
			m, err := gamebackendRepo.CreateMap(nil, faker.Username(), faker.Username(), 4, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(m).NotTo(BeNil())
			Expect(m.Id).NotTo(BeNil())

			dimension, err := gamebackendRepo.CreateDimension(nil, faker.Username(), faker.Username(), faker.Username(), []*uuid.UUID{m.Id})
			Expect(err).NotTo(HaveOccurred())
			Expect(dimension).NotTo(BeNil())
			Expect(dimension.Id).NotTo(BeNil())

			pendingCon, err := gamebackendRepo.CreatePendingConnection(nil, faker.Username(), faker.Username())
			Expect(err).NotTo(HaveOccurred())
			Expect(pendingCon).NotTo(BeNil())
			Expect(pendingCon.Id).NotTo(BeNil())

			return m, dimension, pendingCon
		}
	)

	Describe("CreatePendingConnection", func() {
		When("given valid input", func() {
			It("should work", func() {
				out, err := gamebackendRepo.CreatePendingConnection(nil, faker.Username(), faker.Username())
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should require valid character", func() {
				out, err := gamebackendRepo.CreatePendingConnection(nil, "", faker.Username())
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("CreateDimension", func() {
		When("given valid input", func() {
			It("should work", func() {
				_, dimension, _ := createModels()
				out, err := gamebackendRepo.CreateDimension(nil, dimension.Name+"a", faker.Username(), faker.Username(), []*uuid.UUID{})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).NotTo(BeNil())
			})

			It("should allow duplicate if name was deleted", func() {
				_, dimension, _ := createModels()
				Expect(gamebackendRepo.DeleteDimensionById(nil, dimension.Id)).To(Succeed())
				out, err := gamebackendRepo.CreateDimension(nil, dimension.Name, faker.Username(), faker.Username(), []*uuid.UUID{})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			// @TODO: Figure unique composite index with name and deleted
			// It("should not allow duplicate names", func() {
			// 	out, err := repo.CreateDimension(nil, dimension.Name, faker.Username(), faker.Username(), []*uuid.UUID{})
			// 	Expect(err).To(HaveOccurred())
			// 	Expect(out).To(BeNil())
			// })
		})
	})

	Describe("CreateMap", func() {
		When("given valid input", func() {
			It("should work", func() {
				m, _, _ := createModels()
				out, err := gamebackendRepo.CreateMap(nil, m.Name+"a", faker.Username(), 4, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error on duplicate name", func() {
				m, _, _ := createModels()
				out, err := gamebackendRepo.CreateMap(nil, m.Name, faker.Username(), 4, false)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("DeletePendingConnection", func() {
		When("given valid input", func() {
			It("should work", func() {
				_, _, pendingCon := createModels()
				Expect(gamebackendRepo.FindPendingConnection(nil, pendingCon.Id)).NotTo(BeNil())
				Expect(gamebackendRepo.DeletePendingConnection(nil, pendingCon.Id)).To(Succeed())
				Expect(gamebackendRepo.FindPendingConnection(nil, pendingCon.Id)).To(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error on nil id", func() {
				_, _, pendingCon := createModels()
				Expect(gamebackendRepo.FindPendingConnection(nil, pendingCon.Id)).NotTo(BeNil())
				Expect(gamebackendRepo.DeletePendingConnection(nil, nil)).NotTo(Succeed())
				Expect(gamebackendRepo.FindPendingConnection(nil, pendingCon.Id)).NotTo(BeNil())
			})
		})
	})

	Describe("FindPendingConnection", func() {
		When("given valid input", func() {
			It("should work", func() {
				_, _, pendingCon := createModels()
				Expect(gamebackendRepo.FindPendingConnection(nil, pendingCon.Id)).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should return no match", func() {
				Expect(gamebackendRepo.FindPendingConnection(nil, nil)).To(BeNil())
			})
		})
	})

	Describe("DeleteDimensionById", func() {
		When("given valid input", func() {
			It("should work", func() {
				_, dimension, _ := createModels()
				Expect(gamebackendRepo.DeleteDimensionById(nil, dimension.Id)).To(Succeed())
				Expect(gamebackendRepo.DeleteDimensionById(nil, dimension.Id)).To(Succeed())
			})
		})
	})

	Describe("DeleteDimensionByName", func() {
		When("given valid input", func() {
			It("should work", func() {
				_, dimension, _ := createModels()
				Expect(gamebackendRepo.DeleteDimensionByName(nil, dimension.Name)).To(Succeed())
				Expect(gamebackendRepo.DeleteDimensionByName(nil, dimension.Name)).To(Succeed())
			})
		})
	})

	Describe("DeleteMapById", func() {
		When("given valid input", func() {
			It("should work", func() {
				m, _, _ := createModels()
				Expect(gamebackendRepo.DeleteMapById(nil, m.Id)).To(Succeed())
				Expect(gamebackendRepo.DeleteMapById(nil, m.Id)).To(Succeed())
			})
		})
	})

	Describe("DeleteMapByName", func() {
		When("given valid input", func() {
			It("should work", func() {
				m, _, _ := createModels()
				Expect(gamebackendRepo.DeleteMapByName(nil, m.Name)).To(Succeed())
				Expect(gamebackendRepo.DeleteMapByName(nil, m.Name)).To(Succeed())
			})
		})
	})

	Describe("SaveDimension", func() {
		When("given valid input", func() {
			It("should work", func() {
				_, dimension, _ := createModels()
				Expect(dimension.Maps).To(HaveLen(1))
				dimension.Name += "a"
				out, err := gamebackendRepo.SaveDimension(nil, dimension)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).To(Equal(dimension.Id))
				Expect(out.Maps).To(HaveLen(1))
			})
		})

		When("given invalid input", func() {
			It("should error with nil dimension", func() {
				out, err := gamebackendRepo.SaveDimension(nil, nil)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("SaveMap", func() {
		When("given valid input", func() {
			It("should work", func() {
				m, _, _ := createModels()

				m.Name += "a"
				out, err := gamebackendRepo.SaveMap(nil, m)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).To(Equal(m.Id))
				Expect(out.Name).To(Equal(m.Name))
			})
		})

		When("given invalid input", func() {
			It("should error with nil m", func() {
				out, err := gamebackendRepo.SaveMap(nil, nil)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("FindAllDimensions", func() {
		When("given valid input", func() {
			It("should work", func() {
				createModels()
				out, err := gamebackendRepo.FindAllDimensions(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(len(out) >= 1).To(BeTrue())
			})
		})
	})

	Describe("FindAllMaps", func() {
		When("given valid input", func() {
			It("should work", func() {
				createModels()
				out, err := gamebackendRepo.FindAllMaps(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(len(out) >= 1).To(BeTrue())
			})
		})
	})

	Describe("FindDimensionByName", func() {
		When("given valid input", func() {
			It("should work", func() {
				_, dimension, _ := createModels()
				out, err := gamebackendRepo.FindDimensionByName(nil, dimension.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).To(Equal(dimension.Id))
				Expect(out.Name).To(Equal(dimension.Name))
			})
		})

		When("given invalid input", func() {
			It("should return nil for no match", func() {
				_, dimension, _ := createModels()
				out, err := gamebackendRepo.FindDimensionByName(nil, dimension.Name+"a")
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(BeNil())

				out, err = gamebackendRepo.FindDimensionByName(nil, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("FindMapByName", func() {
		When("given valid input", func() {
			It("should work", func() {
				m, _, _ := createModels()
				out, err := gamebackendRepo.FindMapByName(nil, m.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).To(Equal(m.Id))
				Expect(out.Name).To(Equal(m.Name))
			})
		})

		When("given invalid input", func() {
			It("should return nil for no match", func() {
				m, _, _ := createModels()
				out, err := gamebackendRepo.FindMapByName(nil, m.Name+"a")
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(BeNil())

				out, err = gamebackendRepo.FindMapByName(nil, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("FindMapById", func() {
		When("given valid input", func() {
			It("should work", func() {
				m, _, _ := createModels()
				out, err := gamebackendRepo.FindMapById(nil, m.Id)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.Id).To(Equal(m.Id))
			})
		})

		When("given invalid input", func() {
			It("should return nil for no match", func() {
				id, err := uuid.NewRandom()
				Expect(err).NotTo(HaveOccurred())
				out, err := gamebackendRepo.FindMapById(nil, &id)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should return err for nil id", func() {
				out, err := gamebackendRepo.FindMapById(nil, nil)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("FindDimensionsByIds", func() {
		When("given valid input", func() {
			It("should work", func() {
				_, dimension, _ := createModels()

				id, err := uuid.NewRandom()
				Expect(err).NotTo(HaveOccurred())
				out, err := gamebackendRepo.FindDimensionsByIds(nil, []*uuid.UUID{dimension.Id, &id})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(1))

				out, err = gamebackendRepo.FindDimensionsByIds(nil, []*uuid.UUID{&id})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(0))

				out, err = gamebackendRepo.FindDimensionsByIds(nil, []*uuid.UUID{})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(0))
			})
		})

		When("given invalid input", func() {

		})
	})

	Describe("FindDimensionsByNames", func() {
		When("given valid input", func() {
			It("should work", func() {
				_, dimension, _ := createModels()

				name := faker.Username() + "a"
				out, err := gamebackendRepo.FindDimensionsByNames(nil, []string{dimension.Name, name})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(1))

				out, err = gamebackendRepo.FindDimensionsByNames(nil, []string{name})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(0))

				out, err = gamebackendRepo.FindDimensionsByNames(nil, []string{})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(0))
			})
		})

		When("given invalid input", func() {

		})
	})

	Describe("FindMapsByNames", func() {
		When("given valid input", func() {
			It("should work", func() {
				m, _, _ := createModels()
				name := faker.Username() + "a"
				out, err := gamebackendRepo.FindMapsByNames(nil, []string{m.Name, name})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(1))

				out, err = gamebackendRepo.FindMapsByNames(nil, []string{name})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(0))

				out, err = gamebackendRepo.FindMapsByNames(nil, []string{})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(0))
			})
		})

		When("given invalid input", func() {

		})
	})

	Describe("FindDimensionsWithMapIds", func() {
		When("given valid input", func() {
			It("should work", func() {
				m, _, _ := createModels()

				id, err := uuid.NewRandom()
				Expect(err).NotTo(HaveOccurred())
				out, err := gamebackendRepo.FindDimensionsWithMapIds(nil, []*uuid.UUID{m.Id, &id})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(1))

				out, err = gamebackendRepo.FindDimensionsWithMapIds(nil, []*uuid.UUID{&id})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(0))

				out, err = gamebackendRepo.FindDimensionsWithMapIds(nil, []*uuid.UUID{})
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out).To(HaveLen(0))
			})
		})

		When("given invalid input", func() {

		})
	})

	Describe("WithTrx", func() {
		When("given valid input", func() {
			It("should work", func() {
				Expect(gamebackendRepo.WithTrx(&gorm.DB{})).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should use original trx", func() {
				Expect(gamebackendRepo.WithTrx(nil)).NotTo(BeNil())
			})
		})
	})
})

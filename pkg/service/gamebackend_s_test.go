package service_test

import (
	"context"
	"fmt"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/service"
)

var _ = Describe("Gamebackend service", func() {
	var (
		hook           *test.Hook
		mockController *gomock.Controller
		mockRepository *mocks.MockGamebackendRepository

		gbService service.GamebackendService

		character *model.Character
		dimension *model.Dimension
		m         *model.Map
		ctx       context.Context
		fakeErr   error
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
		mockController = gomock.NewController(GinkgoT())
		mockRepository = mocks.NewMockGamebackendRepository(mockController)
		hook.Reset()

		var err error
		mockRepository.EXPECT().Migrate(ctx).Return(nil)
		gbService, err = service.NewGamebackendService(ctx, mockRepository)
		Expect(err).To(BeNil())
		Expect(gbService).NotTo(BeNil())
		hook.Reset()

		ctx = context.Background()
		fakeErr = fmt.Errorf("error: %s", faker.Username())
		character = &model.Character{
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

		uuid1, err := uuid.NewRandom()
		Expect(err).NotTo(HaveOccurred())
		uuid2, err := uuid.NewRandom()
		Expect(err).NotTo(HaveOccurred())
		m = &model.Map{
			Model: model.Model{
				Id:        &uuid2,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: gorm.DeletedAt{},
			},
			Name:       faker.Username(),
			Path:       faker.Username(),
			MaxPlayers: 40,
			Instanced:  false,
		}
		dimension = &model.Dimension{
			Model: model.Model{
				Id:        &uuid1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: gorm.DeletedAt{},
			},
			Name:     faker.Username(),
			Location: "us-central",
			Version:  faker.Username(),
			Maps:     []*model.Map{m},
		}
		m.Dimensions = []*model.Dimension{dimension}
	})

	Describe("NewGamebackendService", func() {
		It("should fail if migration fails", func() {
			mockRepository.EXPECT().Migrate(ctx).Return(fakeErr)
			gbService, err := service.NewGamebackendService(ctx, mockRepository)
			Expect(err).To(MatchError(fakeErr))
			Expect(gbService).To(BeNil())
		})
	})

	Describe("CreatePendingConnection", func() {
		It("should work", func() {
			pendingConnection := &model.PendingConnection{}
			Expect(faker.FakeData(pendingConnection)).To(Succeed())
			mockRepository.EXPECT().CreatePendingConnection(ctx, character.Name, pendingConnection.ServerName).Return(pendingConnection, fakeErr)
			out, err := gbService.CreatePendingConnection(ctx, character.Name, pendingConnection.ServerName)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(Equal(pendingConnection))
		})
	})

	Describe("CheckPlayerConnection", func() {
		When("given valid input", func() {
			It("should work", func() {
				pendingConnection := &model.PendingConnection{}
				Expect(faker.FakeData(pendingConnection)).To(Succeed())
				pendingConnection.CreatedAt = time.Now()

				mockRepository.EXPECT().FindPendingConnection(ctx, pendingConnection.Id).Return(pendingConnection)
				mockRepository.EXPECT().DeletePendingConnection(ctx, pendingConnection.Id)
				out, err := gbService.CheckPlayerConnection(ctx, pendingConnection.Id, pendingConnection.ServerName)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(Equal(pendingConnection))
			})
		})

		When("given invalid input", func() {
			It("should fail due to no pending connection found", func() {
				pendingConnection := &model.PendingConnection{}
				Expect(faker.FakeData(pendingConnection)).To(Succeed())
				pendingConnection.CreatedAt = time.Now()

				mockRepository.EXPECT().FindPendingConnection(ctx, pendingConnection.Id).Return(nil)
				out, err := gbService.CheckPlayerConnection(ctx, pendingConnection.Id, pendingConnection.ServerName)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should fail if the pending connection server names don't match", func() {
				pendingConnection := &model.PendingConnection{}
				Expect(faker.FakeData(pendingConnection)).To(Succeed())
				pendingConnection.CreatedAt = time.Now()

				mockRepository.EXPECT().FindPendingConnection(ctx, pendingConnection.Id).Return(pendingConnection)
				out, err := gbService.CheckPlayerConnection(ctx, pendingConnection.Id, pendingConnection.ServerName+"a")
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should fail if the connection expired", func() {
				pendingConnection := &model.PendingConnection{}
				Expect(faker.FakeData(pendingConnection)).To(Succeed())
				pendingConnection.CreatedAt = time.Now().Add(-time.Minute)

				mockRepository.EXPECT().FindPendingConnection(ctx, pendingConnection.Id).Return(pendingConnection)
				mockRepository.EXPECT().DeletePendingConnection(ctx, pendingConnection.Id)
				out, err := gbService.CheckPlayerConnection(ctx, pendingConnection.Id, pendingConnection.ServerName)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("CreateDimension", func() {
		It("should work", func() {
			mockRepository.EXPECT().CreateDimension(ctx, dimension.Name, dimension.Location, dimension.Version, []*uuid.UUID{}).Return(dimension, fakeErr)
			out, err := gbService.CreateDimension(ctx, dimension.Name, dimension.Location, dimension.Version, []*uuid.UUID{})
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(Equal(dimension))
		})
	})

	Describe("CreateMap", func() {
		It("should work", func() {
			mockRepository.EXPECT().CreateMap(ctx, m.Name, m.Path, m.MaxPlayers, m.Instanced).Return(m, fakeErr)
			out, err := gbService.CreateMap(ctx, m.Name, m.Path, m.MaxPlayers, m.Instanced)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(Equal(m))
		})
	})

	Describe("DeleteDimensionById", func() {
		It("should work", func() {
			mockRepository.EXPECT().DeleteDimensionById(ctx, dimension.Id).Return(fakeErr)
			err := gbService.DeleteDimensionById(ctx, dimension.Id)
			Expect(err).To(MatchError(fakeErr))
		})
	})

	Describe("DeleteMapById", func() {
		It("should work", func() {
			mockRepository.EXPECT().DeleteMapById(ctx, m.Id).Return(fakeErr)
			err := gbService.DeleteMapById(ctx, m.Id)
			Expect(err).To(MatchError(fakeErr))
		})
	})

	Describe("DeleteMapByName", func() {
		It("should work", func() {
			mockRepository.EXPECT().DeleteDimensionByName(ctx, dimension.Name).Return(fakeErr)
			err := gbService.DeleteDimensionByName(ctx, dimension.Name)
			Expect(err).To(MatchError(fakeErr))
		})
	})

	Describe("DeleteMapByName", func() {
		It("should work", func() {
			mockRepository.EXPECT().DeleteMapByName(ctx, m.Name).Return(fakeErr)
			err := gbService.DeleteMapByName(ctx, m.Name)
			Expect(err).To(MatchError(fakeErr))
		})
	})

	Describe("DuplicateDimension", func() {
		It("should work", func() {
			mockRepository.EXPECT().DuplicateDimension(ctx, dimension.Id, dimension.Name+"a").Return(dimension, fakeErr)
			out, err := gbService.DuplicateDimension(ctx, dimension.Id, dimension.Name+"a")
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(Equal(dimension))
		})
	})

	Describe("EditDimension", func() {
		It("should work", func() {

		})
	})

	Describe("EditMap", func() {
		It("should work", func() {

		})
	})

	Describe("FindAllDimensions", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindAllDimensions(ctx).Return(model.Dimensions{dimension}, fakeErr)
			out, err := gbService.FindAllDimensions(ctx)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(dimension))
		})
	})

	Describe("FindAllMap", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindAllMaps(ctx).Return(model.Maps{m}, fakeErr)
			out, err := gbService.FindAllMaps(ctx)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(m))
		})
	})

	Describe("FindDimensionById", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindDimensionById(ctx, dimension.Id).Return(dimension, fakeErr)
			out, err := gbService.FindDimensionById(ctx, dimension.Id)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(Equal(dimension))
		})
	})

	Describe("FindDimensionByName", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindDimensionByName(ctx, dimension.Name).Return(dimension, fakeErr)
			out, err := gbService.FindDimensionByName(ctx, dimension.Name)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(Equal(dimension))
		})
	})

	Describe("FindDimensionsByIds", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindDimensionsByIds(ctx, []*uuid.UUID{dimension.Id}).Return(model.Dimensions{dimension}, fakeErr)
			out, err := gbService.FindDimensionsByIds(ctx, []*uuid.UUID{dimension.Id})
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(dimension))
		})
	})

	Describe("FindDimensionsByNames", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindDimensionsByNames(ctx, []string{dimension.Name}).Return(model.Dimensions{dimension}, fakeErr)
			out, err := gbService.FindDimensionsByNames(ctx, []string{dimension.Name})
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(dimension))
		})
	})

	Describe("FindDimensionsWithmapIds", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindDimensionsWithMapIds(ctx, []*uuid.UUID{m.Id}).Return(model.Dimensions{dimension}, fakeErr)
			out, err := gbService.FindDimensionsWithMapIds(ctx, []*uuid.UUID{m.Id})
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(dimension))
		})
	})

	Describe("FindMapById", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindMapById(ctx, m.Id).Return(m, fakeErr)
			out, err := gbService.FindMapById(ctx, m.Id)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(Equal(m))
		})
	})

	Describe("FindMapByName", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindMapByName(ctx, m.Name).Return(m, fakeErr)
			out, err := gbService.FindMapByName(ctx, m.Name)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(Equal(m))
		})
	})

	Describe("FindMapsByIds", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindMapsByIds(ctx, []*uuid.UUID{m.Id}).Return(model.Maps{m}, fakeErr)
			out, err := gbService.FindMapsByIds(ctx, []*uuid.UUID{m.Id})
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(m))
		})
	})

	Describe("FindMapsByNames", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindMapsByNames(ctx, []string{m.Name}).Return(model.Maps{m}, fakeErr)
			out, err := gbService.FindMapsByNames(ctx, []string{m.Name})
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(m))
		})
	})
})

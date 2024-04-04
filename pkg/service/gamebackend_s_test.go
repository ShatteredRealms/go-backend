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
	"github.com/ShatteredRealms/go-backend/pkg/model/character"
	"github.com/ShatteredRealms/go-backend/pkg/model/game"
	"github.com/ShatteredRealms/go-backend/pkg/model/gamebackend"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/service"
)

var _ = Describe("Gamebackend service", func() {
	var (
		hook           *test.Hook
		mockController *gomock.Controller
		mockRepository *mocks.MockGamebackendRepository

		gbService service.GamebackendService

		char      *character.Character
		dimension *game.Dimension
		m         *game.Map
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
		char = &character.Character{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: 0,
			OwnerId:   faker.Username(),
			Name:      faker.Username(),
			Gender:    "Male",
			Realm:     "Human",
			PlayTime:  100,
			Location: game.Location{
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
		m = &game.Map{
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
		dimension = &game.Dimension{
			Model: model.Model{
				Id:        &uuid1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: gorm.DeletedAt{},
			},
			Name:     faker.Username(),
			Location: "us-central",
			Version:  faker.Username(),
			Maps:     []*game.Map{m},
		}
		m.Dimensions = []*game.Dimension{dimension}
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
			pendingConnection := &gamebackend.PendingConnection{}
			Expect(faker.FakeData(pendingConnection)).To(Succeed())
			mockRepository.EXPECT().CreatePendingConnection(ctx, char.Name, pendingConnection.ServerName).Return(pendingConnection, fakeErr)
			out, err := gbService.CreatePendingConnection(ctx, char.Name, pendingConnection.ServerName)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(Equal(pendingConnection))
		})
	})

	Describe("CheckPlayerConnection", func() {
		When("given valid input", func() {
			It("should work", func() {
				pendingConnection := &gamebackend.PendingConnection{}
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
				pendingConnection := &gamebackend.PendingConnection{}
				Expect(faker.FakeData(pendingConnection)).To(Succeed())
				pendingConnection.CreatedAt = time.Now()

				mockRepository.EXPECT().FindPendingConnection(ctx, pendingConnection.Id).Return(nil)
				out, err := gbService.CheckPlayerConnection(ctx, pendingConnection.Id, pendingConnection.ServerName)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should fail if the pending connection server names don't match", func() {
				pendingConnection := &gamebackend.PendingConnection{}
				Expect(faker.FakeData(pendingConnection)).To(Succeed())
				pendingConnection.CreatedAt = time.Now()

				mockRepository.EXPECT().FindPendingConnection(ctx, pendingConnection.Id).Return(pendingConnection)
				out, err := gbService.CheckPlayerConnection(ctx, pendingConnection.Id, pendingConnection.ServerName+"a")
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should fail if the connection expired", func() {
				pendingConnection := &gamebackend.PendingConnection{}
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
		var (
			idTarget   *pb.DimensionTarget
			nameTarget *pb.DimensionTarget
		)
		BeforeEach(func() {

			idTarget = &pb.DimensionTarget{
				FindBy: &pb.DimensionTarget_Id{
					Id: dimension.Id.String(),
				},
			}
			nameTarget = &pb.DimensionTarget{
				FindBy: &pb.DimensionTarget_Name{
					Name: dimension.Name,
				},
			}
		})

		When("given valid input", func() {
			It("should work for ids", func() {
				mockRepository.EXPECT().CreateDimension(ctx, dimension.Name+"a", dimension.Location, dimension.Version, gomock.Any()).Return(dimension, fakeErr)
				mockRepository.EXPECT().FindDimensionById(ctx, dimension.Id).Return(dimension, nil)
				out, err := gbService.DuplicateDimension(ctx, idTarget, dimension.Name+"a")
				Expect(err).To(MatchError(fakeErr))
				Expect(out).To(Equal(dimension))
			})

			It("should work for names", func() {
				mockRepository.EXPECT().CreateDimension(ctx, dimension.Name+"a", dimension.Location, dimension.Version, gomock.Any()).Return(dimension, fakeErr)
				mockRepository.EXPECT().FindDimensionByName(ctx, dimension.Name).Return(dimension, nil)
				out, err := gbService.DuplicateDimension(ctx, nameTarget, dimension.Name+"a")
				Expect(err).To(MatchError(fakeErr))
				Expect(out).To(Equal(dimension))
			})
		})

		When("given invalid input", func() {
			It("should error if dimension not found by id", func() {
				mockRepository.EXPECT().FindDimensionById(ctx, dimension.Id).Return(nil, nil)
				out, err := gbService.DuplicateDimension(ctx, idTarget, dimension.Name+"a")
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if dimension not found by name", func() {
				mockRepository.EXPECT().FindDimensionByName(ctx, dimension.Name).Return(nil, nil)
				out, err := gbService.DuplicateDimension(ctx, nameTarget, dimension.Name+"a")
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if error find by id ", func() {
				mockRepository.EXPECT().FindDimensionById(ctx, dimension.Id).Return(nil, fakeErr)
				out, err := gbService.DuplicateDimension(ctx, idTarget, dimension.Name+"a")
				Expect(err).To(MatchError(fakeErr))
				Expect(out).To(BeNil())
			})

			It("should error if error find by name", func() {
				mockRepository.EXPECT().FindDimensionByName(ctx, dimension.Name).Return(nil, fakeErr)
				out, err := gbService.DuplicateDimension(ctx, nameTarget, dimension.Name+"a")
				Expect(err).To(MatchError(fakeErr))
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("EditMap", func() {
		When("given valid input", func() {
			It("should work", func() {
				req := &pb.EditMapRequest{
					Target:             &pb.MapTarget{FindBy: &pb.MapTarget_Name{Name: m.Name}},
					OptionalName:       &pb.EditMapRequest_Name{Name: faker.Username()},
					OptionalPath:       &pb.EditMapRequest_Path{Path: faker.Username()},
					OptionalMaxPlayers: &pb.EditMapRequest_MaxPlayers{MaxPlayers: 123},
					OptionalInstanced:  &pb.EditMapRequest_Instanced{Instanced: true},
				}
				// Note: Maps were not changed so do not need to change
				expectedOut := &game.Map{}
				*expectedOut = *m
				expectedOut.Name = req.GetName()
				expectedOut.Path = req.GetPath()
				expectedOut.MaxPlayers = req.GetMaxPlayers()
				expectedOut.Instanced = req.GetInstanced()

				mockRepository.EXPECT().FindMapByName(gomock.Any(), m.Name).Return(m, nil)
				mockRepository.EXPECT().SaveMap(gomock.Any(), expectedOut).Return(expectedOut, fakeErr)
				out, err := gbService.EditMap(ctx, req)
				Expect(err).To(MatchError(fakeErr))
				Expect(out).To(Equal(expectedOut))
			})
		})

		When("given invalid input", func() {
			It("should fail if invalid m id", func() {
				req := &pb.EditMapRequest{
					Target: &pb.MapTarget{
						FindBy: &pb.MapTarget_Id{Id: "id"},
					},
				}
				out, err := gbService.EditMap(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if m not found by id", func() {
				req := &pb.EditMapRequest{
					Target: &pb.MapTarget{
						FindBy: &pb.MapTarget_Id{Id: m.Id.String()},
					},
				}
				mockRepository.EXPECT().FindMapById(gomock.Any(), m.Id).Return(nil, nil)
				out, err := gbService.EditMap(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if target type is unknown", func() {
				req := &pb.EditMapRequest{
					Target: &pb.MapTarget{},
					OptionalName: &pb.EditMapRequest_Name{
						Name: faker.Username(),
					},
				}
				out, err := gbService.EditMap(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("EditDimension", func() {
		When("given valid input", func() {
			It("should work", func() {
				req := &pb.EditDimensionRequest{
					Target: &pb.DimensionTarget{
						FindBy: &pb.DimensionTarget_Name{Name: dimension.Name},
					},
					OptionalName: &pb.EditDimensionRequest_Name{
						Name: faker.Username(),
					},
					OptionalVersion: &pb.EditDimensionRequest_Version{
						Version: faker.Username(),
					},
					EditMaps: true,
					MapIds:   []string{m.Id.String()},
					OptionalLocation: &pb.EditDimensionRequest_Location{
						Location: "us-central",
					},
				}
				// Note: Maps were not changed so do not need to change
				expectedOut := &game.Dimension{}
				*expectedOut = *dimension
				expectedOut.Name = req.GetName()
				expectedOut.Version = req.GetVersion()
				expectedOut.Location = req.GetLocation()

				mockRepository.EXPECT().FindDimensionByName(gomock.Any(), dimension.Name).Return(dimension, nil)
				mockRepository.EXPECT().SaveDimension(gomock.Any(), expectedOut).Return(expectedOut, fakeErr)
				mockRepository.EXPECT().FindMapsByIds(gomock.Any(), gomock.Any()).Return(game.Maps{m}, nil)
				out, err := gbService.EditDimension(ctx, req)
				Expect(err).To(MatchError(fakeErr))
				Expect(out).To(Equal(expectedOut))
			})
		})

		When("given invalid input", func() {
			It("should fail if invalid dimension id", func() {
				req := &pb.EditDimensionRequest{
					Target: &pb.DimensionTarget{
						FindBy: &pb.DimensionTarget_Id{Id: "id"},
					},
				}
				out, err := gbService.EditDimension(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if dimension not found by id", func() {
				req := &pb.EditDimensionRequest{
					Target: &pb.DimensionTarget{
						FindBy: &pb.DimensionTarget_Id{Id: dimension.Id.String()},
					},
				}
				mockRepository.EXPECT().FindDimensionById(gomock.Any(), dimension.Id).Return(nil, nil)
				out, err := gbService.EditDimension(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if target type is unknown", func() {
				req := &pb.EditDimensionRequest{
					Target: &pb.DimensionTarget{},
					OptionalName: &pb.EditDimensionRequest_Name{
						Name: faker.Username(),
					},
					OptionalVersion: &pb.EditDimensionRequest_Version{
						Version: faker.Username(),
					},
					EditMaps: true,
					MapIds:   []string{m.Id.String()},
					OptionalLocation: &pb.EditDimensionRequest_Location{
						Location: "us-central",
					},
				}
				out, err := gbService.EditDimension(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if map id is invalid", func() {
				req := &pb.EditDimensionRequest{
					Target: &pb.DimensionTarget{
						FindBy: &pb.DimensionTarget_Name{
							Name: dimension.Name,
						},
					},
					OptionalName: &pb.EditDimensionRequest_Name{
						Name: faker.Username(),
					},
					OptionalVersion: &pb.EditDimensionRequest_Version{
						Version: faker.Username(),
					},
					EditMaps: true,
					MapIds:   []string{"id"},
					OptionalLocation: &pb.EditDimensionRequest_Location{
						Location: "us-central",
					},
				}
				mockRepository.EXPECT().FindDimensionByName(gomock.Any(), dimension.Name).Return(dimension, nil)
				out, err := gbService.EditDimension(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if find map by ids errored", func() {
				req := &pb.EditDimensionRequest{
					Target: &pb.DimensionTarget{
						FindBy: &pb.DimensionTarget_Name{
							Name: dimension.Name,
						},
					},
					OptionalName: &pb.EditDimensionRequest_Name{
						Name: faker.Username(),
					},
					OptionalVersion: &pb.EditDimensionRequest_Version{
						Version: faker.Username(),
					},
					EditMaps: true,
					MapIds:   []string{m.Id.String()},
					OptionalLocation: &pb.EditDimensionRequest_Location{
						Location: "us-central",
					},
				}
				mockRepository.EXPECT().FindDimensionByName(gomock.Any(), dimension.Name).Return(dimension, nil)
				mockRepository.EXPECT().FindMapsByIds(gomock.Any(), gomock.Any()).Return(nil, fakeErr)
				out, err := gbService.EditDimension(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if any map ids aren't found", func() {
				req := &pb.EditDimensionRequest{
					Target: &pb.DimensionTarget{
						FindBy: &pb.DimensionTarget_Name{
							Name: dimension.Name,
						},
					},
					OptionalName: &pb.EditDimensionRequest_Name{
						Name: faker.Username(),
					},
					OptionalVersion: &pb.EditDimensionRequest_Version{
						Version: faker.Username(),
					},
					EditMaps: true,
					MapIds:   []string{m.Id.String()},
					OptionalLocation: &pb.EditDimensionRequest_Location{
						Location: "us-central",
					},
				}
				randUUID := uuid.New()
				mockRepository.EXPECT().FindDimensionByName(gomock.Any(), dimension.Name).Return(dimension, nil)
				mockRepository.EXPECT().FindMapsByIds(gomock.Any(), gomock.Any()).Return(game.Maps{
					&game.Map{Model: model.Model{Id: &randUUID}},
					&game.Map{Model: model.Model{Id: &randUUID}},
				}, nil)
				out, err := gbService.EditDimension(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if target type is unknown", func() {
				req := &pb.EditDimensionRequest{
					Target: &pb.DimensionTarget{
						FindBy: &pb.DimensionTarget_Name{
							Name: dimension.Name,
						},
					},
					OptionalName: &pb.EditDimensionRequest_Name{
						Name: faker.Username(),
					},
					OptionalVersion: &pb.EditDimensionRequest_Version{
						Version: faker.Username(),
					},
					EditMaps: true,
					MapIds:   []string{m.Id.String()},
					OptionalLocation: &pb.EditDimensionRequest_Location{
						Location: "us-central",
					},
				}
				mockRepository.EXPECT().FindDimensionByName(gomock.Any(), dimension.Name).Return(dimension, nil)
				mockRepository.EXPECT().FindMapsByIds(gomock.Any(), gomock.Any()).Return(nil, fakeErr)
				out, err := gbService.EditDimension(ctx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

		})
	})

	Describe("EditMap", func() {
		It("should work", func() {

		})
	})

	Describe("FindAllDimensions", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindAllDimensions(ctx).Return(game.Dimensions{dimension}, fakeErr)
			out, err := gbService.FindAllDimensions(ctx)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(dimension))
		})
	})

	Describe("FindAllMap", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindAllMaps(ctx).Return(game.Maps{m}, fakeErr)
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
			mockRepository.EXPECT().FindDimensionsByIds(ctx, []*uuid.UUID{dimension.Id}).Return(game.Dimensions{dimension}, fakeErr)
			out, err := gbService.FindDimensionsByIds(ctx, []*uuid.UUID{dimension.Id})
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(dimension))
		})
	})

	Describe("FindDimensionsByNames", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindDimensionsByNames(ctx, []string{dimension.Name}).Return(game.Dimensions{dimension}, fakeErr)
			out, err := gbService.FindDimensionsByNames(ctx, []string{dimension.Name})
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(dimension))
		})
	})

	Describe("FindDimensionsWithmapIds", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindDimensionsWithMapIds(ctx, []*uuid.UUID{m.Id}).Return(game.Dimensions{dimension}, fakeErr)
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
			mockRepository.EXPECT().FindMapsByIds(ctx, []*uuid.UUID{m.Id}).Return(game.Maps{m}, fakeErr)
			out, err := gbService.FindMapsByIds(ctx, []*uuid.UUID{m.Id})
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(m))
		})
	})

	Describe("FindMapsByNames", func() {
		It("should work", func() {
			mockRepository.EXPECT().FindMapsByNames(ctx, []string{m.Name}).Return(game.Maps{m}, fakeErr)
			out, err := gbService.FindMapsByNames(ctx, []string{m.Name})
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(ContainElement(m))
		})
	})
})

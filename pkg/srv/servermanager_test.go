package srv_test

import (
	context "context"
	"time"

	app "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.opentelemetry.io/otel"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

var _ = Describe("Servermanager server (local)", func() {
	var (
		hook           *test.Hook
		mockController *gomock.Controller
		ctx            context.Context

		globalConfig   *config.GlobalConfig
		conf           *app.GameBackendServerContext
		mockCharClient *mocks.MockCharacterServiceClient
		mockChatClient *mocks.MockChatServiceClient
		mockService    *mocks.MockGamebackendService
		server         pb.ServerManagerServiceServer

		dimension *model.Dimension
		m         *model.Map
	)

	BeforeEach(func() {
		var err error
		ctx = context.Background()
		log.Logger, hook = test.NewNullLogger()
		globalConfig = config.NewGlobalConfig(ctx)
		mockController = gomock.NewController(GinkgoT())
		mockCharClient = mocks.NewMockCharacterServiceClient(mockController)
		mockChatClient = mocks.NewMockChatServiceClient(mockController)
		mockService = mocks.NewMockGamebackendService(mockController)

		globalConfig.GameBackend.Mode = config.LocalMode
		conf = &app.GameBackendServerContext{
			GlobalConfig:       globalConfig,
			CharacterClient:    mockCharClient,
			ChatClient:         mockChatClient,
			GamebackendService: mockService,
			AgonesClient:       nil,
			KeycloakClient:     keycloak,
			Tracer:             otel.Tracer("test-servermanager"),
		}

		server, err = srv.NewServerManagerServiceServer(ctx, conf)
		Expect(err).NotTo(HaveOccurred())
		Expect(server).NotTo(BeNil())

		mapId := uuid.New()
		dimensionId := uuid.New()
		m = &model.Map{
			Model: model.Model{
				Id:        &mapId,
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
				Id:        &dimensionId,
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

		hook.Reset()
	})

	Describe("CreateDimension", func() {
		var (
			req *pb.CreateDimensionRequest
		)
		BeforeEach(func() {
			req = &pb.CreateDimensionRequest{
				Name:     dimension.Name,
				Version:  dimension.Version,
				MapIds:   []string{m.Id.String()},
				Location: dimension.Location,
			}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockService.EXPECT().
					CreateDimension(gomock.Any(), req.Name, req.Location, req.Version, []*uuid.UUID{m.Id}).
					Return(dimension, nil)
				out, err := server.CreateDimension(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error for invalid ctx (nil)", func() {
				out, err := server.CreateDimension(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error for empty ctx (empty)", func() {
				out, err := server.CreateDimension(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.CreateDimension(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (player)", func() {
				out, err := server.CreateDimension(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid map id", func() {
				req.MapIds = []string{"asdf"}
				out, err := server.CreateDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for error in creationg", func() {
				mockService.EXPECT().
					CreateDimension(gomock.Any(), req.Name, req.Location, req.Version, []*uuid.UUID{m.Id}).
					Return(nil, fakeErr)
				out, err := server.CreateDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("CreateMap", func() {
		var (
			req *pb.CreateMapRequest
		)
		BeforeEach(func() {
			req = &pb.CreateMapRequest{
				Name:       m.Name,
				Path:       m.Path,
				MaxPlayers: m.MaxPlayers,
				Instanced:  m.Instanced,
			}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockService.EXPECT().
					CreateMap(gomock.Any(), req.Name, req.Path, req.MaxPlayers, req.Instanced).
					Return(m, nil)
				out, err := server.CreateMap(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error for invalid ctx (nil)", func() {
				out, err := server.CreateMap(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error for empty ctx (empty)", func() {
				out, err := server.CreateMap(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.CreateMap(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.CreateMap(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for error in creationg", func() {
				mockService.EXPECT().
					CreateMap(gomock.Any(), req.Name, req.Path, req.MaxPlayers, req.Instanced).
					Return(nil, fakeErr)
				out, err := server.CreateMap(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("DeleteDimension", func() {
		var (
			req *pb.DimensionTarget
		)
		BeforeEach(func() {
			req = &pb.DimensionTarget{
				FindBy: &pb.DimensionTarget_Id{
					Id: dimension.Id.String(),
				},
			}
		})
		When("given valid input", func() {
			It("should work for Id target (admin)", func() {
				mockService.EXPECT().
					FindDimension(gomock.Any(), req).
					Return(dimension, nil)
				mockService.EXPECT().
					DeleteDimensionById(gomock.Any(), dimension.Id).
					Return(nil)
				out, err := server.DeleteDimension(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work for Name target (admin)", func() {
				req.FindBy = &pb.DimensionTarget_Name{
					Name: dimension.Name,
				}
				mockService.EXPECT().
					FindDimension(gomock.Any(), req).
					Return(dimension, nil)
				mockService.EXPECT().
					DeleteDimensionById(gomock.Any(), dimension.Id).
					Return(nil)
				out, err := server.DeleteDimension(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error for invalid ctx (nil)", func() {
				out, err := server.DeleteDimension(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error for empty ctx (empty)", func() {
				out, err := server.DeleteDimension(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.DeleteDimension(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.DeleteDimension(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindDimension error", func() {
				mockService.EXPECT().
					FindDimension(gomock.Any(), req).
					Return(nil, fakeErr)
				out, err := server.DeleteDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for not found", func() {
				mockService.EXPECT().
					FindDimension(gomock.Any(), req).
					Return(nil, nil)
				out, err := server.DeleteDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindByName not found", func() {
				req.FindBy = &pb.DimensionTarget_Name{
					Name: dimension.Name,
				}
				mockService.EXPECT().
					FindDimension(gomock.Any(), req).
					Return(nil, nil)
				out, err := server.DeleteDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for error on deletion", func() {
				mockService.EXPECT().
					FindDimension(gomock.Any(), req).
					Return(dimension, nil)
				mockService.EXPECT().
					DeleteDimensionById(gomock.Any(), dimension.Id).
					Return(fakeErr)
				out, err := server.DeleteDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("DeleteMap", func() {
		var (
			req *pb.MapTarget
		)
		BeforeEach(func() {
			req = &pb.MapTarget{
				FindBy: &pb.MapTarget_Id{
					Id: m.Id.String(),
				},
			}
		})
		When("given valid input", func() {
			It("should work for Id target (admin)", func() {
				mockService.EXPECT().
					FindMapById(gomock.Any(), m.Id).
					Return(m, nil)
				mockService.EXPECT().
					DeleteMapById(gomock.Any(), m.Id).
					Return(nil)
				out, err := server.DeleteMap(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work for Name target (admin)", func() {
				req.FindBy = &pb.MapTarget_Name{
					Name: m.Name,
				}
				mockService.EXPECT().
					FindMapByName(gomock.Any(), m.Name).
					Return(m, nil)
				mockService.EXPECT().
					DeleteMapById(gomock.Any(), m.Id).
					Return(nil)
				out, err := server.DeleteMap(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error for invalid ctx (nil)", func() {
				out, err := server.DeleteMap(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error for empty ctx (empty)", func() {
				out, err := server.DeleteMap(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.DeleteMap(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.DeleteMap(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for error on FindById", func() {
				mockService.EXPECT().
					FindMapById(gomock.Any(), m.Id).
					Return(nil, fakeErr)
				out, err := server.DeleteMap(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for error on FindByName", func() {
				req.FindBy = &pb.MapTarget_Name{
					Name: m.Name,
				}
				mockService.EXPECT().
					FindMapByName(gomock.Any(), m.Name).
					Return(nil, fakeErr)
				out, err := server.DeleteMap(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindById not found", func() {
				mockService.EXPECT().
					FindMapById(gomock.Any(), m.Id).
					Return(nil, nil)
				out, err := server.DeleteMap(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid id", func() {
				req.FindBy = &pb.MapTarget_Id{
					Id: "asdf",
				}
				out, err := server.DeleteMap(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindByName not found", func() {
				req.FindBy = &pb.MapTarget_Name{
					Name: m.Name,
				}
				mockService.EXPECT().
					FindMapByName(gomock.Any(), m.Name).
					Return(nil, nil)
				out, err := server.DeleteMap(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for error on unknown target type", func() {
				req.FindBy = nil
				out, err := server.DeleteMap(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for error on deletion", func() {
				mockService.EXPECT().
					FindMapById(gomock.Any(), m.Id).
					Return(m, nil)
				mockService.EXPECT().
					DeleteMapById(gomock.Any(), m.Id).
					Return(fakeErr)
				out, err := server.DeleteMap(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("DuplicateDimension", func() {
		var (
			req *pb.DuplicateDimensionRequest
		)
		BeforeEach(func() {
			req = &pb.DuplicateDimensionRequest{
				Target: &pb.DimensionTarget{
					FindBy: &pb.DimensionTarget_Id{
						Id: dimension.Id.String(),
					},
				},
				Name: faker.Username(),
			}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockService.EXPECT().
					DuplicateDimension(gomock.Any(), req.Target, req.Name).
					Return(dimension, nil)
				out, err := server.DuplicateDimension(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error for invalid ctx (nil)", func() {
				out, err := server.DuplicateDimension(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error for empty ctx (empty)", func() {
				out, err := server.DuplicateDimension(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.DuplicateDimension(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (player)", func() {
				out, err := server.DuplicateDimension(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for duplicate dimension service errors", func() {
				mockService.EXPECT().
					DuplicateDimension(gomock.Any(), req.Target, req.Name).
					Return(nil, fakeErr)
				out, err := server.DuplicateDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for duplicate dimension service no dimension returned", func() {
				mockService.EXPECT().
					DuplicateDimension(gomock.Any(), req.Target, req.Name).
					Return(nil, nil)
				out, err := server.DuplicateDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("EditDimension", func() {
		var (
			req *pb.EditDimensionRequest
		)
		BeforeEach(func() {
			req = &pb.EditDimensionRequest{
				Target: &pb.DimensionTarget{
					FindBy: &pb.DimensionTarget_Id{
						Id: dimension.Id.String(),
					},
				},
				OptionalName: &pb.EditDimensionRequest_Name{
					Name: faker.Username(),
				},
				OptionalVersion: &pb.EditDimensionRequest_Version{
					Version: faker.Username(),
				},
				EditMaps: true,
				MapIds:   []string{},
				OptionalLocation: &pb.EditDimensionRequest_Location{
					Location: faker.Username(),
				},
			}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockService.EXPECT().FindDimension(gomock.Any(), req.Target).Return(dimension, nil)
				mockService.EXPECT().EditDimension(gomock.Any(), req).Return(dimension, nil)
				out, err := server.EditDimension(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error for invalid ctx (nil)", func() {
				out, err := server.EditDimension(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error for empty ctx (empty)", func() {
				out, err := server.EditDimension(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.EditDimension(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (player)", func() {
				out, err := server.EditDimension(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindDimension errors", func() {
				mockService.EXPECT().FindDimension(gomock.Any(), req.Target).Return(nil, fakeErr)
				out, err := server.EditDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindDimension returns no matches", func() {
				mockService.EXPECT().FindDimension(gomock.Any(), req.Target).Return(nil, nil)
				out, err := server.EditDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for EditDimension returns an error", func() {
				mockService.EXPECT().FindDimension(gomock.Any(), req.Target).Return(dimension, nil)
				mockService.EXPECT().EditDimension(gomock.Any(), req).Return(nil, fakeErr)
				out, err := server.EditDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	// Describe("", func() {
	// 	var (
	// 		req *pb.
	// 	)
	// 	BeforeEach(func() {
	// 		req =
	// 	})
	// 	When("given valid input", func() {
	// 		It("should work (admin)", func() {
	// 			out, err := server.(ctx, req)
	// 			Expect(err).NotTo(HaveOccurred())
	// 			Expect(out).NotTo(BeNil())
	// 		})
	// 	})
	//
	// 	When("given invalid input", func() {
	// 		It("should error for invalid ctx (nil)", func() {
	// 			out, err := server.(nil, req)
	// 			Expect(err).To(HaveOccurred())
	// 			Expect(out).To(BeNil())
	// 		})
	// 		It("should error for empty ctx (empty)", func() {
	// 			out, err := server.(context.Background(), req)
	// 			Expect(err).To(HaveOccurred())
	// 			Expect(out).To(BeNil())
	// 		})
	//
	// 		It("should error for invalid permission (guest)", func() {
	// 			out, err := server.(incGuestCtx, req)
	// 			Expect(err).To(HaveOccurred())
	// 			Expect(out).To(BeNil())
	// 		})
	//
	// 		It("should error for invalid permission (player)", func() {
	// 			out, err := server.(incPlayerCtx, req)
	// 			Expect(err).To(HaveOccurred())
	// 			Expect(out).To(BeNil())
	// 		})
	// 	})
	// })

})

package srv_test

import (
	context "context"
	"strings"
	"time"

	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	autoscalingv1 "agones.dev/agones/pkg/apis/autoscaling/v1"
	"agones.dev/agones/pkg/testing"
	app "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/model/game"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.opentelemetry.io/otel"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	k8stesting "k8s.io/client-go/testing"
)

var _ = Describe("Servermanager server", func() {
	var (
		hook           *test.Hook
		mockController *gomock.Controller
		ctx            context.Context

		globalConfig   *config.GlobalConfig
		conf           *app.GameBackendServerContext
		mockCharClient *mocks.MockCharacterServiceClient
		mockChatClient *mocks.MockChatServiceClient
		mockService    *mocks.MockGamebackendService
		mockAgones     testing.Mocks
		server         pb.ServerManagerServiceServer

		dimension *game.Dimension
		m         *game.Map
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
		mockAgones = testing.NewMocks()

		conf = &app.GameBackendServerContext{
			GlobalConfig:       globalConfig,
			CharacterClient:    mockCharClient,
			ChatClient:         mockChatClient,
			GamebackendService: mockService,
			AgonesClient:       mockAgones.AgonesClient,
			KeycloakClient:     keycloak,
			Tracer:             otel.Tracer("test-servermanager"),
		}
		conf.GlobalConfig.GameBackend.Mode = config.LocalMode

		server, err = srv.NewServerManagerServiceServer(ctx, conf)
		Expect(err).NotTo(HaveOccurred())
		Expect(server).NotTo(BeNil())

		mapId := uuid.New()
		dimensionId := uuid.New()
		m = &game.Map{
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
		dimension = &game.Dimension{
			Model: model.Model{
				Id:        &dimensionId,
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
		Context("local mode", func() {
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

		Context("with agones", func() {
			var (
				fleet           *agonesv1.Fleet
				fleetAutoscaler *autoscalingv1.FleetAutoscaler
			)
			BeforeEach(func() {
				conf.GlobalConfig.GameBackend.Mode = config.ModeProduction
				mockService.EXPECT().
					CreateDimension(gomock.Any(), req.Name, req.Location, req.Version, []*uuid.UUID{m.Id}).
					Return(dimension, nil)
			})

			When("agones is working", func() {
				It("should work", func() {
					mockAgones.AgonesClient.AddReactor("create", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.CreateAction)
						fleet = ua.GetObject().(*agonesv1.Fleet)
						return true, fleet, nil
					})
					mockAgones.AgonesClient.AddReactor("create", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.CreateAction)
						fleetAutoscaler = ua.GetObject().(*autoscalingv1.FleetAutoscaler)
						return true, fleetAutoscaler, nil
					})
					out, err := server.CreateDimension(incAdminCtx, req)
					Expect(err).NotTo(HaveOccurred())
					Expect(out).NotTo(BeNil())

					Expect(fleet).NotTo(BeNil())
					Expect(fleet.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(fleet.Name).To(ContainSubstring(strings.ToLower(m.Name)))
					Expect(fleetAutoscaler).NotTo(BeNil())
					Expect(fleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(fleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(m.Name)))
					Expect(fleet.Validate(testing.FakeAPIHooks{})).To(BeNil())
					Expect(fleetAutoscaler.Validate()).To(BeNil())
				})
			})

			When("agones isn't working", func() {
				It("should err on fleet creation error", func() {
					mockAgones.AgonesClient.AddReactor("create", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						return true, nil, fakeErr
					})
					mockAgones.AgonesClient.AddReactor("create", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.CreateAction)
						fleetAutoscaler = ua.GetObject().(*autoscalingv1.FleetAutoscaler)
						return true, fleetAutoscaler, nil
					})
					out, err := server.CreateDimension(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).NotTo(BeNil())
					Expect(fleet).To(BeNil())
					Expect(fleetAutoscaler).To(BeNil())
				})

				It("should err on fleet autoscaling creation error", func() {
					mockAgones.AgonesClient.AddReactor("create", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.CreateAction)
						fleet = ua.GetObject().(*agonesv1.Fleet)
						return true, fleet, nil
					})
					mockAgones.AgonesClient.AddReactor("create", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						return true, nil, fakeErr
					})
					out, err := server.CreateDimension(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).NotTo(BeNil())
					Expect(fleet).NotTo(BeNil())
					Expect(fleet.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(fleet.Name).To(ContainSubstring(strings.ToLower(m.Name)))
					Expect(fleetAutoscaler).To(BeNil())
				})
			})

			It("should err if agones not setup", func() {
				conf.AgonesClient = nil
				out, err := server.CreateDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(fleet).To(BeNil())
				Expect(fleetAutoscaler).To(BeNil())
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
		Context("local mode", func() {
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

		Context("with agones", func() {
			var (
				fleetDeleted           bool
				fleetAutoscalerDeleted bool
			)
			BeforeEach(func() {
				conf.GlobalConfig.GameBackend.Mode = config.ModeProduction
				mockService.EXPECT().
					FindDimension(gomock.Any(), req).
					Return(dimension, nil)
				mockService.EXPECT().
					DeleteDimensionById(gomock.Any(), dimension.Id).
					Return(nil)
				fleetDeleted = false
				fleetAutoscalerDeleted = false
			})

			When("agones is working", func() {
				It("should work", func() {
					mockAgones.AgonesClient.AddReactor("delete", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						fleetAutoscalerDeleted = true
						return true, nil, nil
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						fleetDeleted = true
						return true, nil, nil
					})
					out, err := server.DeleteDimension(incAdminCtx, req)
					Expect(err).NotTo(HaveOccurred())
					Expect(out).NotTo(BeNil())
					Expect(fleetDeleted).To(BeTrue())
					Expect(fleetAutoscalerDeleted).To(BeTrue())
				})
			})

			When("agones isn't working", func() {
				It("should error on autoscaling deleting issue", func() {
					mockAgones.AgonesClient.AddReactor("delete", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						fleetAutoscalerDeleted = true
						return true, nil, fakeErr
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						fleetDeleted = true
						return true, nil, nil
					})
					out, err := server.DeleteDimension(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).NotTo(BeNil())
					Expect(fleetDeleted).To(BeFalse())
					Expect(fleetAutoscalerDeleted).To(BeTrue())
				})

				It("should error on fleet deleting issue", func() {
					mockAgones.AgonesClient.AddReactor("delete", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						fleetAutoscalerDeleted = true
						return true, nil, nil
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						fleetDeleted = true
						return true, nil, fakeErr
					})
					out, err := server.DeleteDimension(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).NotTo(BeNil())
					Expect(fleetDeleted).To(BeTrue())
					Expect(fleetAutoscalerDeleted).To(BeTrue())
				})

				It("should error if agones is not setup", func() {
					conf.AgonesClient = nil
					mockAgones.AgonesClient.AddReactor("delete", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						fleetAutoscalerDeleted = true
						return true, nil, nil
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						fleetDeleted = true
						return true, nil, fakeErr
					})
					out, err := server.DeleteDimension(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).NotTo(BeNil())
					Expect(fleetDeleted).To(BeFalse())
					Expect(fleetAutoscalerDeleted).To(BeFalse())
				})
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

		Context("local mode", func() {
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

		Context("with agones", func() {
			It("should delete gameservers running this map", func() {
				conf.GlobalConfig.GameBackend.Mode = config.ModeProduction

				var (
					deletedFleet           string
					deletedFleetAutoscaler string
				)

				mockAgones.AgonesClient.AddReactor("delete", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
					da := action.(k8stesting.DeleteAction)
					deletedFleet = da.GetName()
					return true, nil, fakeErr
				})
				mockAgones.AgonesClient.AddReactor("delete", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
					da := action.(k8stesting.DeleteAction)
					deletedFleetAutoscaler = da.GetName()
					return true, nil, nil
				})

				mockService.EXPECT().
					FindMapById(gomock.Any(), m.Id).
					Return(m, nil)
				mockService.EXPECT().
					DeleteMapById(gomock.Any(), m.Id).
					Return(nil)
				mockService.EXPECT().
					FindDimensionsWithMapIds(gomock.Any(), gomock.Any()).
					Return(game.Dimensions{dimension}, nil)
				out, err := server.DeleteMap(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).NotTo(BeNil())

				Expect(deletedFleet).To(ContainSubstring(strings.ToLower(dimension.Name)))
				Expect(deletedFleet).To(ContainSubstring(strings.ToLower(m.Name)))
				Expect(deletedFleetAutoscaler).To(ContainSubstring(strings.ToLower(dimension.Name)))
				Expect(deletedFleetAutoscaler).To(ContainSubstring(strings.ToLower(m.Name)))
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

		Context("local mode", func() {
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

		Context("with agones", func() {
			It("should setup new dimension", func() {
				conf.GlobalConfig.GameBackend.Mode = config.ModeProduction

				var (
					fleet           *agonesv1.Fleet
					fleetAutoscaler *autoscalingv1.FleetAutoscaler
				)
				mockAgones.AgonesClient.AddReactor("create", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
					ua := action.(k8stesting.CreateAction)
					fleet = ua.GetObject().(*agonesv1.Fleet)
					return true, fleet, nil
				})
				mockAgones.AgonesClient.AddReactor("create", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
					ua := action.(k8stesting.CreateAction)
					fleetAutoscaler = ua.GetObject().(*autoscalingv1.FleetAutoscaler)
					return true, fleetAutoscaler, fakeErr
				})
				mockService.EXPECT().
					DuplicateDimension(gomock.Any(), req.Target, req.Name).
					Return(dimension, nil)
				out, err := server.DuplicateDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).NotTo(BeNil())

				Expect(fleet).NotTo(BeNil())
				Expect(fleet.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
				Expect(fleet.Name).To(ContainSubstring(strings.ToLower(m.Name)))
				Expect(fleetAutoscaler).NotTo(BeNil())
				Expect(fleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
				Expect(fleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(m.Name)))
				Expect(fleet.Validate(testing.FakeAPIHooks{})).To(BeNil())
				Expect(fleetAutoscaler.Validate()).To(BeNil())
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
		Context("local mode", func() {
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
		Context("with agones", func() {
			var (
				updatedFleet           *agonesv1.Fleet
				updatedFleetAutoscaler *autoscalingv1.FleetAutoscaler
				createdFleet           *agonesv1.Fleet
				createdFleetAutoscaler *autoscalingv1.FleetAutoscaler
				deletedFleet           string
				deletedFleetAutoscaler string
				m2                     *game.Map
				m3                     *game.Map
				newDimension           *game.Dimension
			)
			BeforeEach(func() {
				conf.GlobalConfig.GameBackend.Mode = config.ModeProduction
				req.OptionalName = nil
				req.EditMaps = true
				req.OptionalVersion = &pb.EditDimensionRequest_Version{
					Version: "v2",
				}
				m2Id := uuid.New()
				m2 = &game.Map{
					Model: model.Model{
						Id:        &m2Id,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						DeletedAt: gorm.DeletedAt{},
					},
					Name:       faker.Username(),
					Path:       faker.Username(),
					MaxPlayers: 40,
					Instanced:  false,
				}
				m2.Dimensions = []*game.Dimension{dimension}
				m3Id := uuid.New()
				m3 = &game.Map{
					Model: model.Model{
						Id:        &m3Id,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						DeletedAt: gorm.DeletedAt{},
					},
					Name:       faker.Username(),
					Path:       faker.Username(),
					MaxPlayers: 40,
					Instanced:  false,
				}
				m3.Dimensions = []*game.Dimension{dimension}
				dimension.Maps = game.Maps{m, m2}
				newDimension = &game.Dimension{
					Model: model.Model{
						Id:        dimension.Id,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Name:     dimension.Name,
					Location: req.GetLocation(),
					Version:  req.GetVersion(),
					Maps:     []*game.Map{m2, m3},
				}
				req.MapIds = []string{m2.Id.String(), m3.Id.String()}
				mockService.EXPECT().FindDimension(gomock.Any(), req.Target).Return(dimension, nil)
				mockService.EXPECT().EditDimension(gomock.Any(), req).Return(newDimension, nil)
			})

			When("agones is working", func() {
				It("should update correctly given map changes", func() {
					mockAgones.AgonesClient.AddReactor("update", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.UpdateAction)
						updatedFleet = ua.GetObject().(*agonesv1.Fleet)
						return true, updatedFleet, nil
					})
					mockAgones.AgonesClient.AddReactor("update", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.UpdateAction)
						updatedFleetAutoscaler = ua.GetObject().(*autoscalingv1.FleetAutoscaler)
						return true, updatedFleetAutoscaler, nil
					})
					mockAgones.AgonesClient.AddReactor("create", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ca := action.(k8stesting.CreateAction)
						createdFleet = ca.GetObject().(*agonesv1.Fleet)
						return true, createdFleet, nil
					})
					mockAgones.AgonesClient.AddReactor("create", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ca := action.(k8stesting.CreateAction)
						createdFleetAutoscaler = ca.GetObject().(*autoscalingv1.FleetAutoscaler)
						return true, createdFleetAutoscaler, nil
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						da := action.(k8stesting.DeleteAction)
						deletedFleet = da.GetName()
						return true, nil, nil
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						da := action.(k8stesting.DeleteAction)
						deletedFleetAutoscaler = da.GetName()
						return true, nil, nil
					})
					out, err := server.EditDimension(incAdminCtx, req)
					Expect(err).NotTo(HaveOccurred())
					Expect(out).NotTo(BeNil())

					// Removed m, kept m2, and m3 is new
					Expect(deletedFleet).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(deletedFleet).To(ContainSubstring(strings.ToLower(m.Name)))
					Expect(deletedFleetAutoscaler).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(deletedFleetAutoscaler).To(ContainSubstring(strings.ToLower(m.Name)))
					Expect(updatedFleet).NotTo(BeNil())
					Expect(updatedFleet.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(updatedFleet.Name).To(ContainSubstring(strings.ToLower(m2.Name)))
					Expect(updatedFleetAutoscaler).NotTo(BeNil())
					Expect(updatedFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(updatedFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(m2.Name)))
					Expect(updatedFleet.Validate(testing.FakeAPIHooks{})).To(BeNil())
					Expect(updatedFleetAutoscaler.Validate()).To(BeNil())
					Expect(createdFleet).NotTo(BeNil())
					Expect(createdFleet.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(createdFleet.Name).To(ContainSubstring(strings.ToLower(m3.Name)))
					Expect(createdFleetAutoscaler).NotTo(BeNil())
					Expect(createdFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(createdFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(m3.Name)))
					Expect(createdFleet.Validate(testing.FakeAPIHooks{})).To(BeNil())
					Expect(createdFleetAutoscaler.Validate()).To(BeNil())
				})

				It("should update correctly given name change", func() {
					req.OptionalName = &pb.EditDimensionRequest_Name{
						Name: "newname",
					}
					newDimension.Name = req.GetName()
					mockAgones.AgonesClient.AddReactor("create", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ca := action.(k8stesting.CreateAction)
						createdFleet = ca.GetObject().(*agonesv1.Fleet)
						Expect(createdFleet).NotTo(BeNil())
						Expect(createdFleet.Name).To(Or(ContainSubstring(strings.ToLower(m2.Name)), ContainSubstring(strings.ToLower(m3.Name))))
						Expect(createdFleet.Name).To(ContainSubstring(strings.ToLower(req.GetName())))
						Expect(createdFleet.Validate(testing.FakeAPIHooks{})).To(BeNil())
						return true, createdFleet, nil
					})
					mockAgones.AgonesClient.AddReactor("create", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ca := action.(k8stesting.CreateAction)
						createdFleetAutoscaler = ca.GetObject().(*autoscalingv1.FleetAutoscaler)
						Expect(createdFleetAutoscaler).NotTo(BeNil())
						Expect(createdFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(req.GetName())))
						Expect(createdFleetAutoscaler.Name).To(Or(ContainSubstring(strings.ToLower(m2.Name)), ContainSubstring(strings.ToLower(m3.Name))))
						Expect(createdFleetAutoscaler.Validate()).To(BeNil())
						return true, createdFleetAutoscaler, fakeErr
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						da := action.(k8stesting.DeleteAction)
						deletedFleet = da.GetName()
						Expect(deletedFleet).To(Or(ContainSubstring(strings.ToLower(m2.Name)), ContainSubstring(strings.ToLower(m.Name))))
						Expect(deletedFleet).To(ContainSubstring(strings.ToLower(dimension.Name)))
						return true, nil, nil
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						da := action.(k8stesting.DeleteAction)
						deletedFleetAutoscaler = da.GetName()
						Expect(deletedFleetAutoscaler).To(Or(ContainSubstring(strings.ToLower(m2.Name)), ContainSubstring(strings.ToLower(m.Name))))
						Expect(deletedFleetAutoscaler).To(ContainSubstring(strings.ToLower(dimension.Name)))
						return true, nil, fakeErr
					})
					out, err := server.EditDimension(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).NotTo(BeNil())
				})

				It("should update correctly given only a version change", func() {
					req.EditMaps = false
					newDimension.Maps = game.Maps{m, m2}
					mockAgones.AgonesClient.AddReactor("update", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ca := action.(k8stesting.UpdateAction)
						updatedFleet = ca.GetObject().(*agonesv1.Fleet)
						Expect(updatedFleet).NotTo(BeNil())
						Expect(updatedFleet.Name).To(Or(ContainSubstring(strings.ToLower(m.Name)), ContainSubstring(strings.ToLower(m2.Name))))
						Expect(updatedFleet.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
						Expect(updatedFleet.Validate(testing.FakeAPIHooks{})).To(BeNil())
						return true, updatedFleet, nil
					})
					mockAgones.AgonesClient.AddReactor("update", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ca := action.(k8stesting.UpdateAction)
						updatedFleetAutoscaler = ca.GetObject().(*autoscalingv1.FleetAutoscaler)
						Expect(updatedFleetAutoscaler).NotTo(BeNil())
						Expect(updatedFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
						Expect(updatedFleetAutoscaler.Name).To(Or(ContainSubstring(strings.ToLower(m.Name)), ContainSubstring(strings.ToLower(m2.Name))))
						Expect(updatedFleetAutoscaler.Validate()).To(BeNil())
						return true, updatedFleetAutoscaler, fakeErr
					})
					out, err := server.EditDimension(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).NotTo(BeNil())
				})
			})

			When("agones isn't working", func() {
				It("should combine all errs that occured", func() {
					mockAgones.AgonesClient.AddReactor("update", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.UpdateAction)
						updatedFleet = ua.GetObject().(*agonesv1.Fleet)
						return true, updatedFleet, nil
					})
					mockAgones.AgonesClient.AddReactor("update", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.UpdateAction)
						updatedFleetAutoscaler = ua.GetObject().(*autoscalingv1.FleetAutoscaler)
						return true, updatedFleetAutoscaler, fakeErr
					})
					mockAgones.AgonesClient.AddReactor("create", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.CreateAction)
						createdFleet = ua.GetObject().(*agonesv1.Fleet)
						return true, createdFleet, nil
					})
					mockAgones.AgonesClient.AddReactor("create", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.CreateAction)
						createdFleetAutoscaler = ua.GetObject().(*autoscalingv1.FleetAutoscaler)
						return true, createdFleetAutoscaler, fakeErr
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						da := action.(k8stesting.DeleteAction)
						deletedFleet = da.GetName()
						return true, nil, fakeErr
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						da := action.(k8stesting.DeleteAction)
						deletedFleetAutoscaler = da.GetName()
						return true, nil, nil
					})
					out, err := server.EditDimension(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).NotTo(BeNil())

					// Removed m, kept m2, and m3 is new
					Expect(deletedFleet).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(deletedFleet).To(ContainSubstring(strings.ToLower(m.Name)))
					Expect(deletedFleetAutoscaler).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(deletedFleetAutoscaler).To(ContainSubstring(strings.ToLower(m.Name)))
					Expect(updatedFleet).NotTo(BeNil())
					Expect(updatedFleet.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(updatedFleet.Name).To(ContainSubstring(strings.ToLower(m2.Name)))
					Expect(updatedFleetAutoscaler).NotTo(BeNil())
					Expect(updatedFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(updatedFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(m2.Name)))
					Expect(updatedFleet.Validate(testing.FakeAPIHooks{})).To(BeNil())
					Expect(updatedFleetAutoscaler.Validate()).To(BeNil())
					Expect(createdFleet).NotTo(BeNil())
					Expect(createdFleet.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(createdFleet.Name).To(ContainSubstring(strings.ToLower(m3.Name)))
					Expect(createdFleetAutoscaler).NotTo(BeNil())
					Expect(createdFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(createdFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(m3.Name)))
					Expect(createdFleet.Validate(testing.FakeAPIHooks{})).To(BeNil())
					Expect(createdFleetAutoscaler.Validate()).To(BeNil())
				})

				It("should err if agones not setup", func() {
					conf.AgonesClient = nil
					out, err := server.EditDimension(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).NotTo(BeNil())
					Expect(updatedFleet).To(BeNil())
					Expect(updatedFleetAutoscaler).To(BeNil())
				})
			})
		})
	})

	Describe("EditMap", func() {
		var (
			req *pb.EditMapRequest
		)
		BeforeEach(func() {
			req = &pb.EditMapRequest{
				Target: &pb.MapTarget{
					FindBy: &pb.MapTarget_Id{
						Id: dimension.Id.String(),
					},
				},
				OptionalName: &pb.EditMapRequest_Name{
					Name: faker.Username(),
				},
				OptionalPath: &pb.EditMapRequest_Path{
					Path: faker.Username(),
				},
				OptionalInstanced: &pb.EditMapRequest_Instanced{
					Instanced: true,
				},
				OptionalMaxPlayers: &pb.EditMapRequest_MaxPlayers{
					MaxPlayers: 5,
				},
			}
		})
		Context("local mode", func() {
			When("given valid input", func() {
				It("should work (admin)", func() {
					mockService.EXPECT().FindMap(gomock.Any(), req.Target).Return(m, nil)
					mockService.EXPECT().EditMap(gomock.Any(), req).Return(m, nil)
					out, err := server.EditMap(incAdminCtx, req)
					Expect(err).NotTo(HaveOccurred())
					Expect(out).NotTo(BeNil())
				})
			})

			When("given invalid input", func() {
				It("should error for invalid ctx (nil)", func() {
					out, err := server.EditMap(nil, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})
				It("should error for empty ctx (empty)", func() {
					out, err := server.EditMap(context.Background(), req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should error for invalid permission (guest)", func() {
					out, err := server.EditMap(incGuestCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should error for invalid permission (player)", func() {
					out, err := server.EditMap(incPlayerCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should error for FindMap errors", func() {
					mockService.EXPECT().FindMap(gomock.Any(), req.Target).Return(nil, fakeErr)
					out, err := server.EditMap(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should error for FindMap returns no matches", func() {
					mockService.EXPECT().FindMap(gomock.Any(), req.Target).Return(nil, nil)
					out, err := server.EditMap(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should error for EditMap returns an error", func() {
					mockService.EXPECT().FindMap(gomock.Any(), req.Target).Return(m, nil)
					mockService.EXPECT().EditMap(gomock.Any(), req).Return(nil, fakeErr)
					out, err := server.EditMap(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})
			})
		})
		Context("with agones", func() {
			var (
				updatedFleet           *agonesv1.Fleet
				updatedFleetAutoscaler *autoscalingv1.FleetAutoscaler
				createdFleet           *agonesv1.Fleet
				createdFleetAutoscaler *autoscalingv1.FleetAutoscaler
				deletedFleet           string
				deletedFleetAutoscaler string
				m2                     *game.Map
			)
			BeforeEach(func() {
				conf.GlobalConfig.GameBackend.Mode = config.ModeProduction
				m2 = &game.Map{
					Model: model.Model{
						Id:        m.Id,
						CreatedAt: m.CreatedAt,
						UpdatedAt: time.Now(),
					},
					Name:       req.GetName(),
					Path:       req.GetPath(),
					MaxPlayers: req.GetMaxPlayers(),
					Instanced:  req.GetInstanced(),
					Dimensions: []*game.Dimension{dimension},
				}
				mockService.EXPECT().FindMap(gomock.Any(), req.Target).Return(m, nil)
				mockService.EXPECT().EditMap(gomock.Any(), req).Return(m2, nil)

				req.OptionalName = nil
			})

			When("name changes", func() {
				It("should delete and recreate servers", func() {
					mockAgones.AgonesClient.AddReactor("create", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ca := action.(k8stesting.CreateAction)
						createdFleet = ca.GetObject().(*agonesv1.Fleet)
						return true, createdFleet, nil
					})
					mockAgones.AgonesClient.AddReactor("create", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ca := action.(k8stesting.CreateAction)
						createdFleetAutoscaler = ca.GetObject().(*autoscalingv1.FleetAutoscaler)
						return true, createdFleetAutoscaler, fakeErr
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						da := action.(k8stesting.DeleteAction)
						deletedFleet = da.GetName()
						return true, nil, fakeErr
					})
					mockAgones.AgonesClient.AddReactor("delete", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						da := action.(k8stesting.DeleteAction)
						deletedFleetAutoscaler = da.GetName()
						return true, nil, nil
					})
					out, err := server.EditMap(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).NotTo(BeNil())

					Expect(deletedFleet).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(deletedFleet).To(ContainSubstring(strings.ToLower(m.Name)))
					Expect(deletedFleetAutoscaler).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(deletedFleetAutoscaler).To(ContainSubstring(strings.ToLower(m.Name)))
					Expect(createdFleet).NotTo(BeNil())
					Expect(createdFleet.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(createdFleet.Name).To(ContainSubstring(strings.ToLower(m2.Name)))
					Expect(createdFleetAutoscaler).NotTo(BeNil())
					Expect(createdFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(createdFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(m2.Name)))
					Expect(createdFleet.Validate(testing.FakeAPIHooks{})).To(BeNil())
					Expect(createdFleetAutoscaler.Validate()).To(BeNil())
				})
			})

			When("name doesn't change", func() {
				It("should update for all dimension", func() {
					m2.Name = m.Name
					req.OptionalName = &pb.EditMapRequest_Name{
						Name: m2.Name,
					}
					mockAgones.AgonesClient.AddReactor("update", "fleets", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.UpdateAction)
						updatedFleet = ua.GetObject().(*agonesv1.Fleet)
						return true, updatedFleet, nil
					})
					mockAgones.AgonesClient.AddReactor("update", "fleetautoscalers", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
						ua := action.(k8stesting.UpdateAction)
						updatedFleetAutoscaler = ua.GetObject().(*autoscalingv1.FleetAutoscaler)
						return true, updatedFleetAutoscaler, fakeErr
					})

					out, err := server.EditMap(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).NotTo(BeNil())

					Expect(updatedFleet).NotTo(BeNil())
					Expect(updatedFleet.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(updatedFleet.Name).To(ContainSubstring(strings.ToLower(m.Name)))
					Expect(updatedFleetAutoscaler).NotTo(BeNil())
					Expect(updatedFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(dimension.Name)))
					Expect(updatedFleetAutoscaler.Name).To(ContainSubstring(strings.ToLower(m.Name)))
					Expect(updatedFleet.Validate(testing.FakeAPIHooks{})).To(BeNil())
					Expect(updatedFleetAutoscaler.Validate()).To(BeNil())
				})
			})
		})
	})

	Describe("GetAllDimension", func() {
		var (
			req *emptypb.Empty
		)
		BeforeEach(func() {
			req = &emptypb.Empty{}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockService.EXPECT().FindAllDimensions(gomock.Any()).Return(game.Dimensions{dimension}, nil)
				out, err := server.GetAllDimension(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out.Dimensions).To(HaveLen(1))
			})
		})

		When("given invalid input", func() {
			It("should error for invalid ctx (nil)", func() {
				out, err := server.GetAllDimension(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error for empty ctx (empty)", func() {
				out, err := server.GetAllDimension(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.GetAllDimension(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (player)", func() {
				out, err := server.GetAllDimension(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindAllDimensions error", func() {
				mockService.EXPECT().FindAllDimensions(gomock.Any()).Return(game.Dimensions{}, fakeErr)
				out, err := server.GetAllDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("GetAllMap", func() {
		var (
			req *emptypb.Empty
		)
		BeforeEach(func() {
			req = &emptypb.Empty{}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockService.EXPECT().FindAllMaps(gomock.Any()).Return(game.Maps{m}, nil)
				out, err := server.GetAllMaps(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out.Maps).To(HaveLen(1))
			})
		})

		When("given invalid input", func() {
			It("should error for invalid ctx (nil)", func() {
				out, err := server.GetAllMaps(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error for empty ctx (empty)", func() {
				out, err := server.GetAllMaps(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.GetAllMaps(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (player)", func() {
				out, err := server.GetAllMaps(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindAllMapss error", func() {
				mockService.EXPECT().FindAllMaps(gomock.Any()).Return(game.Maps{}, fakeErr)
				out, err := server.GetAllMaps(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("GetDimension", func() {
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
			It("should work (admin)", func() {
				mockService.EXPECT().FindDimension(gomock.Any(), req).Return(dimension, nil)
				out, err := server.GetDimension(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out.Name).To(BeEquivalentTo(dimension.Name))
			})
		})

		When("given invalid input", func() {
			It("should error for invalid ctx (nil)", func() {
				out, err := server.GetDimension(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error for empty ctx (empty)", func() {
				out, err := server.GetDimension(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.GetDimension(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (player)", func() {
				out, err := server.GetDimension(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindDimension error", func() {
				mockService.EXPECT().FindDimension(gomock.Any(), req).Return(nil, fakeErr)
				out, err := server.GetDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindDimension no results", func() {
				mockService.EXPECT().FindDimension(gomock.Any(), req).Return(nil, nil)
				out, err := server.GetDimension(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("GetMap", func() {
		var (
			req *pb.MapTarget
		)
		BeforeEach(func() {
			req = &pb.MapTarget{
				FindBy: &pb.MapTarget_Id{
					Id: dimension.Id.String(),
				},
			}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockService.EXPECT().FindMap(gomock.Any(), req).Return(m, nil)
				out, err := server.GetMap(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out.Name).To(BeEquivalentTo(m.Name))
			})
		})

		When("given invalid input", func() {
			It("should error for invalid ctx (nil)", func() {
				out, err := server.GetMap(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error for empty ctx (empty)", func() {
				out, err := server.GetMap(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.GetMap(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (player)", func() {
				out, err := server.GetMap(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindMap error", func() {
				mockService.EXPECT().FindMap(gomock.Any(), req).Return(nil, fakeErr)
				out, err := server.GetMap(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for FindMap no results", func() {
				mockService.EXPECT().FindMap(gomock.Any(), req).Return(nil, nil)
				out, err := server.GetMap(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})
})

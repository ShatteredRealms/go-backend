package srv_test

import (
	"context"
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
)

var _ = Describe("Connection server (local)", func() {
	var (
		hook           *test.Hook
		mockController *gomock.Controller
		ctx            context.Context

		globalConfig   *config.GlobalConfig
		conf           *app.GameBackendServerContext
		mockCharClient *mocks.MockCharacterServiceClient
		mockChatClient *mocks.MockChatServiceClient
		mockService    *mocks.MockGamebackendService
		server         pb.ConnectionServiceServer

		character   *model.Character
		pendingConn *model.PendingConnection
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
			Tracer:             otel.Tracer("test-connection"),
		}

		server, err = srv.NewConnectionServiceServer(ctx, conf)
		Expect(err).NotTo(HaveOccurred())
		Expect(server).NotTo(BeNil())

		character = &model.Character{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: 0,
			OwnerId:   *player.ID,
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

		id := uuid.New()
		pendingConn = &model.PendingConnection{
			Id:         &id,
			Character:  character.Name,
			ServerName: faker.Username(),
			CreatedAt:  time.Now(),
		}

		hook.Reset()
	})

	Describe("ConnectGameServer", func() {
		var (
			req *pb.CharacterTarget
		)
		BeforeEach(func() {
			req = &pb.CharacterTarget{
				Type: &pb.CharacterTarget_Id{
					Id: uint64(character.ID),
				},
			}
		})

		When("given valid input", func() {
			It("should work (admin)", func() {
				mockCharClient.EXPECT().GetCharacter(gomock.Any(), req).Return(character.ToPb(), nil)
				mockService.EXPECT().CreatePendingConnection(gomock.Any(), character.Name, "localhost").Return(pendingConn, nil)
				out, err := server.ConnectGameServer(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out.Address).To(Equal("127.0.0.1"))
				Expect(out.Port).To(BeEquivalentTo(7777))
				Expect(out.ConnectionId).To(Equal(pendingConn.Id.String()))
			})

			It("should work (player)", func() {
				mockCharClient.EXPECT().GetCharacter(gomock.Any(), req).Return(character.ToPb(), nil)
				mockService.EXPECT().CreatePendingConnection(gomock.Any(), character.Name, "localhost").Return(pendingConn, nil)
				out, err := server.ConnectGameServer(incPlayerCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out.Address).To(Equal("127.0.0.1"))
				Expect(out.Port).To(BeEquivalentTo(7777))
				Expect(out.ConnectionId).To(Equal(pendingConn.Id.String()))
			})
		})

		When("given invalid input", func() {
			It("should error for invalid context", func() {
				out, err := server.ConnectGameServer(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for empty context", func() {
				out, err := server.ConnectGameServer(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error for invalid permission (guest)", func() {
				out, err := server.ConnectGameServer(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error on get character err", func() {
				mockCharClient.EXPECT().GetCharacter(gomock.Any(), req).Return(character.ToPb(), fakeErr)
				out, err := server.ConnectGameServer(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error on character not found", func() {
				mockCharClient.EXPECT().GetCharacter(gomock.Any(), req).Return(nil, nil)
				out, err := server.ConnectGameServer(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error on creating pending connection error", func() {
				mockCharClient.EXPECT().GetCharacter(gomock.Any(), req).Return(character.ToPb(), nil)
				mockService.EXPECT().CreatePendingConnection(gomock.Any(), character.Name, "localhost").Return(nil, fakeErr)
				out, err := server.ConnectGameServer(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})
})

package srv_test

import (
	"context"
	"time"

	app "github.com/ShatteredRealms/go-backend/cmd/character/app"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.opentelemetry.io/otel"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Character server", func() {
	var (
		hook            *test.Hook
		mockController  *gomock.Controller
		mockCharService *mocks.MockCharacterService
		mockInvService  *mocks.MockInventoryService
		charCtx         *app.CharactersServerContext

		server pb.CharacterServiceServer
		ctx    = context.Background()

		character *model.Character
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
		mockController = gomock.NewController(GinkgoT())

		mockCharService = mocks.NewMockCharacterService(mockController)
		mockInvService = mocks.NewMockInventoryService(mockController)

		charCtx = &app.CharactersServerContext{
			GlobalConfig:     conf,
			CharacterService: mockCharService,
			InventoryService: mockInvService,
			KeycloakClient:   keycloak,
			Tracer:           otel.Tracer("test-character"),
		}

		var err error
		server, err = srv.NewCharacterServiceServer(ctx, charCtx)
		Expect(err).NotTo(HaveOccurred())
		Expect(server).NotTo(BeNil())

		character = &model.Character{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: 0,
			OwnerId:   faker.Username(),
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

		hook.Reset()
	})

	Describe("AddCharacterPlayTime", func() {
		When("given valid input", func() {
			It("should work", func() {
				_ = character
				// req := &pb.AddPlayTimeRequest{
				// 	Character: &pb.CharacterTarget{
				// 		Type: &pb.CharacterTarget_Id{Id: uint64(character.ID)},
				// 	},
				// 	Time: 100,
				// }
				// mockCharService.EXPECT().FindByName(gomock.Any(), character.Name).Return(character, nil)
				// mockCharService.EXPECT().AddPlayTime(gomock.Any(), character.ID, req.Time).Return(uint64(200), nil)
				// out, err := server.AddCharacterPlayTime(incClientCtx, req)
				// Expect(err).NotTo(HaveOccurred())
				// Expect(out).To(BeEquivalentTo(200))
			})
		})
		When("given invalid input", func() {
			It("should error", func() {

			})
		})
	})

	Describe("AddCharacterPlayTime", func() {
		When("given valid input", func() {
			It("should work", func() {

			})
		})
		When("given invalid input", func() {
			It("should error", func() {

			})
		})
	})
})

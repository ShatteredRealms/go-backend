package srv_test

import (
	"context"

	character "github.com/ShatteredRealms/go-backend/cmd/character/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.opentelemetry.io/otel"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Character server", func() {
	var (
		hook           *test.Hook
		mockController *gomock.Controller
		charCtx        *character.CharactersServerContext

		server pb.CharacterServiceServer
		ctx    = context.Background()
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
		mockController = gomock.NewController(GinkgoT())

		conf := config.NewGlobalConfig(ctx)
		mockCharService := mocks.NewMockCharacterService(mockController)
		mockInvService := mocks.NewMockInventoryService(mockController)

		charCtx = &character.CharactersServerContext{
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

		hook.Reset()
	})
	It("should work", func() {
		Expect(keycloak).NotTo(BeNil())
	})
})

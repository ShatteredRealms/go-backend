package service_test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/sirupsen/logrus/hooks/test"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

var _ = Describe("Chat service", Ordered, func() {

	var (
		hook *test.Hook

		mockController *gomock.Controller
		mockRepository *mocks.MockChatRepository
		cleanupFunc    func()
		kResource      *dockertest.Resource

		chatService service.ChatService

		kafkaPort uint

		err       error
		fakeError = fmt.Errorf("error")
		channels  = []*model.ChatChannel{
			{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Name:      faker.Username(),
				Dimension: faker.Username(),
			},
			{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: time.Now().Add(time.Minute),
					UpdatedAt: time.Now().Add(time.Minute),
				},
				Name:      faker.Username(),
				Dimension: faker.Username(),
			},
		}
	)

	BeforeAll(func() {
		_, hook = test.NewNullLogger()
		mockController = gomock.NewController(GinkgoT())
		mockRepository = mocks.NewMockChatRepository(mockController)

		cleanupFunc, _, kResource = testdb.SetupKafkaWithDocker()
		sPort := kResource.GetPort("29093/tcp")
		port, err := strconv.ParseUint(sPort, 10, 64)
		kafkaPort = uint(port)
		Expect(err).NotTo(HaveOccurred())
		hook.Reset()
	})

	Describe("NewChatService", Ordered, func() {
		When("given invalid input", func() {
			It("should fail if migration fails", func() {
				mockRepository.EXPECT().Migrate(gomock.Any()).Return(fakeError)
				chatService, err = service.NewChatService(context.Background(), mockRepository, config.ServerAddress{
					Port: kafkaPort,
					Host: "localhost",
				})

				Expect(err).To(MatchError(fakeError))
				Expect(chatService).To(BeNil())
			})

			It("should fail if kafka connect fails", func() {
				mockRepository.EXPECT().Migrate(gomock.Any()).Return(nil)
				chatService, err = service.NewChatService(context.Background(), mockRepository, config.ServerAddress{
					Port: 0,
					Host: "nowhere",
				})

				Expect(err).To(HaveOccurred())
				Expect(chatService).To(BeNil())
			})

			It("should fail if kafka all channels", func() {
				mockRepository.EXPECT().Migrate(gomock.Any()).Return(nil)
				mockRepository.EXPECT().AllChannels(gomock.Any()).Return(model.ChatChannels{}, fakeError)
				chatService, err = service.NewChatService(context.Background(), mockRepository, config.ServerAddress{
					Port: kafkaPort,
					Host: "localhost",
				})

				Expect(err).To(HaveOccurred())
				Expect(chatService).To(BeNil())
			})
		})

		When("valid input is given", func() {
			It("should succeed", func() {
				mockRepository.EXPECT().Migrate(gomock.Any()).Return(nil)
				mockRepository.EXPECT().AllChannels(gomock.Any()).Return(channels, nil)
				chatService, err = service.NewChatService(context.Background(), mockRepository, config.ServerAddress{
					Port: kafkaPort,
					Host: "localhost",
				})

				Expect(err).To(BeNil())
				Expect(chatService).NotTo(BeNil())
			})
		})
	})

	AfterAll(func() {
		cleanupFunc()
	})
})

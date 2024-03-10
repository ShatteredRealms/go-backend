package service_test

import (
	"context"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		log.Logger, hook = test.NewNullLogger()
		mockController = gomock.NewController(GinkgoT())
		mockRepository = mocks.NewMockChatRepository(mockController)

		cleanupFunc, kafkaPort = testdb.SetupKafkaWithDocker()
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
				Eventually(func(g Gomega) error {
					mockRepository.EXPECT().Migrate(gomock.Any()).Return(nil)
					mockRepository.EXPECT().AllChannels(gomock.Any()).Return(channels, nil)
					chatService, err = service.NewChatService(context.Background(), mockRepository, config.ServerAddress{
						Port: kafkaPort,
						Host: "localhost",
					})
					return err
				}).Within(time.Minute).Should(Succeed())
			})
		})
	})

	Describe("AllChannels", func() {
		It("should directly call the repo", func() {
			mockRepository.EXPECT().AllChannels(gomock.Any()).Return(channels, fakeError)
			out, err := chatService.AllChannels(nil)
			Expect(err).To(HaveOccurred())
			Expect(out).To(ContainElements(channels))
		})
	})

	Describe("GetChannel", func() {
		It("should directly call the repo", func() {
			mockRepository.EXPECT().FindChannelById(gomock.Any(), uint(1)).Return(channels[0], fakeError)
			out, err := chatService.GetChannel(nil, uint(1))
			Expect(err).To(HaveOccurred())
			Expect(out).To(Equal(channels[0]))
		})
	})

	Describe("CreateChannel", func() {
		When("given invalid data", func() {
			It("should directly call the repo", func() {
				mockRepository.EXPECT().CreateChannel(gomock.Any(), channels[0]).Return(channels[0], fakeError)
				out, err := chatService.CreateChannel(nil, channels[0])
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})

		When("given valid data", func() {
			It("should directly call the repo", func() {
				mockRepository.EXPECT().CreateChannel(gomock.Any(), channels[0]).Return(channels[0], nil)
				out, err := chatService.CreateChannel(nil, channels[0])
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(Equal(channels[0]))
			})
		})
	})

	Describe("DeleteChannel", func() {
		When("given invalid data", func() {
			It("should directly call the repo", func() {
				mockRepository.EXPECT().DeleteChannel(gomock.Any(), channels[0]).Return(fakeError)
				err := chatService.DeleteChannel(nil, channels[0])
				Expect(err).To(HaveOccurred())
			})
		})

		When("given valid data", func() {
			It("should directly call the repo", func() {
				mockRepository.EXPECT().DeleteChannel(gomock.Any(), channels[1]).Return(nil)
				err := chatService.DeleteChannel(nil, channels[1])
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("Sending channel messages", func() {
		It("should work", func() {
			reader := chatService.ChannelMessagesReader(context.Background(), channels[0].ID)
			Expect(reader).NotTo(BeNil())

			user := faker.Username()
			messageA := faker.Email()
			messageB := faker.Email()
			Expect(chatService.SendChannelMessage(context.Background(), user, messageA, channels[0].ID)).To(Succeed())
			Expect(chatService.SendChannelMessage(context.Background(), user, messageB, channels[0].ID)).To(Succeed())
			eventuallyFunc := func(g Gomega) (string, error) {
				message, err := reader.ReadMessage(context.Background())
				Expect(err).NotTo(HaveOccurred())
				return fmt.Sprintf("%s: %s", string(message.Key), string(message.Value)), err
			}
			Eventually(eventuallyFunc).Within(time.Second).Should(Equal(fmt.Sprintf("%s: %s", user, messageA)))
			Eventually(eventuallyFunc).Within(time.Second).Should(Equal(fmt.Sprintf("%s: %s", user, messageB)))
		})

		It("should fail with empty message", func() {
			Expect(chatService.SendChannelMessage(context.Background(), faker.Username(), "", channels[0].ID)).NotTo(Succeed())
		})
	})

	Describe("Sending channel messages", func() {
		It("should work", func() {
			// user := faker.Username()
			// sender := faker.Username() + "a"
			// Expect(chatService.RegisterCharacterChatTopic(context.Background(), user)).To(Succeed())
			// reader := chatService.DirectMessagesReader(context.Background(), user)
			// Expect(reader).NotTo(BeNil())
			//
			// messageA := faker.Email()
			// messageB := faker.Email()
			// Eventually(chatService.SendDirectMessage(context.Background(), sender, messageA, user)).WithTimeout(time.Second * 5).WithPolling(time.Second).Should(Succeed())
			// Eventually(chatService.SendDirectMessage(context.Background(), sender, messageB, user)).Should(Succeed())
			// eventuallyFunc := func(g Gomega) (string, error) {
			// 	message, err := reader.ReadMessage(context.Background())
			// 	Expect(err).NotTo(HaveOccurred())
			// 	return fmt.Sprintf("%s: %s", string(message.Key), string(message.Value)), err
			// }
			// Eventually(eventuallyFunc).Within(time.Second).Should(Equal(fmt.Sprintf("%s: %s", sender, messageA)))
			// Eventually(eventuallyFunc).Within(time.Second).Should(Equal(fmt.Sprintf("%s: %s", sender, messageB)))
		})

		It("should fail with empty message", func() {
			Expect(chatService.SendDirectMessage(context.Background(), faker.Username(), "", faker.Username())).NotTo(Succeed())
		})

		It("should fail if no character exists", func() {
			Expect(chatService.SendDirectMessage(context.Background(), faker.Username(), faker.Email(), faker.Username())).NotTo(Succeed())
		})
	})

	AfterAll(func() {
		cleanupFunc()
	})
})

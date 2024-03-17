package service_test

import (
	"context"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
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
		hook.Reset()
	})

	Describe("NewChatService", Ordered, func() {
		When("given invalid input", func() {
			It("should fail if migration fails", func() {
				mockRepository.EXPECT().Migrate(gomock.Any()).Return(fakeError)
				chatService, err = service.NewChatService(context.Background(), mockRepository, config.ServerAddress{
					Port: kafkaPort,
					Host: "127.0.0.1",
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
					Host: "127.0.0.1",
				})

				Expect(err).To(HaveOccurred())
				Expect(chatService).To(BeNil())
			})
		})

		When("valid input is given", func() {
			It("should succeed", func() {
				Eventually(func(g Gomega) error {
					mockRepository.EXPECT().Migrate(gomock.Any()).Return(nil).AnyTimes()
					mockRepository.EXPECT().AllChannels(gomock.Any()).Return(channels, nil).AnyTimes()
					chatService, err = service.NewChatService(context.Background(), mockRepository, config.ServerAddress{
						Port: kafkaPort,
						Host: "127.0.0.1",
					})
					return err
				}).Within(time.Minute).Should(Succeed())
			})
		})
	})

	Describe("ChangeAuthorizationForCharacter", func() {
		It("should work", func() {
			ctx := context.Background()
			id := uint(1)
			ids := []uint{1, 2, 3}
			mockRepository.EXPECT().ChangeAuthorizationForCharacter(ctx, id, ids, true).Return(fakeError)
			err := chatService.ChangeAuthorizationForCharacter(ctx, id, ids, true)
			Expect(err).To(MatchError(fakeError))
		})
	})

	Describe("AuthorizedChannelsForCharacter", func() {
		It("should work", func() {
			ctx := context.Background()
			id := uint(1)
			mockRepository.EXPECT().AuthorizedChannelsForCharacter(ctx, id).Return(channels, fakeError)
			out, err := chatService.AuthorizedChannelsForCharacter(ctx, id)
			Expect(err).To(MatchError(fakeError))
			Expect(out).To(ContainElements(channels))
		})
	})

	Describe("UpdateChannel", func() {
		When("given valid input", func() {
			It("should work", func() {
				ctx := context.Background()
				req := &pb.UpdateChatChannelRequest{
					ChannelId: uint64(channels[0].ID),
					OptionalName: &pb.UpdateChatChannelRequest_Name{
						Name: faker.Username(),
					},
					OptionalDimension: &pb.UpdateChatChannelRequest_Dimension{
						Dimension: faker.Username(),
					},
				}
				expectedUpdate := &model.ChatChannel{}
				*expectedUpdate = *channels[0]
				expectedUpdate.Name = req.GetName()
				expectedUpdate.Dimension = req.GetDimension()
				mockRepository.EXPECT().FindChannelById(gomock.Any(), channels[0].ID).Return(channels[0], nil)
				mockRepository.EXPECT().UpdateChannel(gomock.Any(), expectedUpdate).Return(expectedUpdate, fakeError)
				out, err := chatService.UpdateChannel(ctx, req)
				Expect(err).To(MatchError(fakeError))
				Expect(out).To(Equal(expectedUpdate))
			})
		})
		When("given valid input", func() {
			It("should error on not found", func() {
				ctx := context.Background()
				req := &pb.UpdateChatChannelRequest{
					ChannelId: uint64(channels[0].ID),
				}
				mockRepository.EXPECT().FindChannelById(gomock.Any(), channels[0].ID).Return(channels[0], fakeError)
				out, err := chatService.UpdateChannel(ctx, req)
				Expect(err).To(MatchError(fakeError))
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("AllChannels", func() {
		It("should directly call the repo", func() {
			mockRepository.EXPECT().AllChannels(gomock.Any()).Return(channels, nil).AnyTimes()
			out, err := chatService.AllChannels(nil)
			Expect(err).NotTo(HaveOccurred())
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
			user := faker.Username()
			sender := faker.Username() + "a"
			Expect(chatService.RegisterCharacterChatTopic(context.Background(), user)).To(Succeed())
			reader := chatService.DirectMessagesReader(context.Background(), user)
			Expect(reader).NotTo(BeNil())

			messageA := faker.Email()
			messageB := faker.Email()
			Eventually(func(g Gomega) error {
				Expect(chatService.RegisterCharacterChatTopic(context.Background(), user)).To(Succeed())
				err := chatService.SendDirectMessage(context.Background(), sender, messageA, user)
				return err
			}).WithTimeout(time.Second * 15).WithPolling(time.Second).Should(Succeed())
			Expect(chatService.SendDirectMessage(context.Background(), sender, messageB, user)).To(Succeed())
			eventuallyFunc := func(g Gomega) (string, error) {
				message, err := reader.ReadMessage(context.Background())
				Expect(err).NotTo(HaveOccurred())
				return fmt.Sprintf("%s: %s", string(message.Key), string(message.Value)), err
			}
			Eventually(eventuallyFunc).Within(time.Second).Should(Equal(fmt.Sprintf("%s: %s", sender, messageA)))
			Eventually(eventuallyFunc).Within(time.Second).Should(Equal(fmt.Sprintf("%s: %s", sender, messageB)))
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

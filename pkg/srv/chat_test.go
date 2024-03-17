package srv_test

import (
	"context"
	"fmt"
	"io"
	"time"

	app "github.com/ShatteredRealms/go-backend/cmd/chat/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus/hooks/test"
	"go.opentelemetry.io/otel"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

var _ = Describe("Chat", func() {
	var (
		hook            *test.Hook
		mockController  *gomock.Controller
		mockCharService *mocks.MockCharacterServiceClient
		mockChatService *mocks.MockChatService
		// mockChatRepo    *mocks.MockChatRepository

		ctx = context.Background()

		chatCtx *app.ChatServerContext
		server  pb.ChatServiceServer

		kafkaConn    *kafka.Conn
		readerConfig = kafka.ReaderConfig{
			Topic:    topicName,
			MinBytes: 1,
			MaxBytes: 10e3,
		}
		kafkaWriter *kafka.Writer
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
		mockController = gomock.NewController(GinkgoT())
		mockCharService = mocks.NewMockCharacterServiceClient(mockController)
		mockChatService = mocks.NewMockChatService(mockController)

		chatCtx = &app.ChatServerContext{
			GlobalConfig:     conf,
			ChatService:      mockChatService,
			CharacterService: mockCharService,
			KeycloakClient:   keycloak,
			Tracer:           otel.Tracer("test-character"),
		}

		var err error
		server, err = srv.NewChatServiceServer(ctx, chatCtx)
		Expect(err).NotTo(HaveOccurred())
		Expect(server).NotTo(BeNil())

		hook.Reset()
	})

	Context("with kafka", func() {
		var character *model.Character
		var chatChannel *model.ChatChannel
		BeforeEach(func() {
			Eventually(func(g Gomega) error {
				var err error
				kafkaConn, err = repository.ConnectKafka(config.ServerAddress{
					Port: uint(kafkaPort),
					Host: "127.0.0.1",
				})
				if err != nil {
					return err
				}
				Expect(kafkaConn).NotTo(BeNil())

				readerConfig.Brokers = []string{fmt.Sprintf("127.0.0.1:%d", kafkaPort)}
				kafkaWriter = &kafka.Writer{
					Addr:     kafkaConn.RemoteAddr(),
					Topic:    topicName,
					Balancer: &kafka.LeastBytes{},
					Async:    true,
				}

				return kafkaConn.CreateTopics(kafka.TopicConfig{
					Topic:             topicName,
					NumPartitions:     1,
					ReplicationFactor: 1,
				})
			}).Within(time.Minute).Should(Succeed())
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
			chatChannel = &model.ChatChannel{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Name:      faker.Username(),
				Dimension: faker.Username(),
			}
		})

		Describe("ConnectChannel", func() {
			var req *pb.ChatChannelTarget
			var mockInSrv *mocks.MockChatService_ConnectChannelServer
			var msg *pb.ChatMessage
			BeforeEach(func() {
				req = &pb.ChatChannelTarget{
					Id: 1,
				}
				mockInSrv = mocks.NewMockChatService_ConnectChannelServer(mockController)
				msg = &pb.ChatMessage{
					Message:       faker.Username(),
					CharacterName: faker.Username(),
				}
			})

			When("given valid input", func() {
				It("should work for users with chat manager permissions (admin)", func() {
					mockChatService.EXPECT().ChannelMessagesReader(gomock.Any(), uint(req.Id)).Return(kafka.NewReader(readerConfig))
					mockInSrv.EXPECT().Context().Return(incAdminCtx).AnyTimes()
					mockInSrv.EXPECT().Send(gomock.Any()).Return(io.EOF)
					Eventually(kafkaWriter.WriteMessages(context.Background(), kafka.Message{
						Key:   []byte(msg.CharacterName),
						Value: []byte(msg.Message),
					})).Within(time.Second * 5).Should(Succeed())
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(io.EOF))
				})

				It("should work for users with chat permissions (player)", func() {
					mockChatService.EXPECT().ChannelMessagesReader(gomock.Any(), uint(req.Id)).Return(kafka.NewReader(readerConfig))
					mockCharService.EXPECT().GetAllCharactersForUser(gomock.Any(), gomock.Any()).Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().AuthorizedChannelsForCharacter(gomock.Any(), character.ID).Return(model.ChatChannels{chatChannel}, nil)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					mockInSrv.EXPECT().Send(gomock.Any()).Return(io.EOF)
					Eventually(kafkaWriter.WriteMessages(context.Background(), kafka.Message{
						Key:   []byte(msg.CharacterName),
						Value: []byte(msg.Message),
					})).Within(time.Second * 5).Should(Succeed())
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(io.EOF))
				})
			})

			When("given invalid input", func() {
				It("should error on invalid claims", func() {
					mockInSrv.EXPECT().Context().Return(nil).AnyTimes()
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(model.ErrUnauthorized))
				})

				It("should error on no permissions (guest)", func() {
					mockInSrv.EXPECT().Context().Return(incGuestCtx).AnyTimes()
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(model.ErrUnauthorized))
				})

				It("should error on no permissions for chat channel (player)", func() {
					mockCharService.EXPECT().GetAllCharactersForUser(gomock.Any(), gomock.Any()).Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().AuthorizedChannelsForCharacter(gomock.Any(), character.ID).Return(model.ChatChannels{}, nil)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(model.ErrUnauthorized))
				})

				It("should error on no permissions for chat channel due to no characters (player)", func() {
					mockCharService.EXPECT().GetAllCharactersForUser(gomock.Any(), gomock.Any()).Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{}}, nil)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(model.ErrUnauthorized))
				})

				It("should error if getting characters has errors", func() {
					mockCharService.EXPECT().GetAllCharactersForUser(gomock.Any(), gomock.Any()).Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{}}, fakeErr)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(model.ErrHandleRequest))
				})

				It("should error if getting authorized channels for character has errors", func() {
					mockCharService.EXPECT().GetAllCharactersForUser(gomock.Any(), gomock.Any()).Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().AuthorizedChannelsForCharacter(gomock.Any(), character.ID).Return(nil, fakeErr)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(model.ErrHandleRequest))
				})
			})
		})

		Describe("ConnectDirectMessage", func() {
			var req *pb.CharacterTarget
			var mockInSrv *mocks.MockChatService_ConnectDirectMessageServer
			var msg *pb.ChatMessage
			BeforeEach(func() {
				req = &pb.CharacterTarget{
					Type: &pb.CharacterTarget_Id{Id: uint64(character.ID)},
				}
				mockInSrv = mocks.NewMockChatService_ConnectDirectMessageServer(mockController)
				msg = &pb.ChatMessage{
					Message:       faker.Username(),
					CharacterName: faker.Username(),
				}
			})

			When("given valid input", func() {
				It("should work for users with chat manager permissions (admin)", func() {
					mockChatService.EXPECT().DirectMessagesReader(gomock.Any(), character.Name).Return(kafka.NewReader(readerConfig))
					mockCharService.EXPECT().GetAllCharactersForUser(gomock.Any(), gomock.Any()).Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockInSrv.EXPECT().Context().Return(incAdminCtx).AnyTimes()
					mockInSrv.EXPECT().Send(gomock.Any()).Return(io.EOF)
					Eventually(kafkaWriter.WriteMessages(context.Background(), kafka.Message{
						Key:   []byte(msg.CharacterName),
						Value: []byte(msg.Message),
					})).Should(Succeed())
					Expect(server.ConnectDirectMessage(req, mockInSrv)).To(MatchError(io.EOF))
				})

				It("should work for users with chat permissions (player)", func() {
					mockChatService.EXPECT().DirectMessagesReader(gomock.Any(), character.Name).Return(kafka.NewReader(readerConfig))
					mockCharService.EXPECT().GetAllCharactersForUser(gomock.Any(), gomock.Any()).Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					mockInSrv.EXPECT().Send(gomock.Any()).Return(io.EOF)
					Eventually(kafkaWriter.WriteMessages(context.Background(), kafka.Message{
						Key:   []byte(msg.CharacterName),
						Value: []byte(msg.Message),
					})).Within(time.Second * 5).Should(Succeed())
					Expect(server.ConnectDirectMessage(req, mockInSrv)).To(MatchError(io.EOF))
				})
			})

			When("given invalid input", func() {
				It("should error on invalid claims", func() {
					mockInSrv.EXPECT().Context().Return(nil).AnyTimes()
					Expect(server.ConnectDirectMessage(req, mockInSrv)).To(MatchError(model.ErrUnauthorized))
				})

				It("should error on no permissions (guest)", func() {
					mockInSrv.EXPECT().Context().Return(incGuestCtx).AnyTimes()
					Expect(server.ConnectDirectMessage(req, mockInSrv)).To(MatchError(model.ErrUnauthorized))
				})

				It("should error if not owner of character (player)", func() {
					mockCharService.EXPECT().GetAllCharactersForUser(gomock.Any(), gomock.Any()).Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{}}, nil)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectDirectMessage(req, mockInSrv)).To(MatchError(model.ErrUnauthorized))
				})

				It("should error if not owner of character (adin)", func() {
					mockCharService.EXPECT().GetAllCharactersForUser(gomock.Any(), gomock.Any()).Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{}}, nil)
					mockInSrv.EXPECT().Context().Return(incAdminCtx).AnyTimes()
					Expect(server.ConnectDirectMessage(req, mockInSrv)).To(MatchError(model.ErrUnauthorized))
				})

				It("should error if getting characters has errors", func() {
					mockCharService.EXPECT().GetAllCharactersForUser(gomock.Any(), gomock.Any()).Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{}}, fakeErr)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectDirectMessage(req, mockInSrv)).NotTo(Succeed())
				})

				It("should error on invalid target type", func() {
					req.Type = nil
					mockCharService.EXPECT().GetAllCharactersForUser(gomock.Any(), gomock.Any()).Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{}}, fakeErr)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectDirectMessage(req, mockInSrv)).NotTo(Succeed())
				})
			})
		})

		Describe("SendChatMessage", func() {
			var (
				req *pb.SendChatMessageRequest
			)
			BeforeEach(func() {
				req = &pb.SendChatMessageRequest{
					ChannelId: uint64(chatChannel.ID),
					ChatMessage: &pb.ChatMessage{
						Message:       faker.Username(),
						CharacterName: character.Name,
					},
				}
			})
			When("given valid input", func() {
				It("should work (admin)", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().
						SendChannelMessage(gomock.Any(), gomock.Any(), gomock.Any(), uint(req.ChannelId)).
						Return(nil)
					out, err := server.SendChatMessage(incAdminCtx, req)
					Expect(err).NotTo(HaveOccurred())
					Expect(out).NotTo(BeNil())
				})

				It("should work (player)", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().
						AuthorizedChannelsForCharacter(gomock.Any(), character.ID).
						Return(model.ChatChannels{chatChannel}, nil)
					mockChatService.EXPECT().
						SendChannelMessage(gomock.Any(), gomock.Any(), gomock.Any(), uint(req.ChannelId)).
						Return(nil)
					out, err := server.SendChatMessage(incPlayerCtx, req)
					Expect(err).NotTo(HaveOccurred())
					Expect(out).NotTo(BeNil())
				})
			})

			When("given invalid input", func() {
				It("should err on invalid context", func() {
					out, err := server.SendChatMessage(nil, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())

					out, err = server.SendChatMessage(context.Background(), req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err on invalid permission (guest)", func() {
					out, err := server.SendChatMessage(incGuestCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if not owner (admin)", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{}}, nil)
					out, err := server.SendChatMessage(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if not owner (player)", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{}}, nil)
					out, err := server.SendChatMessage(incPlayerCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if no channel permission (player)", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().
						AuthorizedChannelsForCharacter(gomock.Any(), character.ID).
						Return(model.ChatChannels{}, nil)
					out, err := server.SendChatMessage(incPlayerCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if verify user errors", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(nil, fakeErr)
					out, err := server.SendChatMessage(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if getting authorized channels fails", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().
						AuthorizedChannelsForCharacter(gomock.Any(), character.ID).
						Return(nil, fakeErr)
					out, err := server.SendChatMessage(incPlayerCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should error if sending message fails", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().
						SendChannelMessage(gomock.Any(), gomock.Any(), gomock.Any(), uint(req.ChannelId)).
						Return(fakeErr)
					out, err := server.SendChatMessage(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})
			})
		})

		Describe("SendDirectMessage", func() {
			var (
				req *pb.SendDirectMessageRequest
			)
			BeforeEach(func() {
				req = &pb.SendDirectMessageRequest{
					Target: &pb.CharacterTarget{
						Type: &pb.CharacterTarget_Name{
							Name: character.Name,
						},
					},
					ChatMessage: &pb.ChatMessage{Message: faker.Username(), CharacterName: character.Name},
				}
			})
			When("given valid input", func() {
				It("should work (admin)", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().
						SendDirectMessage(gomock.Any(), gomock.Any(), gomock.Any(), character.Name).
						Return(nil)
					out, err := server.SendDirectMessage(incAdminCtx, req)
					Expect(err).NotTo(HaveOccurred())
					Expect(out).NotTo(BeNil())
				})

				It("should work (player)", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().
						SendDirectMessage(gomock.Any(), gomock.Any(), gomock.Any(), character.Name).
						Return(nil)
					out, err := server.SendDirectMessage(incPlayerCtx, req)
					Expect(err).NotTo(HaveOccurred())
					Expect(out).NotTo(BeNil())
				})
			})

			When("given invalid input", func() {
				It("should err on invalid context", func() {
					out, err := server.SendDirectMessage(nil, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())

					out, err = server.SendDirectMessage(context.Background(), req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err on invalid permission (guest)", func() {
					out, err := server.SendDirectMessage(incGuestCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if not owner (admin)", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{}}, nil)
					out, err := server.SendDirectMessage(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if not owner (player)", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{}}, nil)
					out, err := server.SendDirectMessage(incPlayerCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if verify user errors", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(nil, fakeErr)
					out, err := server.SendDirectMessage(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should error if sending message fails", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().
						SendDirectMessage(gomock.Any(), gomock.Any(), gomock.Any(), character.Name).
						Return(fakeErr)
					out, err := server.SendDirectMessage(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})
			})
		})
	})
})

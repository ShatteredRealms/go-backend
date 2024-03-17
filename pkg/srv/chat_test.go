package srv_test

import (
	"context"
	"fmt"
	"io"
	"time"

	app "github.com/ShatteredRealms/go-backend/cmd/chat/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
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
	"google.golang.org/protobuf/types/known/emptypb"
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

		chatChannel *model.ChatChannel
		character   *model.Character
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

		character = &model.Character{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: 0,
			OwnerId:   *player.ID,
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

		hook.Reset()
	})

	Context("with kafka", func() {
		var (
			msg              *pb.ChatMessage
			writeMessageFunc = func(g Gomega) error {
				err := kafkaConn.CreateTopics(kafka.TopicConfig{
					Topic:             topicName,
					NumPartitions:     1,
					ReplicationFactor: 1,
				})
				if err != nil {
					return err
				}
				return kafkaWriter.WriteMessages(context.Background(), kafka.Message{
					Key:   []byte(msg.CharacterName),
					Value: []byte(msg.Message),
				})
			}
		)
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
			msg = &pb.ChatMessage{
				Message:       faker.Username(),
				CharacterName: character.Name,
			}
		})

		Describe("ConnectChannel", func() {
			var req *pb.ChatChannelTarget
			var mockInSrv *mocks.MockChatService_ConnectChannelServer
			BeforeEach(func() {
				req = &pb.ChatChannelTarget{
					Id: 1,
				}
				mockInSrv = mocks.NewMockChatService_ConnectChannelServer(mockController)
			})

			When("given valid input", func() {
				It("should work for users with chat manager permissions (admin)", func() {
					mockChatService.EXPECT().ChannelMessagesReader(gomock.Any(), uint(req.Id)).Return(kafka.NewReader(readerConfig))
					mockInSrv.EXPECT().Context().Return(incAdminCtx).AnyTimes()
					mockInSrv.EXPECT().Send(gomock.Any()).Return(io.EOF)
					Eventually(writeMessageFunc).Within(time.Second * 15).Should(Succeed())
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(io.EOF))
				})

				It("should work for users with chat permissions (player)", func() {
					mockChatService.EXPECT().ChannelMessagesReader(gomock.Any(), uint(req.Id)).Return(kafka.NewReader(readerConfig))
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().AuthorizedChannelsForCharacter(gomock.Any(), character.ID).Return(model.ChatChannels{chatChannel}, nil)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					mockInSrv.EXPECT().Send(gomock.Any()).Return(io.EOF)
					Eventually(writeMessageFunc).Within(time.Second * 15).Should(Succeed())
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
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().AuthorizedChannelsForCharacter(gomock.Any(), character.ID).Return(model.ChatChannels{}, nil)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(model.ErrUnauthorized))
				})

				It("should error if getting characters has errors", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(nil, fakeErr)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(model.ErrHandleRequest))
				})

				It("should error if getting authorized channels for character has errors", func() {
					mockCharService.EXPECT().
						GetAllCharactersForUser(gomock.Any(), gomock.Any()).
						Return(&pb.CharactersDetails{Characters: []*pb.CharacterDetails{character.ToPb()}}, nil)
					mockChatService.EXPECT().AuthorizedChannelsForCharacter(gomock.Any(), character.ID).Return(nil, fakeErr)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectChannel(req, mockInSrv)).To(MatchError(model.ErrHandleRequest))
				})
			})
		})

		Describe("ConnectDirectMessage", func() {
			var req *pb.CharacterTarget
			var mockInSrv *mocks.MockChatService_ConnectDirectMessageServer
			BeforeEach(func() {
				req = &pb.CharacterTarget{
					Type: &pb.CharacterTarget_Id{Id: uint64(character.ID)},
				}
				mockInSrv = mocks.NewMockChatService_ConnectDirectMessageServer(mockController)
			})

			When("given valid input", func() {
				It("should work for users with chat manager permissions (admin)", func() {
					character.OwnerId = *admin.ID
					mockChatService.EXPECT().DirectMessagesReader(gomock.Any(), character.Name).Return(kafka.NewReader(readerConfig))
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					mockInSrv.EXPECT().Context().Return(incAdminCtx).AnyTimes()
					mockInSrv.EXPECT().Send(gomock.Any()).Return(io.EOF)
					Eventually(writeMessageFunc).Within(time.Second * 15).Should(Succeed())
					Expect(server.ConnectDirectMessage(req, mockInSrv)).To(MatchError(io.EOF))
				})

				It("should work for users with chat manager permissions (admin other)", func() {
					mockChatService.EXPECT().DirectMessagesReader(gomock.Any(), character.Name).Return(kafka.NewReader(readerConfig))
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					mockInSrv.EXPECT().Context().Return(incAdminCtx).AnyTimes()
					mockInSrv.EXPECT().Send(gomock.Any()).Return(io.EOF)
					Eventually(writeMessageFunc).Within(time.Second * 15).Should(Succeed())
					Expect(server.ConnectDirectMessage(req, mockInSrv)).To(MatchError(io.EOF))
				})

				It("should work for users with chat permissions (player)", func() {
					_, claims, err := helpers.VerifyClaims(incPlayerCtx, keycloak, conf.Chat.Keycloak.Realm)
					Expect(err).NotTo(HaveOccurred())
					Expect(claims.Subject).To(Equal(character.OwnerId))
					mockChatService.EXPECT().DirectMessagesReader(gomock.Any(), character.Name).Return(kafka.NewReader(readerConfig))
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					mockInSrv.EXPECT().Send(gomock.Any()).Return(io.EOF)
					Eventually(writeMessageFunc).Within(time.Second * 15).Should(Succeed())
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
					character.OwnerId = *admin.ID
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectDirectMessage(req, mockInSrv)).To(MatchError(model.ErrUnauthorized))
				})

				It("should error if getting characters has errors", func() {
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(nil, fakeErr)
					mockInSrv.EXPECT().Context().Return(incPlayerCtx).AnyTimes()
					Expect(server.ConnectDirectMessage(req, mockInSrv)).NotTo(Succeed())
				})

				It("should error on invalid target type", func() {
					req.Type = nil
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(nil, fakeErr)
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
					ChannelId:   uint64(chatChannel.ID),
					ChatMessage: msg,
				}
			})
			When("given valid input", func() {
				It("should work (admin)", func() {
					character.OwnerId = *admin.ID
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					mockChatService.EXPECT().
						SendChannelMessage(gomock.Any(), gomock.Any(), gomock.Any(), uint(req.ChannelId)).
						Return(nil)
					out, err := server.SendChatMessage(incAdminCtx, req)
					Expect(err).NotTo(HaveOccurred())
					Expect(out).NotTo(BeNil())
				})

				It("should work (player)", func() {
					Expect(character.OwnerId).To(Equal(*player.ID))
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
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
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					out, err := server.SendChatMessage(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if not owner (player)", func() {
					character.OwnerId = *admin.ID
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					out, err := server.SendChatMessage(incPlayerCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if no channel permission (player)", func() {
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					mockChatService.EXPECT().
						AuthorizedChannelsForCharacter(gomock.Any(), character.ID).
						Return(model.ChatChannels{}, nil)
					out, err := server.SendChatMessage(incPlayerCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if verify user errors", func() {
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(nil, fakeErr)
					out, err := server.SendChatMessage(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if getting authorized channels fails", func() {
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					mockChatService.EXPECT().
						AuthorizedChannelsForCharacter(gomock.Any(), character.ID).
						Return(nil, fakeErr)
					out, err := server.SendChatMessage(incPlayerCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should error if sending message fails", func() {
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					mockChatService.EXPECT().
						AuthorizedChannelsForCharacter(gomock.Any(), character.ID).
						Return(model.ChatChannels{chatChannel}, nil)
					mockChatService.EXPECT().
						SendChannelMessage(gomock.Any(), gomock.Any(), gomock.Any(), uint(req.ChannelId)).
						Return(fakeErr)
					out, err := server.SendChatMessage(incPlayerCtx, req)
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
					ChatMessage: msg,
				}
			})
			When("given valid input", func() {
				It("should work (admin)", func() {
					character.OwnerId = *admin.ID
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					mockChatService.EXPECT().
						SendDirectMessage(gomock.Any(), gomock.Any(), gomock.Any(), character.Name).
						Return(nil)
					out, err := server.SendDirectMessage(incAdminCtx, req)
					Expect(err).NotTo(HaveOccurred())
					Expect(out).NotTo(BeNil())
				})

				It("should work (player)", func() {
					character.OwnerId = *player.ID
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
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
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					out, err := server.SendDirectMessage(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if not owner (player)", func() {
					character.OwnerId = *admin.ID
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					out, err := server.SendDirectMessage(incPlayerCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should err if verify user errors", func() {
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(nil, fakeErr)
					out, err := server.SendDirectMessage(incAdminCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})

				It("should error if sending message fails", func() {
					mockCharService.EXPECT().
						GetCharacter(gomock.Any(), gomock.Any()).
						Return(character.ToPb(), nil)
					mockChatService.EXPECT().
						SendDirectMessage(gomock.Any(), gomock.Any(), gomock.Any(), character.Name).
						Return(fakeErr)
					out, err := server.SendDirectMessage(incPlayerCtx, req)
					Expect(err).To(HaveOccurred())
					Expect(out).To(BeNil())
				})
			})
		})
	})

	Describe("GetChannel", func() {
		var (
			req *pb.ChatChannelTarget
		)
		BeforeEach(func() {
			req = &pb.ChatChannelTarget{
				Id: uint64(chatChannel.ID),
			}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockChatService.EXPECT().GetChannel(gomock.Any(), chatChannel.ID).Return(chatChannel, nil)
				out, err := server.GetChannel(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(BeEquivalentTo(chatChannel.ToPb()))
			})
		})

		When("given invalid input", func() {
			It("should error if invalid context (nil)", func() {
				out, err := server.GetChannel(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid context (no claims)", func() {
				out, err := server.GetChannel(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (player)", func() {
				out, err := server.GetChannel(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (guest)", func() {
				out, err := server.GetChannel(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if getting channel errors", func() {
				mockChatService.EXPECT().GetChannel(gomock.Any(), chatChannel.ID).Return(nil, fakeErr)
				out, err := server.GetChannel(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if no channel exists", func() {
				mockChatService.EXPECT().GetChannel(gomock.Any(), chatChannel.ID).Return(nil, nil)
				out, err := server.GetChannel(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("CreateChannel", func() {
		var (
			req *pb.CreateChannelMessage
		)
		BeforeEach(func() {
			req = &pb.CreateChannelMessage{
				Name:      faker.Username(),
				Dimension: faker.Username(),
			}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockChatService.EXPECT().
					CreateChannel(gomock.Any(), gomock.Eq(&model.ChatChannel{Name: req.Name, Dimension: req.Dimension})).
					Return(nil, nil)
				out, err := server.CreateChannel(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error if invalid context (nil)", func() {
				out, err := server.CreateChannel(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid context (no claims)", func() {
				out, err := server.CreateChannel(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (player)", func() {
				out, err := server.CreateChannel(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (guest)", func() {
				out, err := server.CreateChannel(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if error during creation", func() {
				mockChatService.EXPECT().
					CreateChannel(gomock.Any(), gomock.Eq(&model.ChatChannel{Name: req.Name, Dimension: req.Dimension})).
					Return(nil, fakeErr)
				out, err := server.CreateChannel(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should unique error if name is taken", func() {
				mockChatService.EXPECT().
					CreateChannel(gomock.Any(), gomock.Eq(&model.ChatChannel{Name: req.Name, Dimension: req.Dimension})).
					Return(nil, gorm.ErrDuplicatedKey)
				out, err := server.CreateChannel(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("DeleteChannel", func() {
		var (
			req *pb.ChatChannelTarget
		)
		BeforeEach(func() {
			req = &pb.ChatChannelTarget{
				Id: uint64(chatChannel.ID),
			}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockChatService.EXPECT().
					DeleteChannel(gomock.Any(), gomock.Any()).
					Return(nil)
				out, err := server.DeleteChannel(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error if invalid context (nil)", func() {
				out, err := server.DeleteChannel(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid context (no claims)", func() {
				out, err := server.DeleteChannel(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (player)", func() {
				out, err := server.DeleteChannel(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (guest)", func() {
				out, err := server.DeleteChannel(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if error during creation", func() {
				mockChatService.EXPECT().
					DeleteChannel(gomock.Any(), gomock.Any()).
					Return(fakeErr)
				out, err := server.DeleteChannel(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should unique error if not found is taken", func() {
				mockChatService.EXPECT().
					DeleteChannel(gomock.Any(), gomock.Any()).
					Return(gorm.ErrRecordNotFound)
				out, err := server.DeleteChannel(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("EditChannel", func() {
		var (
			req *pb.UpdateChatChannelRequest
		)
		BeforeEach(func() {
			req = &pb.UpdateChatChannelRequest{
				ChannelId:         uint64(chatChannel.ID),
				OptionalName:      &pb.UpdateChatChannelRequest_Name{Name: faker.Username()},
				OptionalDimension: &pb.UpdateChatChannelRequest_Dimension{Dimension: faker.Username()},
			}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockChatService.EXPECT().
					UpdateChannel(gomock.Any(), req).
					Return(chatChannel, nil)
				out, err := server.EditChannel(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error if invalid context (nil)", func() {
				out, err := server.EditChannel(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid context (no claims)", func() {
				out, err := server.EditChannel(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (player)", func() {
				out, err := server.EditChannel(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (guest)", func() {
				out, err := server.EditChannel(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if error during creation", func() {
				mockChatService.EXPECT().
					UpdateChannel(gomock.Any(), req).
					Return(nil, fakeErr)
				out, err := server.EditChannel(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should unique error if not found is taken", func() {
				mockChatService.EXPECT().
					UpdateChannel(gomock.Any(), req).
					Return(nil, gorm.ErrRecordNotFound)
				out, err := server.EditChannel(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("AllChatChannel", func() {
		var (
			req *emptypb.Empty
		)
		BeforeEach(func() {
			req = &emptypb.Empty{}
		})

		When("given valid input", func() {
			It("should work (admin)", func() {
				mockChatService.EXPECT().
					AllChannels(gomock.Any()).
					Return(model.ChatChannels{chatChannel}, nil)
				out, err := server.AllChatChannels(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error if invalid context (nil)", func() {
				out, err := server.AllChatChannels(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid context (no claims)", func() {
				out, err := server.AllChatChannels(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (player)", func() {
				out, err := server.AllChatChannels(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (guest)", func() {
				out, err := server.AllChatChannels(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if error during creation", func() {
				mockChatService.EXPECT().
					AllChannels(gomock.Any()).
					Return(nil, fakeErr)
				out, err := server.AllChatChannels(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("GetAuthorizedChatChannel", func() {
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
			It("should work (admin self)", func() {
				character.OwnerId = *admin.ID
				mockChatService.EXPECT().
					AuthorizedChannelsForCharacter(gomock.Any(), character.ID).
					Return(model.ChatChannels{chatChannel}, nil)
				mockCharService.EXPECT().
					GetCharacter(gomock.Any(), req).
					Return(character.ToPb(), nil)
				out, err := server.GetAuthorizedChatChannels(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out.Channels).To(HaveLen(1))
			})

			It("should work (admin other)", func() {
				mockChatService.EXPECT().
					AuthorizedChannelsForCharacter(gomock.Any(), character.ID).
					Return(model.ChatChannels{chatChannel}, nil)
				mockCharService.EXPECT().
					GetCharacter(gomock.Any(), req).
					Return(character.ToPb(), nil)
				out, err := server.GetAuthorizedChatChannels(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out.Channels).To(HaveLen(1))
			})

			It("should work (player self)", func() {
				mockChatService.EXPECT().
					AuthorizedChannelsForCharacter(gomock.Any(), character.ID).
					Return(model.ChatChannels{chatChannel}, nil)
				mockCharService.EXPECT().
					GetCharacter(gomock.Any(), req).
					Return(character.ToPb(), nil)
				out, err := server.GetAuthorizedChatChannels(incPlayerCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out.Channels).To(HaveLen(1))
			})
		})

		When("given invalid input", func() {
			It("should error if invalid context (nil)", func() {
				out, err := server.GetAuthorizedChatChannels(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid context (no claims)", func() {
				out, err := server.GetAuthorizedChatChannels(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (guest)", func() {
				out, err := server.GetAuthorizedChatChannels(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (player other)", func() {
				character.OwnerId = *admin.ID
				mockChatService.EXPECT().
					AuthorizedChannelsForCharacter(gomock.Any(), character.ID).
					Return(nil, fakeErr)
				mockCharService.EXPECT().
					GetCharacter(gomock.Any(), gomock.Any()).
					Return(character.ToPb(), nil)
				out, err := server.GetAuthorizedChatChannels(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if getting character failed", func() {
				character.OwnerId = *admin.ID
				mockCharService.EXPECT().
					GetCharacter(gomock.Any(), req).
					Return(nil, fakeErr)
				out, err := server.GetAuthorizedChatChannels(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if no character is found", func() {
				character.OwnerId = *admin.ID
				mockCharService.EXPECT().
					GetCharacter(gomock.Any(), req).
					Return(nil, nil)
				out, err := server.GetAuthorizedChatChannels(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("UpdateUserchatChannelAuthorizations", func() {
		var (
			req *pb.RequestChatChannelAuthChange
		)
		BeforeEach(func() {
			req = &pb.RequestChatChannelAuthChange{
				Character: &pb.CharacterTarget{
					Type: &pb.CharacterTarget_Id{
						Id: uint64(character.ID),
					},
				},
				Add: true,
				Ids: []uint64{1, 2, 3},
			}
		})
		When("given valid input", func() {
			It("should work (admin)", func() {
				mockChatService.EXPECT().
					ChangeAuthorizationForCharacter(gomock.Any(), character.ID, *helpers.ArrayOfUint64ToUint(&req.Ids), req.Add).
					Return(nil)
				out, err := server.UpdateUserChatChannelAuthorizations(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error if invalid context (nil)", func() {
				out, err := server.UpdateUserChatChannelAuthorizations(nil, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid context (no claims)", func() {
				out, err := server.UpdateUserChatChannelAuthorizations(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (player)", func() {
				out, err := server.UpdateUserChatChannelAuthorizations(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if invalid permission (guest)", func() {
				out, err := server.UpdateUserChatChannelAuthorizations(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if error during getting character target", func() {
				req.Character.Type = nil
				out, err := server.UpdateUserChatChannelAuthorizations(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if error during update", func() {
				mockChatService.EXPECT().
					ChangeAuthorizationForCharacter(gomock.Any(), character.ID, *helpers.ArrayOfUint64ToUint(&req.Ids), req.Add).
					Return(fakeErr)
				out, err := server.UpdateUserChatChannelAuthorizations(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should unique error if conflict exists", func() {
				mockChatService.EXPECT().
					ChangeAuthorizationForCharacter(gomock.Any(), character.ID, *helpers.ArrayOfUint64ToUint(&req.Ids), req.Add).
					Return(gorm.ErrDuplicatedKey)
				out, err := server.UpdateUserChatChannelAuthorizations(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})
})

package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("Inner-ChatService")

	// ErrEmptyMessage thrown when a request to send a message is made with an empty message
	ErrEmptyMessage = errors.New("message cannot be empty")
)

type ChatService interface {
	AllChannels(ctx context.Context) (model.ChatChannels, error)
	GetChannel(ctx context.Context, id uint) (*model.ChatChannel, error)
	UpdateChannel(ctx context.Context, pb *pb.UpdateChatChannelRequest) error
	CreateChannel(ctx context.Context, channel *model.ChatChannel) (*model.ChatChannel, error)
	DeleteChannel(ctx context.Context, channel *model.ChatChannel) error

	SendChannelMessage(ctx context.Context, username string, message string, channelId uint) error
	SendDirectMessage(ctx context.Context, senderCharacter string, message string, targetCharacter string) error

	ChannelMessagesReader(ctx context.Context, channelId uint) *kafka.Reader
	DirectMessagesReader(ctx context.Context, username string) *kafka.Reader

	AuthorizedChannelsForCharacter(ctx context.Context, character string) (model.ChatChannels, error)
	ChangeAuthorizationForCharacter(ctx context.Context, character string, channelIds []uint, addAuth bool) error
}

type chatService struct {
	chatRepo  repository.ChatRepository
	kafkaConn *kafka.Conn

	channelMessageWriters map[uint]*kafka.Writer
	directMessageWriters  map[string]*kafka.Writer
}

func (s chatService) ChangeAuthorizationForCharacter(ctx context.Context, character string, channelIds []uint, addAuth bool) error {
	return s.chatRepo.ChangeAuthorizationForCharacter(ctx, character, channelIds, addAuth)
}

func (s chatService) AuthorizedChannelsForCharacter(ctx context.Context, character string) (model.ChatChannels, error) {
	return s.chatRepo.AuthorizedChannelsForCharacter(ctx, character)
}

func (s chatService) UpdateChannel(ctx context.Context, pb *pb.UpdateChatChannelRequest) error {
	ctx, span := tracer.Start(ctx, "UpdateChannel")
	defer span.End()

	channel, err := s.GetChannel(ctx, uint(pb.ChannelId))
	if err != nil {
		return err
	}

	if pb.Name != nil {
		channel.Name = pb.Name.Value
	}

	if pb.Public != nil {
		channel.Public = pb.Public.Value
	}

	return s.chatRepo.UpdateChannel(ctx, channel)
}

func (s chatService) GetChannel(ctx context.Context, id uint) (*model.ChatChannel, error) {
	return s.chatRepo.GetChannel(ctx, id)
}

func (s chatService) ChannelMessagesReader(ctx context.Context, channelId uint) *kafka.Reader {
	ctx, span := tracer.Start(ctx, "ChannelMessagesReader")
	defer span.End()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{s.kafkaConn.RemoteAddr().String()},
		Topic:    topicNameFromChannel(channelId),
		MinBytes: 1,
		MaxBytes: 10e3,
	})
	_ = r.SetOffset(kafka.LastOffset)

	return r
}

func (s chatService) DirectMessagesReader(ctx context.Context, characterName string) *kafka.Reader {
	ctx, span := tracer.Start(ctx, "DirectMessagesReader")
	defer span.End()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{s.kafkaConn.RemoteAddr().String()},
		Topic:    topicNameFromUser(characterName),
		MinBytes: 1,
		MaxBytes: 10e3,
	})
	_ = r.SetOffset(kafka.LastOffset)

	return r
}
func (s chatService) SendChannelMessage(ctx context.Context, characterName string, message string, channelId uint) error {
	ctx, span := tracer.Start(ctx, "SendChannelMessage")
	defer span.End()

	if len(message) == 0 {
		return ErrEmptyMessage
	}

	w := s.getChannelMessageWriter(ctx, channelId)

	return w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(characterName),
			Value: []byte(message),
		},
	)
}

func (s chatService) SendDirectMessage(ctx context.Context, characterName string, message string, targetCharacterName string) error {
	ctx, span := tracer.Start(ctx, "SendDirectMessage")
	defer span.End()

	if len(message) == 0 {
		return ErrEmptyMessage
	}
	w := s.getUserMessageWriter(ctx, targetCharacterName)

	return w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(characterName),
			Value: []byte(message),
		},
	)
}

func (s chatService) AllChannels(ctx context.Context) (model.ChatChannels, error) {
	return s.chatRepo.AllChannels(ctx)
}

func (s chatService) CreateChannel(ctx context.Context, channel *model.ChatChannel) (*model.ChatChannel, error) {
	newChannel, err := s.chatRepo.CreateChannel(ctx, channel)
	if err != nil {
		return nil, err
	}

	_ = s.kafkaConn.CreateTopics(createTopicConfigFromChannel(newChannel))
	return newChannel, nil
}

func (s chatService) DeleteChannel(ctx context.Context, channel *model.ChatChannel) error {
	err := s.chatRepo.DeleteChannel(ctx, channel)
	if err != nil {
		return err
	}

	_ = s.kafkaConn.DeleteTopics(topicNameFromChannel(channel.ID))

	return nil
}

func NewChatService(ctx context.Context, chatRepo repository.ChatRepository, kafkaAddress config.ServerAddress) (ChatService, error) {
	ctx, span := tracer.Start(ctx, "NewChatService")
	defer span.End()

	conn, err := repository.ConnectKafka(kafkaAddress)
	if err != nil {
		return nil, fmt.Errorf("connecting kafka: %v", err)
	}

	service := chatService{
		chatRepo:              chatRepo,
		kafkaConn:             conn,
		channelMessageWriters: map[uint]*kafka.Writer{},
		directMessageWriters:  map[string]*kafka.Writer{},
	}

	channels, err := service.AllChannels(ctx)
	if err != nil {
		return nil, err
	}

	topicConfigs := make([]kafka.TopicConfig, len(channels))

	for idx, channel := range channels {
		topicConfigs[idx] = createTopicConfigFromChannel(channel)
	}

	err = service.kafkaConn.CreateTopics(topicConfigs...)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func createTopicConfigFromChannel(channel *model.ChatChannel) kafka.TopicConfig {
	return kafka.TopicConfig{
		Topic:             topicNameFromChannel(channel.ID),
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
}

func topicNameFromChannel(channelId uint) string {
	return fmt.Sprintf("chat-channel-%d", channelId)
}

func topicNameFromUser(username string) string {
	return fmt.Sprintf("chat-user-%s", username)
}

func (s chatService) getChannelMessageWriter(ctx context.Context, channelId uint) *kafka.Writer {
	ctx, span := tracer.Start(ctx, "GetChannelMessageWriter")
	defer span.End()

	if s.channelMessageWriters[channelId] == nil {
		s.channelMessageWriters[channelId] = &kafka.Writer{
			Addr:     s.kafkaConn.RemoteAddr(),
			Topic:    topicNameFromChannel(channelId),
			Balancer: &kafka.LeastBytes{},
			Async:    true,
		}
	}

	return s.channelMessageWriters[channelId]
}

func (s chatService) getUserMessageWriter(ctx context.Context, username string) *kafka.Writer {
	ctx, span := tracer.Start(ctx, "GetChannelMessageWriter")
	defer span.End()

	if s.directMessageWriters[username] == nil {
		s.directMessageWriters[username] = &kafka.Writer{
			Addr:     s.kafkaConn.RemoteAddr(),
			Topic:    topicNameFromUser(username),
			Balancer: &kafka.LeastBytes{},
		}
	}

	return s.directMessageWriters[username]
}

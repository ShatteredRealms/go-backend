package chat

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	ServiceName = "chat"
)

type ChatServerContext struct {
	GlobalConfig     *config.GlobalConfig
	ChatService      service.ChatService
	CharacterService pb.CharacterServiceClient
	KeycloakClient   *gocloak.GoCloak
	Tracer           trace.Tracer
}

func NewServerContext(ctx context.Context, conf *config.GlobalConfig) (*ChatServerContext, error) {
	server := &ChatServerContext{
		GlobalConfig:   conf,
		Tracer:         otel.Tracer("ChatService"),
		KeycloakClient: gocloak.NewClient(conf.Keycloak.BaseURL),
	}

	db, err := repository.ConnectDB(conf.Chat.Postgres)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres database: %w", err)
	}

	repo := repository.NewChatRepository(db)
	chatService, err := service.NewChatService(ctx, repo, conf.Chat.Kafka)
	if err != nil {
		return nil, fmt.Errorf("creating chat service: %w", err)
	}
	server.ChatService = chatService

	charactersConn, err := helpers.GrpcClientWithOtel(conf.Character.Remote.Address())
	if err != nil {
		return nil, fmt.Errorf("connecting characters service: %w", err)
	}
	server.CharacterService = pb.NewCharacterServiceClient(charactersConn)

	return server, nil
}

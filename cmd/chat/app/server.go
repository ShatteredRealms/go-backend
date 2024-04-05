package chat

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	"github.com/WilSimpson/gocloak/v13"
	"go.opentelemetry.io/otel/trace"
)

var (
	ServiceName = "chat"
)

type ChatServerContext struct {
	*config.ServerContext
	ChatService      service.ChatService
	CharacterService pb.CharacterServiceClient
}

func NewServerContext(ctx context.Context, conf *config.GlobalConfig, tracer trace.Tracer) (*ChatServerContext, error) {
	server := &ChatServerContext{
		ServerContext: &config.ServerContext{
			GlobalConfig:   conf,
			Tracer:         tracer,
			KeycloakClient: gocloak.NewClient(conf.Keycloak.BaseURL),
			RefSROServer:   &conf.Chat.SROServer,
		},
		ChatService:      nil,
		CharacterService: nil,
	}
	ctx, span := server.Tracer.Start(ctx, "server.connect")
	defer span.End()

	server.KeycloakClient.RegisterMiddlewares(gocloak.OpenTelemetryMiddleware)

	db, err := repository.ConnectDB(conf.Chat.Postgres, conf.Redis)
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

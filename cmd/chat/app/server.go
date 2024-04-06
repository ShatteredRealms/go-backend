package chat

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/config"
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
	ChatService service.ChatService
}

func NewServerContext(ctx context.Context, conf *config.GlobalConfig, tracer trace.Tracer) (*ChatServerContext, error) {
	ctx, span := tracer.Start(ctx, "server.context.new")
	defer span.End()

	server := &ChatServerContext{
		ServerContext: config.NewServerContext(ctx, conf, tracer, &conf.Chat.SROServer),
	}

	server.KeycloakClient.RegisterMiddlewares(gocloak.OpenTelemetryMiddleware)

	db, err := repository.ConnectDB(ctx, conf.Chat.Postgres, conf.Redis)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres database: %w", err)
	}

	repo := repository.NewChatRepository(db)
	chatService, err := service.NewChatService(ctx, repo, conf.Chat.Kafka)
	if err != nil {
		return nil, fmt.Errorf("creating chat service: %w", err)
	}
	server.ChatService = chatService

	return server, nil
}

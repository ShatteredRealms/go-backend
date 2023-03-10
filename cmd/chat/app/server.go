package chat

import (
	"context"
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
	CharacterService pb.CharactersServiceClient
	KeycloakClient   *gocloak.GoCloak
	Tracer           trace.Tracer
}

func NewServerContext(ctx context.Context, conf *config.GlobalConfig) *ChatServerContext {
	server := &ChatServerContext{
		GlobalConfig:   conf,
		Tracer:         otel.Tracer("ChatService"),
		KeycloakClient: gocloak.NewClient(conf.GameBackend.Keycloak.BaseURL),
	}

	db, err := repository.ConnectDB(conf.Chat.DB)
	helpers.Check(ctx, err, "connecting to database")

	repo := repository.NewChatRepository(db)
	chatService, err := service.NewChatService(ctx, repo, conf.Chat.Kafka)
	helpers.Check(ctx, err, "chat service")
	server.ChatService = chatService

	charactersConn, err := helpers.GrpcClientWithOtel(conf.Characters.Remote.Address())
	helpers.Check(ctx, err, "connect characters service")
	server.CharacterService = pb.NewCharactersServiceClient(charactersConn)

	return server
}

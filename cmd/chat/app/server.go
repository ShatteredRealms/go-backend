package chat

import (
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type ChatServer struct {
	GlobalConfig     *config.GlobalSROConfig
	ChatService      service.ChatService
	CharacterService pb.CharactersServiceClient
	Tracer           trace.Tracer
}

func NewServer(ctx context.Context, conf *config.GlobalSROConfig) *ChatServer {
	server := &ChatServer{
		GlobalConfig: conf,
		Tracer:       otel.Tracer("ChatService"),
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

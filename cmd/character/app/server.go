package character

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.opentelemetry.io/otel/trace"
)

var (
	ServiceName = "character"
)

type CharacterServerContext struct {
	*config.ServerContext
	CharacterService service.CharacterService
	InventoryService service.InventoryService
}

func NewServerContext(ctx context.Context, conf *config.GlobalConfig, tracer trace.Tracer) (*CharacterServerContext, error) {
	ctx, span := tracer.Start(ctx, "server.context.new")
	defer span.End()

	server := &CharacterServerContext{
		ServerContext: config.NewServerContext(ctx, conf, tracer, &conf.Character.SROServer),
	}

	postgres, err := repository.ConnectDB(ctx, conf.Character.Postgres, conf.Redis)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	characterRepo, err := repository.NewCharacterRepository(postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	characterService, err := service.NewCharacterService(ctx, characterRepo)
	if err != nil {
		return nil, fmt.Errorf("character service: %w", err)
	}
	server.CharacterService = characterService

	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI(server.GlobalConfig.Character.Mongo.Master.MongoDSN())
	mongoDb, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("connecting to mongo database: %w", err)
	}
	invRepo := repository.NewInventoryRepository(mongoDb.Database(server.GlobalConfig.Character.Mongo.Master.Name))
	server.InventoryService = service.NewInventoryService(invRepo)

	return server, nil
}

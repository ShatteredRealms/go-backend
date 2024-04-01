package character

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	ServiceName = "character"
)

type CharactersServerContext struct {
	GlobalConfig     *config.GlobalConfig
	CharacterService service.CharacterService
	InventoryService service.InventoryService
	KeycloakClient   *gocloak.GoCloak
	Tracer           trace.Tracer
}

func NewServerContext(ctx context.Context, conf *config.GlobalConfig) *CharactersServerContext {
	server := &CharactersServerContext{
		GlobalConfig:   conf,
		Tracer:         otel.Tracer("CharactersService"),
		KeycloakClient: gocloak.NewClient(conf.Keycloak.BaseURL),
	}

	postgres, err := repository.ConnectDB(server.GlobalConfig.Character.Postgres)
	helpers.Check(ctx, err, "connecting to postgres database")

	characterRepo, err := repository.NewCharacterRepository(postgres)
	helpers.Check(ctx, err, "character repo")
	characterService, err := service.NewCharacterService(ctx, characterRepo)
	helpers.Check(ctx, err, "character serivce")
	server.CharacterService = characterService

	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI(server.GlobalConfig.Character.Mongo.Master.MongoDSN())
	mongoDb, err := mongo.Connect(ctx, opts)
	helpers.Check(ctx, err, "connecting to mongo database")
	invRepo := repository.NewInventoryRepository(mongoDb.Database(server.GlobalConfig.Character.Mongo.Master.Name))
	server.InventoryService = service.NewInventoryService(invRepo)

	return server
}

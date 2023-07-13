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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	ServiceName = "characters"
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
		KeycloakClient: gocloak.NewClient(conf.GameBackend.Keycloak.BaseURL),
	}

	postgres, err := repository.ConnectDB(server.GlobalConfig.Character.PostgresConfig)
	helpers.Check(ctx, err, "connecting to postgres database")

	characterRepo := repository.NewCharacterRepository(postgres)
	characterService, err := service.NewCharacterService(ctx, characterRepo)
	helpers.Check(ctx, err, "character serivce")
	server.CharacterService = characterService

	mongoDb, err := mongo.Connect(ctx, options.Client().ApplyURI(server.GlobalConfig.Character.MongoConfig.Master.MongoDSN()))
	helpers.Check(ctx, err, "connecting to mongo database")
	invRepo := repository.NewInventoryRepository(mongoDb.Database(server.GlobalConfig.Character.MongoConfig.Master.Name))
	server.InventoryService = service.NewInventoryService(invRepo)

	return server
}

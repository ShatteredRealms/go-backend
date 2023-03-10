package characters

import (
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	ServiceName = "characters"
)

type CharactersServerContext struct {
	GlobalConfig *config.GlobalConfig
	Service      service.CharacterService
	Tracer       trace.Tracer
}

func NewServerContext(ctx context.Context, conf *config.GlobalConfig) *CharactersServerContext {
	server := &CharactersServerContext{
		GlobalConfig: conf,
		Tracer:       otel.Tracer("CharactersService"),
	}

	db, err := repository.ConnectDB(server.GlobalConfig.Characters.DB)
	helpers.Check(ctx, err, "connecting to database")

	repo := repository.NewCharacterRepository(db)
	server.Service = service.NewCharacterService(repo)
	return server
}

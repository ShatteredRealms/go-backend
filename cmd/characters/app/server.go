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

type CharactersServer struct {
	GlobalConfig *config.GlobalSROConfig
	Service      service.CharacterService
	Tracer       trace.Tracer
}

func NewServer(ctx context.Context, conf *config.GlobalSROConfig) *CharactersServer {
	server := &CharactersServer{
		GlobalConfig: conf,
		Tracer:       otel.Tracer("CharactersServer"),
	}

	db, err := repository.ConnectDB(server.GlobalConfig.Characters.DB)
	helpers.Check(ctx, err, "connecting to database")

	repo := repository.NewCharacterRepository(db)
	server.Service = service.NewCharacterService(repo)
	return server
}

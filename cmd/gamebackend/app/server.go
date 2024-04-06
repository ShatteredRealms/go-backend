package gamebackend

import (
	"context"
	"fmt"

	"agones.dev/agones/pkg/client/clientset/versioned"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/client-go/rest"
)

var (
	ServiceName = "gamebackend"
)

type GameBackendServerContext struct {
	*config.ServerContext
	GamebackendService service.GamebackendService
	AgonesClient       versioned.Interface
}

func NewServerContext(ctx context.Context, conf *config.GlobalConfig, tracer trace.Tracer) (*GameBackendServerContext, error) {
	ctx, span := tracer.Start(ctx, "server.context.new")
	defer span.End()

	server := &GameBackendServerContext{
		ServerContext: config.NewServerContext(ctx, conf, tracer, &conf.GameBackend.SROServer),
	}

	db, err := repository.ConnectDB(ctx, conf.GameBackend.Postgres, conf.Redis)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres database: %w", err)
	}

	repo := repository.NewGamebackendRepository(db)
	gamebackendService, err := service.NewGamebackendService(ctx, repo)
	if err != nil {
		return nil, fmt.Errorf("creating gamebackend service: %w", err)
	}
	server.GamebackendService = gamebackendService

	if conf.GameBackend.Mode != config.LocalMode {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("creating cluster config: %w", err)
		}

		server.AgonesClient, err = versioned.NewForConfig(config)
		if err != nil {
			return nil, fmt.Errorf("creating agones client from config: %w", err)
		}
	}

	return server, nil
}

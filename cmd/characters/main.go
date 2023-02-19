package main

import (
	"context"
	"net"
	"net/http"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/uptrace-go/uptrace"
)

type appConfig struct {
	Characters config.Server        `yaml:"characters"`
	Accounts   config.Server        `yaml:"accounts"`
	KeyDir     string               `yaml:"keyDir"`
	Uptrace    config.UptraceConfig `yaml:"uptrace"`
}

var (
	conf = &appConfig{
		Characters: config.Server{
			Local: config.ServerAddress{
				Port: 8081,
				Host: "",
			},
			Remote: config.ServerAddress{
				Port: 8081,
				Host: "",
			},
			Mode:     "development",
			LogLevel: log.InfoLevel,
			DB: config.DBPoolConfig{
				Master: config.DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "characters",
					Username: "postgres",
					Password: "password",
				},
				Slaves: []config.DBConfig{},
			},
		},
		Accounts: config.Server{
			Remote: config.ServerAddress{
				Port: 8080,
				Host: "",
			},
		},
		KeyDir: "/etc/sro/auth",
	}
)

func init() {
	helpers.SetupLogs()
	config.SetupConfig(conf)
}

func main() {
	ctx := context.Background()
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(conf.Uptrace.DSN()),
		uptrace.WithServiceName("characters_service"),
		uptrace.WithServiceVersion("v1.0.0"),
	)
	defer uptrace.Shutdown(ctx)

	db, err := repository.ConnectDB(conf.Characters.DB)
	helpers.Check(ctx, err, "db connect from file")

	characterRepo := repository.NewCharacterRepository(db)
	helpers.Check(ctx, characterRepo.Migrate(), "character repo")

	characterService := service.NewCharacterService(characterRepo)
	jwtService, err := service.NewJWTService(conf.KeyDir)
	helpers.Check(ctx, err, "jwt service")

	grpcServer, gwmux, err := NewServer(characterService, jwtService)
	helpers.Check(ctx, err, "create grpc server")

	lis, err := net.Listen("tcp", conf.Characters.Local.Address())
	helpers.Check(ctx, err, "listen")

	server := &http.Server{
		Addr:    conf.Characters.Local.Address(),
		Handler: helpers.GRPCHandlerFunc(grpcServer, gwmux),
	}

	log.Info("Server starting")

	err = server.Serve(lis)
	helpers.Check(ctx, err, "serve")
}

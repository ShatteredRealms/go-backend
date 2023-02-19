package main

import (
	"context"
	"net"
	"net/http"

	chat "github.com/ShatteredRealms/go-backend/cmd/chat/global"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
)

type appConfig struct {
	Chat       config.Server        `yaml:"chat"`
	Accounts   config.Server        `yaml:"accounts"`
	Characters config.Server        `yaml:"characters"`
	KeyDir     string               `yaml:"keyDir"`
	Uptrace    config.UptraceConfig `yaml:"uptrace"`
}

var (
	conf = &appConfig{
		Chat: config.Server{
			Local: config.ServerAddress{
				Port: 8180,
				Host: "",
			},
			Remote: config.ServerAddress{
				Port: 8180,
				Host: "",
			},
			Kafka: config.ServerAddress{
				Port: 29092,
				Host: "localhost",
			},
			Mode:     "development",
			LogLevel: log.InfoLevel,
			DB: config.DBPoolConfig{
				Master: config.DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "chat",
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
		Characters: config.Server{
			Remote: config.ServerAddress{
				Port: 8081,
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
		uptrace.WithServiceName("chat_service"),
		uptrace.WithServiceVersion("v1.0.0"),
	)
	defer uptrace.Shutdown(ctx)

	chat.Tracer = otel.Tracer("chat")
	ctx, span := chat.Tracer.Start(ctx, "main")
	db, err := repository.ConnectDB(conf.Chat.DB)
	helpers.Check(ctx, err, "db connect from file")

	chatRepository := repository.NewChatRepository(db)
	helpers.Check(ctx, chatRepository.Migrate(ctx), "role repo")

	chatService, err := service.NewChatService(ctx, chatRepository, conf.Chat.Kafka)
	helpers.Check(ctx, err, "chat service")

	jwtService, err := service.NewJWTService(conf.KeyDir)
	helpers.Check(ctx, err, "jwt service")

	grpcServer, gwmux, err := NewServer(ctx, jwtService, chatService)
	helpers.Check(ctx, err, "create grpc server")

	lis, err := net.Listen("tcp", conf.Chat.Local.Address())
	helpers.Check(ctx, err, "listen")

	server := &http.Server{
		Addr:    conf.Chat.Local.Address(),
		Handler: helpers.GRPCHandlerFunc(grpcServer, gwmux),
	}

	log.Info("Server starting")

	span.End()

	err = server.Serve(lis)
	helpers.Check(ctx, err, "serve")
}

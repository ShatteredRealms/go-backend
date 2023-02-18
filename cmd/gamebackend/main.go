package main

import (
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/kend/pkg/helpers"
	"github.com/kend/pkg/service"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/uptrace-go/uptrace"
	"net"
	"net/http"
)

type appConfig struct {
	GameBackend config.Server        `yaml:"gameBackend"`
	Accounts    config.Server        `yaml:"accounts"`
	Characters  config.Server        `yaml:"characters"`
	KeyDir      string               `yaml:"keyDir"`
	Agones      agonesConfig         `yaml:"agones"`
	Uptrace     config.UptraceConfig `yaml:"uptrace"`
}

type agonesConfig struct {
	KeyFile    string        `yaml:"keyFile"`
	CertFile   string        `yaml:"certFile"`
	CaCertFile string        `yaml:"caCertFile"`
	Namespace  string        `yaml:"namespace"`
	Allocator  config.Server `yaml:"allocator"`
}

var (
	conf = &appConfig{
		GameBackend: config.Server{
			Local: config.ServerAddress{
				Port: 8082,
				Host: "",
			},
			Remote: config.ServerAddress{
				Port: 8082,
				Host: "",
			},
			Mode:     config.ModeDevelopment,
			LogLevel: log.InfoLevel,
			DB: config.DBPoolConfig{
				Master: config.DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "gamebackend",
					Username: "postgres",
					Password: "password",
				},
				Slaves: []config.DBConfig{},
			},
		},
		Characters: config.Server{
			Remote: config.ServerAddress{
				Port: 8081,
				Host: "",
			},
		},
		Accounts: config.Server{
			Remote: config.ServerAddress{
				Port: 8080,
				Host: "",
			},
		},
		KeyDir: "/etc/sro/auth",
		Agones: agonesConfig{
			KeyFile:    "/etc/sro/auth/agones/client/key",
			CertFile:   "/etc/sro/auth/agones/client/cert",
			CaCertFile: "/etc/sro/auth/agones/ca/ca",
			Namespace:  "default",
			Allocator: config.Server{
				Remote: config.ServerAddress{
					Port: 443,
					Host: "",
				},
			},
		},
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
		uptrace.WithServiceName("gamebackend_service"),
		uptrace.WithServiceVersion("v1.0.0"),
	)
	defer uptrace.Shutdown(ctx)

	jwtService, err := service.NewJWTService(conf.KeyDir)
	helpers.Check(ctx, err, "jwt service")

	grpcServer, gwmux, err := NewServer(jwtService)
	helpers.Check(ctx, err, "create grpc server")

	lis, err := net.Listen("tcp", conf.GameBackend.Local.Address())
	helpers.Check(ctx, err, "listen")

	server := &http.Server{
		Addr:    conf.GameBackend.Local.Address(),
		Handler: helpers.GRPCHandlerFunc(grpcServer, gwmux),
	}

	log.Info("Server starting")

	err = server.Serve(lis)
	helpers.Check(ctx, err, "serve")
}

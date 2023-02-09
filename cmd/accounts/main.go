package main

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/config"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/repository"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"net"
	"net/http"
)

type appConfig struct {
	Accounts config.Server        `yaml:"accounts"`
	KeyDir   string               `yaml:"keyDir"`
	Uptrace  config.UptraceConfig `yaml:"uptrace"`
}

var (
	version     = "v0.0.1"
	serviceName = "accounts_service"
	conf        = &appConfig{
		Accounts: config.Server{
			Local: config.ServerAddress{
				Port: 8080,
				Host: "",
			},
			Remote: config.ServerAddress{
				Port: 8080,
				Host: "",
			},
			Mode:     "development",
			LogLevel: log.InfoLevel,
			DB: config.DBPoolConfig{
				Master: config.DBConfig{
					Host:     "localhost",
					Port:     "5432",
					Name:     "accounts",
					Username: "postgres",
					Password: "password",
				},
				Slaves: []config.DBConfig{},
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
		uptrace.WithServiceName(serviceName),
		uptrace.WithServiceVersion(version),
	)
	defer uptrace.Shutdown(ctx)

	ctx, span := otel.Tracer("accounts").Start(ctx, "main")
	db, err := repository.ConnectDB(conf.Accounts.DB)
	helpers.Check(ctx, err, "db connect from file")

	permissionRepository := repository.NewPermissionRepository(db)
	helpers.Check(ctx, permissionRepository.Migrate(), "permissions repo")
	roleRepository := repository.NewRoleRepository(db)
	helpers.Check(ctx, roleRepository.Migrate(), "role repo")
	userRepository := repository.NewUserRepository(db)
	helpers.Check(ctx, userRepository.Migrate(), "user repo")

	permissionService := service.NewPermissionService(permissionRepository)
	roleService := service.NewRoleService(roleRepository)
	userService := service.NewUserService(userRepository)
	jwtService, err := service.NewJWTService(conf.KeyDir)
	helpers.Check(ctx, err, "jwt service")

	grpcServer, gwmux, err := NewServer(userService, permissionService, roleService, jwtService, ctx)
	helpers.Check(ctx, err, "create grpc server")

	seedDatabaseIfNeeded(ctx, userService, permissionService, roleService, grpcServer.GetServiceInfo())

	lis, err := net.Listen("tcp", conf.Accounts.Local.Address())
	helpers.Check(ctx, err, "listen")

	server := &http.Server{
		Addr:    conf.Accounts.Local.Address(),
		Handler: helpers.GRPCHandlerFunc(grpcServer, gwmux),
	}

	log.Info("Server starting")
	span.End()
	err = server.Serve(lis)
	helpers.Check(ctx, err, "serve")
}

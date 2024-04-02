package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	"github.com/ShatteredRealms/go-backend/pkg/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	character "github.com/ShatteredRealms/go-backend/cmd/character/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
)

var (
	conf *config.GlobalConfig
)

func init() {
	var err error
	conf, err = config.NewGlobalConfig(context.Background())
	if err != nil {
		log.Logger.Fatalf("initialization: %v", err)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := telemetry.SetupOTelSDK(ctx, character.ServiceName, conf.Version, conf.OpenTelemetry.Addr)
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
		if err != nil {
			log.Logger.Infof("Error shutting down: %v", err)
		}
	}()

	if err != nil {
		log.Logger.Errorf("connecting to otel: %w", err)
		return
	}

	server, err := character.NewServerContext(ctx, conf)
	if err != nil {
		log.Logger.Errorf("creating server context: %w", err)
		return
	}
	grpcServer, gwmux := helpers.InitServerDefaults(server.KeycloakClient, server.GlobalConfig.Keycloak.Realm)
	address := server.GlobalConfig.Character.Local.Address()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	pb.RegisterHealthServiceServer(grpcServer, srv.NewHealthServiceServer())
	err = pb.RegisterHealthServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	if err != nil {
		log.Logger.Errorf("register health service handler endpoint: %w", err)
		return
	}

	css, err := srv.NewCharacterServiceServer(ctx, server)
	if err != nil {
		log.Logger.Errorf("create character service server: %w", err)
		return
	}
	pb.RegisterCharacterServiceServer(grpcServer, css)
	err = pb.RegisterCharacterServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	if err != nil {
		log.Logger.Errorf("registering character service handler endpoint: %w", err)
		return
	}

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- helpers.StartServer(ctx, grpcServer, gwmux, server.GlobalConfig.Character.Local.Address())
	}()

	select {
	case err = <-srvErr:
		log.Logger.Errorf("listen server: %v", err)

	case <-ctx.Done():
		log.Logger.Info("Server canceled by user input.")
		stop()
	}

}

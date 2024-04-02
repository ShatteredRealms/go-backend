package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	chat "github.com/ShatteredRealms/go-backend/cmd/chat/app"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	"github.com/ShatteredRealms/go-backend/pkg/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ShatteredRealms/go-backend/pkg/config"
)

var (
	conf *config.GlobalConfig
)

func init() {
	conf = config.NewGlobalConfig(context.Background())
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := telemetry.SetupOTelSDK(ctx, chat.ServiceName, config.Version, conf.OpenTelemetry.Addr)
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
		if err != nil {
			log.Logger.Infof("Error shutting down: %v", err)
		}
	}()

	if err != nil {
		log.Logger.Fatal(err)
	}

	server := chat.NewServerContext(ctx, conf)
	grpcServer, gwmux := helpers.InitServerDefaults(server.KeycloakClient, server.GlobalConfig.Keycloak.Realm)
	address := server.GlobalConfig.Chat.Local.Address()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	pb.RegisterHealthServiceServer(grpcServer, srv.NewHealthServiceServer())
	err = pb.RegisterHealthServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "register health service handler endpoint")

	srvService, err := srv.NewChatServiceServer(ctx, server)
	helpers.Check(ctx, err, "create chat service")
	pb.RegisterChatServiceServer(grpcServer, srvService)
	err = pb.RegisterChatServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "register chat service handler endpoint")

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- helpers.StartServer(ctx, grpcServer, gwmux, server.GlobalConfig.Chat.Local.Address())
	}()

	select {
	case err = <-srvErr:
		log.Logger.Fatalf("listen server: %v", err)

	case <-ctx.Done():
		log.Logger.Info("Server canceled by user input.")
		stop()
	}
}

package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	character "github.com/ShatteredRealms/go-backend/cmd/character/app"
	gamebackend "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ShatteredRealms/go-backend/pkg/config"
)

var (
	conf *config.GlobalConfig
)

func init() {
	helpers.SetupLogger(gamebackend.ServiceName)
	conf = config.NewGlobalConfig(context.Background())
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := helpers.SetupOTelSDK(ctx, character.ServiceName, config.Version, conf.OpenTelemetry.Addr)
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
		if err != nil {
			log.Logger.Infof("Error shutting down: %v", err)
		}
	}()

	if err != nil {
		log.Logger.Fatal(err)
	}
	server := gamebackend.NewServerContext(ctx, conf)
	grpcServer, gwmux := helpers.InitServerDefaults()
	address := server.GlobalConfig.GameBackend.Local.Address()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	pb.RegisterHealthServiceServer(grpcServer, srv.NewHealthServiceServer())
	err = pb.RegisterHealthServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "register health service handler endpoint")

	connServer, err := srv.NewConnectionServiceServer(ctx, server)
	helpers.Check(ctx, err, "creating connection service server")
	pb.RegisterConnectionServiceServer(grpcServer, connServer)
	err = pb.RegisterConnectionServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "register connection service handler endpoint")

	serverManagerServer, err := srv.NewServerManagerServiceServer(ctx, server)
	helpers.Check(ctx, err, "creating server manager service server")
	pb.RegisterServerManagerServiceServer(grpcServer, serverManagerServer)
	err = pb.RegisterServerManagerServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "register server manager service handler endpoint")

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- helpers.StartServer(ctx, grpcServer, gwmux, server.GlobalConfig.GameBackend.Local.Address())
	}()

	select {
	case err = <-srvErr:
		log.Logger.Fatalf("listen server: %v", err)

	case <-ctx.Done():
		log.Logger.Info("Server canceled by user input.")
		stop()
	}
}

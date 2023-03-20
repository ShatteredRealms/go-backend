package main

import (
	"context"
	gamebackend "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/uptrace/uptrace-go/uptrace"
)

var (
	conf *config.GlobalConfig
)

func init() {
	helpers.SetupLogger()
	conf = config.NewGlobalConfig()
}

func main() {
	ctx := context.Background()
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(conf.Uptrace.DSN),
		uptrace.WithServiceName(gamebackend.ServiceName),
		uptrace.WithServiceVersion(conf.Version),
	)

	server := gamebackend.NewServerContext(ctx, conf)
	grpcServer, gwmux := helpers.InitServerDefaults()
	address := server.GlobalConfig.GameBackend.Local.Address()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	pb.RegisterHealthServiceServer(grpcServer, srv.NewHealthServiceServer())
	err := pb.RegisterHealthServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "register health service handler endpoint")

	pb.RegisterConnectionServiceServer(grpcServer, srv.NewConnectionServiceServer(server))
	err = pb.RegisterConnectionServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "register connection service handler endpoint")

	helpers.StartServer(ctx, grpcServer, gwmux, server.GlobalConfig.GameBackend.Local.Address())
}

package main

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	"github.com/uptrace/uptrace-go/uptrace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ShatteredRealms/go-backend/cmd/character/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
)

var (
	conf *config.GlobalConfig
)

func init() {
	helpers.SetupLogger()
	conf = config.NewGlobalConfig(context.Background())
}

func main() {
	ctx := context.Background()
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(conf.Uptrace.DSN),
		uptrace.WithServiceName(character.ServiceName),
		uptrace.WithServiceVersion(conf.Version),
	)

	server := character.NewServerContext(ctx, conf)
	grpcServer, gwmux := helpers.InitServerDefaults()
	address := server.GlobalConfig.Character.Local.Address()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	pb.RegisterHealthServiceServer(grpcServer, srv.NewHealthServiceServer())
	err := pb.RegisterHealthServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "register health service handler endpoint")

	css, err := srv.NewCharacterServiceServer(ctx, server)
	helpers.Check(ctx, err, "create characters service server")
	pb.RegisterCharacterServiceServer(grpcServer, css)
	err = pb.RegisterCharacterServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "registering characters service handler endpoint")

	helpers.StartServer(ctx, grpcServer, gwmux, server.GlobalConfig.Character.Local.Address())
}

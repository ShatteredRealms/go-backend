package main

import (
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/uptrace-go/uptrace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ShatteredRealms/go-backend/cmd/characters/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
)

var (
	conf        *config.GlobalSROConfig
	serviceName = "characters_service"
)

func init() {
	config.SetupLogger()
	conf = config.NewGlobalConfig()
}

func main() {
	ctx := context.Background()
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(conf.Uptrace.DSN()),
		uptrace.WithServiceName(serviceName),
		uptrace.WithServiceVersion(conf.Version),
	)
	defer func(ctx context.Context) {
		err := uptrace.Shutdown(ctx)
		if err != nil {
			log.WithContext(ctx).Errorf("shutdown uptrace: %v", err)
		}
	}(ctx)

	server := characters.NewServer(ctx, conf)
	grpcServer, gwmux := helpers.InitServerDefaults()
	address := server.GlobalConfig.Characters.Local.Address()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	pb.RegisterHealthServiceServer(grpcServer, srv.NewHealthServiceServer())
	err := pb.RegisterHealthServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "register health service handler endpoint")

	css, err := srv.NewCharactersServiceServer(server)
	helpers.Check(ctx, err, "create characters service server")
	pb.RegisterCharactersServiceServer(grpcServer, css)
	err = pb.RegisterCharactersServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "registering characters service handler endpoint")

	helpers.StartServer(ctx, grpcServer, gwmux, server.GlobalConfig.Characters.Local.Address())
}

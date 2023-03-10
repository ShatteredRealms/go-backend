package main

import (
	"context"
	chat "github.com/ShatteredRealms/go-backend/cmd/chat/app"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/uptrace/uptrace-go/uptrace"
)

var (
	conf        *config.GlobalSROConfig
	serviceName = "chat_service"
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

	server := chat.NewServer(ctx, conf)
	grpcServer, gwmux := helpers.InitServerDefaults()
	address := server.GlobalConfig.Chat.Local.Address()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	pb.RegisterHealthServiceServer(grpcServer, srv.NewHealthServiceServer())
	err := pb.RegisterHealthServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "register health service handler endpoint")

	pb.RegisterChatServiceServer(grpcServer, srv.NewChatServiceServer(server))
	err = pb.RegisterChatServiceHandlerFromEndpoint(ctx, gwmux, address, opts)
	helpers.Check(ctx, err, "register chat service handler endpoint")

	helpers.StartServer(ctx, grpcServer, gwmux, server.GlobalConfig.Characters.Local.Address())
}

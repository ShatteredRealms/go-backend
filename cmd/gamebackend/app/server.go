package gamebackend

import (
	"context"
	"fmt"

	"agones.dev/agones/pkg/client/clientset/versioned"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	"github.com/WilSimpson/gocloak/v13"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/client-go/rest"
)

var (
	ServiceName = "gamebackend"
)

type GameBackendServerContext struct {
	*config.ServerContext
	CharacterClient    pb.CharacterServiceClient
	ChatClient         pb.ChatServiceClient
	GamebackendService service.GamebackendService
	AgonesClient       versioned.Interface
}

func NewServerContext(ctx context.Context, conf *config.GlobalConfig, tracer trace.Tracer) (*GameBackendServerContext, error) {
	server := &GameBackendServerContext{
		ServerContext: &config.ServerContext{
			GlobalConfig:   conf,
			Tracer:         tracer,
			KeycloakClient: gocloak.NewClient(conf.Keycloak.BaseURL),
			RefSROServer:   &conf.GameBackend.SROServer,
		},
		CharacterClient:    nil,
		ChatClient:         nil,
		GamebackendService: nil,
		AgonesClient:       nil,
	}
	ctx, span := server.Tracer.Start(ctx, "server.new")
	defer span.End()

	server.KeycloakClient.RegisterMiddlewares(gocloak.OpenTelemetryMiddleware)

	charactersService, err := helpers.GrpcClientWithOtel(conf.Character.Remote.Address())
	if err != nil {
		return nil, fmt.Errorf("connecting to characters service: %w", err)
	}
	server.CharacterClient = pb.NewCharacterServiceClient(charactersService)

	chatService, err := helpers.GrpcClientWithOtel(conf.Chat.Remote.Address())
	if err != nil {
		return nil, fmt.Errorf("connecting to chat service: %w", err)
	}

	server.ChatClient = pb.NewChatServiceClient(chatService)

	if server.GlobalConfig.GameBackend.Mode != config.LocalMode {
		conf, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("creating k8s config: %w", err)
		}

		server.AgonesClient, err = versioned.NewForConfig(conf)
		if err != nil {
			return nil, fmt.Errorf("connecting to agones: %w", err)
		}
	}

	db, err := repository.ConnectDB(conf.GameBackend.Postgres, conf.Redis)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres database: %w", err)
	}

	repo := repository.NewGamebackendRepository(db)
	gamebackendService, err := service.NewGamebackendService(ctx, repo)
	if err != nil {
		return nil, fmt.Errorf("creating gamebackend service: %w", err)
	}
	server.GamebackendService = gamebackendService

	return server, nil
}

// func (s *GameBackendServerContext) dialAgonesAllocatorServer() (*grpc.ClientConn, error) {
// 	clientKey, err := os.ReadFile(s.GlobalConfig.Agones.KeyFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("reading agones key: %w", err)
// 	}
//
// 	clientCert, err := os.ReadFile(s.GlobalConfig.Agones.CertFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("reading agones cert: %w", err)
// 	}
//
// 	caCert, err := os.ReadFile(s.GlobalConfig.Agones.CaCertFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("reading agones ca cert: %w", err)
// 	}
//
// 	opt, err := createRemoteClusterDialOption(clientCert, clientKey, caCert)
// 	if err != nil {
// 		return nil, fmt.Errorf("creating agones dial option: %w", err)
// 	}
//
// 	return grpc.Dial(s.GlobalConfig.Agones.Allocator.Address(), opt)
// }
//
// // createRemoteClusterDialOption creates a grpc client dial option with TLS configuration.
// func createRemoteClusterDialOption(clientCert, clientKey, caCert []byte) (grpc.DialOption, error) {
// 	// Load client cert
// 	cert, err := tls.X509KeyPair(clientCert, clientKey)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
// 	if len(caCert) != 0 {
// 		// Load CA cert, if provided and trust the server certificate.
// 		// This is required for self-signed certs.
// 		tlsConfig.RootCAs = x509.NewCertPool()
// 		if !tlsConfig.RootCAs.AppendCertsFromPEM(caCert) {
// 			return nil, errors.New("only PEM format is accepted for server CA")
// 		}
// 	}
//
// 	return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)), nil
// }

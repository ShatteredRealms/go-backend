package gamebackend

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"agones.dev/agones/pkg/client/clientset/versioned"
	"github.com/Nerzal/gocloak/v13"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"k8s.io/client-go/rest"
)

var (
	ServiceName = "gamebackend"
)

type GameBackendServerContext struct {
	GlobalConfig       *config.GlobalConfig
	CharacterClient    pb.CharacterServiceClient
	ChatClient         pb.ChatServiceClient
	GamebackendService service.GamebackendService
	KeycloakClient     *gocloak.GoCloak
	AgonesClient       versioned.Interface
	Tracer             trace.Tracer
}

func NewServerContext(ctx context.Context, conf *config.GlobalConfig) (*GameBackendServerContext, error) {
	server := &GameBackendServerContext{
		GlobalConfig:   conf,
		Tracer:         otel.Tracer("GameBackendService"),
		KeycloakClient: gocloak.NewClient(conf.Keycloak.BaseURL),
	}

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

	db, err := repository.ConnectDB(conf.GameBackend.Postgres)
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

func (s *GameBackendServerContext) dialAgonesAllocatorServer() (*grpc.ClientConn, error) {
	clientKey, err := os.ReadFile(s.GlobalConfig.Agones.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("reading agones key: %w", err)
	}

	clientCert, err := os.ReadFile(s.GlobalConfig.Agones.CertFile)
	if err != nil {
		return nil, fmt.Errorf("reading agones cert: %w", err)
	}

	caCert, err := os.ReadFile(s.GlobalConfig.Agones.CaCertFile)
	if err != nil {
		return nil, fmt.Errorf("reading agones ca cert: %w", err)
	}

	opt, err := createRemoteClusterDialOption(clientCert, clientKey, caCert)
	if err != nil {
		return nil, fmt.Errorf("creating agones dial option: %w", err)
	}

	return grpc.Dial(s.GlobalConfig.Agones.Allocator.Address(), opt)
}

// createRemoteClusterDialOption creates a grpc client dial option with TLS configuration.
func createRemoteClusterDialOption(clientCert, clientKey, caCert []byte) (grpc.DialOption, error) {
	// Load client cert
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	if len(caCert) != 0 {
		// Load CA cert, if provided and trust the server certificate.
		// This is required for self-signed certs.
		tlsConfig.RootCAs = x509.NewCertPool()
		if !tlsConfig.RootCAs.AppendCertsFromPEM(caCert) {
			return nil, errors.New("only PEM format is accepted for server CA")
		}
	}

	return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)), nil
}

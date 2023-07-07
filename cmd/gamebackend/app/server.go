package gamebackend

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	aapb "agones.dev/agones/pkg/allocation/go"
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
)

var (
	ServiceName = "gamebackend"
)

type GameBackendServerContext struct {
	GlobalConfig       *config.GlobalConfig
	CharactersClient   pb.CharactersServiceClient
	GamebackendService service.GamebackendService
	AgonesClient       aapb.AllocationServiceClient
	KeycloakClient     *gocloak.GoCloak
	Tracer             trace.Tracer
}

func NewServerContext(ctx context.Context, conf *config.GlobalConfig) *GameBackendServerContext {
	server := &GameBackendServerContext{
		GlobalConfig:   conf,
		Tracer:         otel.Tracer("GameBackendService"),
		KeycloakClient: gocloak.NewClient(conf.GameBackend.Keycloak.BaseURL),
	}

	cc, err := helpers.GrpcClientWithOtel(conf.Characters.Remote.Address())
	helpers.Check(ctx, err, "connecting to characters")
	server.CharactersClient = pb.NewCharactersServiceClient(cc)

	if conf.GameBackend.Mode != config.LocalMode {
		ac, err := helpers.GrpcClientWithOtel(conf.Agones.Allocator.Address())
		helpers.Check(ctx, err, "connecting to agones")
		server.AgonesClient = aapb.NewAllocationServiceClient(ac)
	}

	db, err := repository.ConnectDB(conf.GameBackend.DB)
	helpers.Check(ctx, err, "connecting to database")

	repo := repository.NewGamebackendRepository(db)
	gamebackendService, err := service.NewGamebackendService(ctx, repo)
	helpers.Check(ctx, err, "gamebackend service")
	server.GamebackendService = gamebackendService

	return server
}

func (s *GameBackendServerContext) dialAgonesAllocatorServer() (*grpc.ClientConn, error) {
	clientKey, err := os.ReadFile(s.GlobalConfig.Agones.KeyFile)
	if err != nil {
		return nil, err
	}

	clientCert, err := os.ReadFile(s.GlobalConfig.Agones.CertFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(s.GlobalConfig.Agones.CaCertFile)
	if err != nil {
		return nil, err
	}

	opt, err := createRemoteClusterDialOption(clientCert, clientKey, caCert)
	if err != nil {
		return nil, err
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

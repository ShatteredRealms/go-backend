package main

import (
	aapb "agones.dev/agones/pkg/allocation/go"
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/config"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/srv"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

func NewServer(
	jwt service.JWTService,
) (*grpc.Server, *runtime.ServeMux, error) {
	ctx := context.Background()

	grpcServer, gwmux, opts, err := srv.CreateGrpcServerWithAuth(
		ctx,
		jwt,
		conf.Accounts.Remote.Address(),
		"gamebackend",
		map[string]struct{}{
			"/sro.gamebackend.ConnectionService/Connect": {},
		},
	)
	if err != nil {
		return nil, nil, err
	}

	localhostMode := conf.GameBackend.Mode == config.ModeDevelopment
	var allocator aapb.AllocationServiceClient
	if !localhostMode {
		client, err := dialAgonesAllocatorServer()
		if err != nil {
			return nil, nil, err
		}
		allocator = aapb.NewAllocationServiceClient(client)
	}

	charactersClient, err := srv.DialOtelGrpc(conf.Characters.Remote.Address())
	if err != nil {
		return nil, nil, err
	}
	characters := pb.NewCharactersServiceClient(charactersClient)

	connectionServiceServer := srv.NewConnectionServiceServer(
		jwt,
		allocator,
		characters,
		conf.Agones.Namespace,
		localhostMode,
	)

	pb.RegisterConnectionServiceServer(grpcServer, connectionServiceServer)
	err = pb.RegisterConnectionServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.GameBackend.Local.Address(),
		opts,
	)

	return grpcServer, gwmux, err
}

func dialAgonesAllocatorServer() (*grpc.ClientConn, error) {
	clientKey, err := ioutil.ReadFile(conf.Agones.KeyFile)
	if err != nil {
		return nil, err
	}

	clientCert, err := ioutil.ReadFile(conf.Agones.CertFile)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(conf.Agones.CaCertFile)
	if err != nil {
		return nil, err
	}

	opt, err := createRemoteClusterDialOption(clientCert, clientKey, caCert)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(conf.Agones.Allocator.Remote.Address(), opt)
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

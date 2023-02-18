package main

import (
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/kend/pkg/service"
	"github.com/kend/pkg/srv"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func NewServer(
	characterService service.CharacterService,
	jwt service.JWTService,
) (*grpc.Server, *runtime.ServeMux, error) {
	ctx := context.Background()

	grpcServer, gwmux, opts, err := srv.CreateGrpcServerWithAuth(
		ctx,
		jwt,
		conf.Accounts.Remote.Address(),
		"characters",
		nil,
	)

	characterServiceServer := srv.NewCharacterServiceServer(characterService, jwt)
	pb.RegisterCharactersServiceServer(grpcServer, characterServiceServer)
	err = pb.RegisterCharactersServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Characters.Local.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	return grpcServer, gwmux, nil
}

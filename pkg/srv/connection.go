package srv

import (
	"context"

	aapb "agones.dev/agones/pkg/allocation/go"
	gamebackend "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type connectionServiceServer struct {
	pb.UnimplementedConnectionServiceServer
	server *gamebackend.GameBackendServerContext
}

func (s connectionServiceServer) ConnectGameServer(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*pb.ConnectGameServerResponse, error) {
	if s.server.GlobalConfig.GameBackend.Mode == config.LocalMode {
		return &pb.ConnectGameServerResponse{
			Address: "127.0.0.1",
			Port:    7777,
		}, nil
	}

	// If the current user can't get the character, then deny the request
	character, err := s.server.CharactersClient.GetCharacter(
		ctx,
		request,
	)
	if err != nil || character == nil {
		log.WithContext(ctx).Errorf("unable to get character %v: %s", request.Target, err)
		return nil, err
	}

	world := "Scene_Demo"
	if character.Location != nil && character.Location.World != "" {
		world = character.Location.World
	}

	log.WithContext(ctx).Debugf("character world: %s", world)

	allocatorReq := &aapb.AllocationRequest{
		Namespace: s.server.GlobalConfig.Agones.Namespace,
		GameServerSelectors: []*aapb.GameServerSelector{
			{
				//MatchLabels: map[string]string{
				//	"world": world,
				//},
				GameServerState: aapb.GameServerSelector_ALLOCATED,
				Players: &aapb.PlayerSelector{
					MinAvailable: 1,
					MaxAvailable: 1000,
				},
			},
			{
				//MatchLabels: map[string]string{
				//	"world": world,
				//},
				GameServerState: aapb.GameServerSelector_READY,
				Players: &aapb.PlayerSelector{
					MinAvailable: 1,
					MaxAvailable: 1000,
				},
			},
		},
	}

	serverCtx := helpers.ContextAddClientAuth(ctx, s.server.GlobalConfig.GameBackend.Keycloak.ClientId, s.server.GlobalConfig.GameBackend.Keycloak.ClientSecret)
	allocatorResp, err := s.server.AgonesClient.Allocate(serverCtx, allocatorReq)
	if err != nil {
		log.WithContext(ctx).Errorf("allocating: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ConnectGameServerResponse{
		Address: allocatorResp.Address,
		Port:    uint32(allocatorResp.Ports[0].Port),
	}, nil
}

func NewConnectionServiceServer(server *gamebackend.GameBackendServerContext) pb.ConnectionServiceServer {
	return &connectionServiceServer{
		server: server,
	}
}

package srv

import (
	aapb "agones.dev/agones/pkg/allocation/go"
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	utilService "github.com/kend/pkg/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type connectionServiceServer struct {
	pb.UnimplementedConnectionServiceServer
	jwtService utilService.JWTService
	allocator  aapb.AllocationServiceClient
	characters pb.CharactersServiceClient

	// localhostMode used to tell whether to search for a server on kubernetes or return a constant localhost connection
	localhostMode bool

	// namespace kubernetes namespace to search for gameservers in
	namespace string
}

func NewConnectionServiceServer(
	jwtService utilService.JWTService,
	allocator aapb.AllocationServiceClient,
	characters pb.CharactersServiceClient,
	namespace string,
	localHostMode bool,
) pb.ConnectionServiceServer {

	return &connectionServiceServer{
		jwtService:    jwtService,
		allocator:     allocator,
		characters:    characters,
		localhostMode: localHostMode,
		namespace:     namespace,
	}
}

func (s *connectionServiceServer) ConnectGameServer(ctx context.Context, request *pb.ConnectGameServerRequest) (*pb.ConnectGameServerResponse, error) {
	if s.localhostMode {
		return &pb.ConnectGameServerResponse{
			Address: "127.0.0.1",
			Port:    7777,
		}, nil
	}

	// If the current user can't get the character, then deny the request
	//character, err := s.characters.GetCharacter(
	//    ctx,
	//    &pb.CharacterTarget{CharacterId: request.CharacterId},
	//)
	//if err != nil {
	//
	//    fmt.Println("err 1")
	//    return nil, err
	//}

	//world := "Scene_Demo"
	//if character.Location != nil && character.Location.World != "" {
	//    world = character.Location.World
	//}

	allocatorReq := &aapb.AllocationRequest{
		Namespace: s.namespace,
		GameServerSelectors: []*aapb.GameServerSelector{
			{
				//MatchLabels: map[string]string{
				//    "world": world,
				//},
				GameServerState: aapb.GameServerSelector_ALLOCATED,
				Players: &aapb.PlayerSelector{
					MinAvailable: 1,
					MaxAvailable: 1000,
				},
			},
			{
				//MatchLabels: map[string]string{
				//    "world": world,
				//},
				GameServerState: aapb.GameServerSelector_READY,
				Players: &aapb.PlayerSelector{
					MinAvailable: 1,
					MaxAvailable: 1000,
				},
			},
		},
	}

	allocatorResp, err := s.allocator.Allocate(serverAuthContext(ctx, s.jwtService, "sro.com/gamebackend/v1/"), allocatorReq)
	if err != nil {
		log.WithContext(ctx).Errorf("allocating: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ConnectGameServerResponse{
		Address: allocatorResp.Address,
		Port:    uint32(allocatorResp.Ports[0].Port),
	}, nil
}

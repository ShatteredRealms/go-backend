package srv

import (
	"context"

	// aapb "agones.dev/agones/pkg/allocation/go"
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	aav1 "agones.dev/agones/pkg/apis/allocation/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	gamebackend "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	// "github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type connectionServiceServer struct {
	pb.UnimplementedConnectionServiceServer
	server *gamebackend.GameBackendServerContext
}

func (s connectionServiceServer) ConnectGameServer(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*pb.ConnectGameServerResponse, error) {

	// If the current user can't get the character, then deny the request
	character, err := s.server.CharactersClient.GetCharacter(
		helpers.PassAuthContext(ctx),
		request,
	)
	if err != nil || character == nil {
		log.WithContext(ctx).Errorf("unable to get character %v: %s", request.Target, err)
		return nil, err
	}

	if s.server.GlobalConfig.GameBackend.Mode == config.LocalMode {
		log.WithContext(ctx).Debug("using local game server")
		return &pb.ConnectGameServerResponse{
			Address: "127.0.0.1",
			Port:    7777,
		}, nil
	}

	world := "Scene_Demo"
	if character.Location != nil && character.Location.World != "" {
		world = character.Location.World
	}

	log.WithContext(ctx).Debugf("character name %s connecting to world: %s", character.Name, world)

	// allocatorReq := &aapb.AllocationRequest{
	// 	Namespace: s.server.GlobalConfig.Agones.Namespace,
	// 	GameServerSelectors: []*aapb.GameServerSelector{
	// 		{
	// 			//MatchLabels: map[string]string{
	// 			//	"world": world,
	// 			//},
	// 			GameServerState: aapb.GameServerSelector_ALLOCATED,
	// 			Players: &aapb.PlayerSelector{
	// 				MinAvailable: 1,
	// 				MaxAvailable: 1000,
	// 			},
	// 		},
	// 		{
	// 			//MatchLabels: map[string]string{
	// 			//	"world": world,
	// 			//},
	// 			GameServerState: aapb.GameServerSelector_READY,
	// 			Players: &aapb.PlayerSelector{
	// 				MinAvailable: 1,
	// 				MaxAvailable: 1000,
	// 			},
	// 		},
	// 	},
	// }

	srvCtx, err := s.serverContext(ctx)
	// if err != nil {
	// 	log.WithContext(ctx).Errorf("create server context: %v", err)
	// 	return nil, model.ErrHandleRequest
	// }
	// allocatorResp, err := s.server.AgonesClient.Allocate(srvCtx, allocatorReq)
	// if err != nil {
	// 	log.WithContext(ctx).Errorf("allocating: %v", err)
	// 	return nil, status.Error(codes.Internal, err.Error())
	// }

	conf, err := rest.InClusterConfig()
	if err != nil {
		log.WithContext(ctx).Errorf("creating config: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	agones, err := versioned.NewForConfig(conf)
	if err != nil {
		log.WithContext(ctx).Errorf("creating agones connection: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	allocatedState := v1.GameServerStateAllocated
	readyState := v1.GameServerStateReady
	resp, err := agones.AllocationV1().GameServerAllocations(s.server.GlobalConfig.Agones.Namespace).Create(
		srvCtx,
		&aav1.GameServerAllocation{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: aav1.GameServerAllocationSpec{
				Selectors: []aav1.GameServerSelector{
					{
						GameServerState: &allocatedState,
						Players: &aav1.PlayerSelector{
							MinAvailable: 1,
							MaxAvailable: 1000,
						},
					},
					{
						GameServerState: &readyState,
						Players: &aav1.PlayerSelector{
							MinAvailable: 1,
							MaxAvailable: 1000,
						},
					},
				},
			},
		},
		metav1.CreateOptions{},
	)
	if err != nil {
		log.WithContext(ctx).Errorf("allocation request: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ConnectGameServerResponse{
		Address: resp.Status.Address,
		Port:    uint32(resp.Status.Ports[0].Port),
	}, nil
}

func NewConnectionServiceServer(server *gamebackend.GameBackendServerContext) pb.ConnectionServiceServer {
	return &connectionServiceServer{
		server: server,
	}
}

func (s connectionServiceServer) serverContext(ctx context.Context) (context.Context, error) {
	token, err := s.server.KeycloakClient.LoginClient(
		ctx,
		s.server.GlobalConfig.Chat.Keycloak.ClientId,
		s.server.GlobalConfig.Chat.Keycloak.ClientSecret,
		s.server.GlobalConfig.Chat.Keycloak.Realm,
	)
	if err != nil {
		return nil, err
	}

	return helpers.ContextAddClientToken(
		ctx,
		token.AccessToken,
	), nil
}

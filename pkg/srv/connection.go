package srv

import (
	"context"
	"fmt"

	// aapb "agones.dev/agones/pkg/allocation/go"
	v1 "agones.dev/agones/pkg/apis/agones/v1"
	aav1 "agones.dev/agones/pkg/apis/allocation/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	gamebackend "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/google/uuid"

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
	agones *versioned.Clientset
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
		pc, err := s.server.GamebackendService.CreatePendingConnection(ctx, character.Name, "localhost")
		if err != nil {
			return nil, fmt.Errorf("create pending connection: %w", err)
		}

		log.WithContext(ctx).Debugf("%s using local game server", character.Name)
		return &pb.ConnectGameServerResponse{
			Address:      "127.0.0.1",
			Port:         7777,
			ConnectionId: pc.Id.String(),
		}, nil
	}

	// Check if player is playing
	out, err := s.agones.AgonesV1().GameServers("sro").List(ctx, metav1.ListOptions{})
	for _, gs := range out.Items {
		for _, pId := range gs.Status.Players.IDs {
			if character.GetOwner() == pId {

			}
		}
	}

	world := "Scene_Demo"
	if character.Location != nil && character.Location.World != "" {
		world = character.Location.World
	}

	log.WithContext(ctx).Debugf("%s requesting connection to gameserver with world %s", character.Name, world)

	// Request allocation
	srvCtx, err := s.serverContext(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("create server context: %v", err)
		return nil, model.ErrHandleRequest
	}

	allocatedState := v1.GameServerStateAllocated
	readyState := v1.GameServerStateReady
	resp, err := s.agones.AllocationV1().GameServerAllocations(s.server.GlobalConfig.Agones.Namespace).Create(
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

	pc, err := s.server.GamebackendService.CreatePendingConnection(ctx, character.Name, resp.Status.NodeName)
	if err != nil {
		return nil, fmt.Errorf("create pending connection: %w", err)
	}

	return &pb.ConnectGameServerResponse{
		Address:      resp.Status.Address,
		Port:         uint32(resp.Status.Ports[0].Port),
		ConnectionId: pc.Id.String(),
	}, nil
}

func (s connectionServiceServer) VerifyConnect(
	ctx context.Context,
	request *pb.VerifyConnectRequest,
) (*pb.CharacterDetails, error) {
	id, err := uuid.FromBytes([]byte(request.ConnectionId))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	pc, err := s.server.GamebackendService.CheckPlayerConnection(ctx, &id, request.ServerName)
	if err != nil {
		return nil, status.Error(codes.OK, err.Error())
	}

	// If the current user can't get the character, then deny the request
	character, err := s.server.CharactersClient.GetCharacter(
		helpers.PassAuthContext(ctx),
		&pb.CharacterTarget{
			Target: &pb.CharacterTarget_Name{
				Name: pc.Character,
			},
		},
	)
	if err != nil || character == nil {
		log.WithContext(ctx).Errorf("unable to get character %v: %s", pc.Character, err)
		return nil, status.Errorf(codes.Internal, "unable to find character")
	}

	return character, nil
}

func NewConnectionServiceServer(
	ctx context.Context,
	server *gamebackend.GameBackendServerContext,
) (pb.ConnectionServiceServer, error) {
	if server.GlobalConfig.GameBackend.Mode != config.LocalMode {
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

		return &connectionServiceServer{
			server: server,
			agones: agones,
		}, nil
	}

	return &connectionServiceServer{
		server: server,
	}, nil
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

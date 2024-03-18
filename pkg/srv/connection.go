package srv

import (
	"context"
	"fmt"

	v1 "agones.dev/agones/pkg/apis/agones/v1"
	aav1 "agones.dev/agones/pkg/apis/allocation/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	"github.com/Nerzal/gocloak/v13"
	gamebackend "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var (
	ConnectionRoles = make([]*gocloak.Role, 0)

	RoleConnect = registerConnectionRole(&gocloak.Role{
		Name:        gocloak.StringP("connect"),
		Description: gocloak.StringP("Allows requests to connect to a server"),
	})

	RoleManageConnections = registerConnectionRole(&gocloak.Role{
		Name:        gocloak.StringP("manage_connections"),
		Description: gocloak.StringP("Allows verifying and transfering connections"),
	})
)

func registerConnectionRole(role *gocloak.Role) *gocloak.Role {
	ConnectionRoles = append(ConnectionRoles, role)
	return role
}

type connectionServiceServer struct {
	pb.UnimplementedConnectionServiceServer
	server *gamebackend.GameBackendServerContext
	agones *versioned.Clientset
}

func (s connectionServiceServer) ConnectGameServer(
	ctx context.Context,
	request *pb.CharacterTarget,
) (*pb.ConnectGameServerResponse, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.GameBackend.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleConnect, model.GamebackendClientId) {
		return nil, errors.Wrapf(model.ErrUnauthorized, "no role %s", *RoleConnect.Name)
	}

	// If the current user can't get the character, then deny the request
	character, err := s.server.CharacterClient.GetCharacter(
		helpers.PassAuthContext(ctx),
		request,
	)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("unable to get character %v: %s", request.Type, err)
		return nil, err
	}
	if character == nil {
		log.Logger.WithContext(ctx).Warnf("%s requested character %v but does not exists", claims.Username, request.Type)
		return nil, model.ErrDoesNotExist
	}

	if s.server.GlobalConfig.GameBackend.Mode == config.LocalMode {
		return s.requestLocalConnection(ctx, character.Name)
	}

	// Check if player is playing
	out, err := s.agones.AgonesV1().GameServers("sro").List(ctx, metav1.ListOptions{})
	for _, gs := range out.Items {
		for _, pId := range gs.Status.Players.IDs {
			if character.GetOwner() == pId {
				return nil, status.Error(codes.FailedPrecondition, "character already playing")
			}
		}
	}

	// Validate location. First time characters can currently be nil.
	// @TODO: Make this check unnecessary by having default character location. Add to validation check.
	if character.Location == nil {
		character.Location = &pb.Location{
			World: "Scene_Demo",
		}
	} else if character.Location.World == "" {
		character.Location.World = "Scene_Demo"
	}

	return s.requestConnection(ctx, character, character.Location, false)
}

func (s connectionServiceServer) VerifyConnect(
	ctx context.Context,
	request *pb.VerifyConnectRequest,
) (*pb.CharacterDetails, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.GameBackend.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleManageConnections, model.GamebackendClientId) {
		return nil, model.ErrUnauthorized
	}

	id, err := uuid.Parse(request.ConnectionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %s", request.ConnectionId)
	}

	pc, err := s.server.GamebackendService.CheckPlayerConnection(ctx, &id, request.ServerName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// If the current user can't get the character, then deny the request
	character, err := s.server.CharacterClient.GetCharacter(
		helpers.PassAuthContext(ctx),
		&pb.CharacterTarget{
			Type: &pb.CharacterTarget_Name{
				Name: pc.Character,
			},
		},
	)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("unable to get character %v: %s", pc.Character, err)
		return nil, err
	}
	if character == nil {
		log.Logger.WithContext(ctx).Errorf("character not found %v: %s", pc.Character, err)
		return nil, status.Errorf(codes.Internal, "unable to find character")
	}

	return character, nil
}

func (s connectionServiceServer) TransferPlayer(
	ctx context.Context,
	request *pb.TransferPlayerRequest,
) (*pb.ConnectGameServerResponse, error) {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.GameBackend.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return nil, model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleManageConnections, model.GamebackendClientId) {
		return nil, model.ErrUnauthorized
	}

	character, err := s.server.CharacterClient.GetCharacter(
		helpers.PassAuthContext(ctx),
		&pb.CharacterTarget{
			Type: &pb.CharacterTarget_Name{
				Name: request.Character,
			},
		},
	)
	if err != nil || character == nil {
		log.Logger.WithContext(ctx).Errorf("unable to get character %v: %s", request.Character, err)
		return nil, status.Errorf(codes.Internal, "unable to find character")
	}

	if s.server.GlobalConfig.GameBackend.Mode == config.LocalMode {
		return s.requestLocalConnection(ctx, character.Name)
	}

	return s.requestConnection(ctx, character, request.Location, true)
}

func (s connectionServiceServer) requestLocalConnection(
	ctx context.Context,
	characterName string,
) (*pb.ConnectGameServerResponse, error) {
	pc, err := s.server.GamebackendService.CreatePendingConnection(ctx, characterName, "localhost")
	if err != nil {
		return nil, fmt.Errorf("create pending connection: %w", err)
	}

	log.Logger.WithContext(ctx).Debugf("%s using local game server", characterName)
	return &pb.ConnectGameServerResponse{
		Address:      "127.0.0.1",
		Port:         7777,
		ConnectionId: pc.Id.String(),
	}, nil
}

func (s connectionServiceServer) requestConnection(
	ctx context.Context,
	character *pb.CharacterDetails,
	location *pb.Location,
	updateCharacter bool,
) (*pb.ConnectGameServerResponse, error) {
	log.Logger.WithContext(ctx).Debugf("%s requesting connection to gameserver with world %s", character.Name, location.World)
	// Request allocation
	srvCtx, err := s.serverContext(ctx)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("create server context: %v", err)
		return nil, model.ErrHandleRequest
	}

	allocatedState := v1.GameServerStateAllocated
	readyState := v1.GameServerStateReady
	gsAlloc, err := s.agones.AllocationV1().GameServerAllocations(s.server.GlobalConfig.Agones.Namespace).Create(
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
						LabelSelector: metav1.LabelSelector{
							MatchLabels: map[string]string{
								"world": location.World,
							},
						},
					},
					{
						GameServerState: &readyState,
						Players: &aav1.PlayerSelector{
							MinAvailable: 1,
							MaxAvailable: 1000,
						},
						LabelSelector: metav1.LabelSelector{
							MatchLabels: map[string]string{
								"world": location.World,
							},
						},
					},
				},
			},
		},
		metav1.CreateOptions{},
	)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("allocation request: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if updateCharacter {
		// Update the players location
		_, err = s.server.CharacterClient.EditCharacter(
			helpers.PassAuthContext(ctx),
			&pb.EditCharacterRequest{
				Target: &pb.CharacterTarget{
					Type: &pb.CharacterTarget_Id{Id: character.Id},
				},
				OptionalLocation: &pb.EditCharacterRequest_Location{
					Location: location,
				},
			},
		)
		if err != nil {
			log.Logger.WithContext(ctx).Errorf("updating character location: %s", err.Error())
			return nil, fmt.Errorf("updating character location: %w", err)
		}

	}

	// Create pending connection
	pc, err := s.server.GamebackendService.CreatePendingConnection(ctx, character.Name, gsAlloc.Status.NodeName)
	if err != nil {
		return nil, fmt.Errorf("create pending connection: %w", err)
	}

	return &pb.ConnectGameServerResponse{
		Address:      gsAlloc.Status.Address,
		Port:         uint32(gsAlloc.Status.Ports[0].Port),
		ConnectionId: pc.Id.String(),
	}, nil

}

func NewConnectionServiceServer(
	ctx context.Context,
	server *gamebackend.GameBackendServerContext,
) (pb.ConnectionServiceServer, error) {
	token, err := server.KeycloakClient.LoginClient(
		ctx,
		server.GlobalConfig.GameBackend.Keycloak.ClientId,
		server.GlobalConfig.GameBackend.Keycloak.ClientSecret,
		server.GlobalConfig.GameBackend.Keycloak.Realm,
	)
	if err != nil {
		return nil, fmt.Errorf("login keycloak: %v", err)
	}

	err = createRoles(ctx,
		server.KeycloakClient,
		token.AccessToken,
		server.GlobalConfig.GameBackend.Keycloak.Realm,
		server.GlobalConfig.GameBackend.Keycloak.Id,
		&ConnectionRoles,
	)
	if err != nil {
		return nil, err
	}

	if server.GlobalConfig.GameBackend.Mode != config.LocalMode {
		conf, err := rest.InClusterConfig()
		if err != nil {
			log.Logger.WithContext(ctx).Errorf("creating config: %v", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		agones, err := versioned.NewForConfig(conf)
		if err != nil {
			log.Logger.WithContext(ctx).Errorf("creating agones connection: %v", err)
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
		s.server.GlobalConfig.GameBackend.Keycloak.ClientId,
		s.server.GlobalConfig.GameBackend.Keycloak.ClientSecret,
		s.server.GlobalConfig.GameBackend.Keycloak.Realm,
	)
	if err != nil {
		return nil, err
	}

	return helpers.ContextAddClientToken(
		ctx,
		token.AccessToken,
	), nil
}

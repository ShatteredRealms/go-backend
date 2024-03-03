package srv

import (
	context "context"
	"fmt"
	"reflect"
	"strings"

	v1 "agones.dev/agones/pkg/apis/agones/v1"
	autoscalingv1 "agones.dev/agones/pkg/apis/autoscaling/v1"
	"agones.dev/agones/pkg/client/clientset/versioned"
	"github.com/Nerzal/gocloak/v13"
	gamebackend "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/google/uuid"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
)

var (
	ServerManagerRoles = make([]*gocloak.Role, 0)

	RoleServerManager = registerServerManagerRole(&gocloak.Role{
		Name:        gocloak.StringP("server_manager"),
		Description: gocloak.StringP("Allow the use of the server manager service"),
	})

	RoleServerStatistics = registerServerManagerRole(&gocloak.Role{
		Name:        gocloak.StringP("server_statistics"),
		Description: gocloak.StringP("Allow viewing of gameserver status and statistics"),
	})

	ErrNoAgonesConnect = fmt.Errorf("not connected to agones")
)

func registerServerManagerRole(role *gocloak.Role) *gocloak.Role {
	ConnectionRoles = append(ConnectionRoles, role)
	return role
}

type serverManagerServiceServer struct {
	pb.UnimplementedServerManagerServiceServer
	server *gamebackend.GameBackendServerContext
	agones *versioned.Clientset
}

// CreateDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) CreateDimension(
	ctx context.Context,
	request *pb.CreateDimensionRequest,
) (*pb.Dimension, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	mapIds, err := helpers.ParseUUIDs(request.MapIds)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dimension, err := s.server.GamebackendService.CreateDimension(
		ctx,
		request.Name,
		request.Location,
		request.Version,
		mapIds,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "creating: %s", err.Error())
	}

	err = s.setupNewDimension(ctx, dimension)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "setup dimension: %s", err.Error())
	}

	return dimension.ToPb(), nil
}

// CreateMap implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) CreateMap(
	ctx context.Context,
	request *pb.CreateMapRequest,
) (*pb.Map, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	m, err := s.server.GamebackendService.CreateMap(
		ctx,
		request.Name,
		request.Path,
		request.MaxPlayers,
		request.Instanced,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "creating: %s", err.Error())
	}

	return m.ToPb(), nil
}

// DeleteDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) DeleteDimension(
	ctx context.Context,
	request *pb.DimensionTarget,
) (*emptypb.Empty, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	dimension, err := s.findDimensionByNameOrId(ctx, request)
	if err != nil {
		return nil, err
	}

	err = s.server.GamebackendService.DeleteDimensionById(ctx, dimension.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "deleteing dimension %s: %s", dimension.Name, err.Error())
	}

	for _, m := range dimension.Maps {
		err = s.deleteGameServers(ctx, dimension, m)
		if err != nil {
			return nil,
				status.Errorf(codes.Internal, "deleting game servers for dimension %s, map %s: %s", dimension.Name, m.Name, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}

// DeleteMap implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) DeleteMap(
	ctx context.Context,
	request *pb.MapTarget,
) (*emptypb.Empty, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	m, err := s.findMapByNameOrId(ctx, request)
	if err != nil {
		return nil, err
	}

	err = s.server.GamebackendService.DeleteMapById(ctx, m.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "deleteing map %s: %s", m.Name, err.Error())
	}

	dimensions, err := s.server.GamebackendService.FindDimensionsWithMapIds(ctx, []*uuid.UUID{m.Id})

	for _, dimension := range dimensions {
		err = s.deleteGameServers(ctx, dimension, m)
		if err != nil {
			return nil,
				status.Errorf(codes.Internal, "deleting game servers for dimension %s, map %s: %s", dimension.Name, m.Name, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}

// DuplicateDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) DuplicateDimension(
	ctx context.Context,
	request *pb.DuplicateDimensionRequest,
) (*pb.Dimension, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	var newDimension *model.Dimension
	switch target := request.Target.FindBy.(type) {
	case *pb.DimensionTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, err
		}

		newDimension, err = s.server.GamebackendService.DuplicateDimension(ctx, &id, request.Name)

	case *pb.DimensionTarget_Name:
		dimension, err := s.server.GamebackendService.FindDimensionByName(ctx, target.Name)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		newDimension, err = s.server.GamebackendService.DuplicateDimension(ctx, dimension.Id, request.Name)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if newDimension == nil {
		return nil, status.Error(codes.Internal, "failed to create new dimension")
	}

	err = s.setupNewDimension(ctx, newDimension)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "setup dimension: %s", err.Error())
	}

	return newDimension.ToPb(), err
}

// EditDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) EditDimension(
	ctx context.Context,
	request *pb.EditDimensionRequest,
) (*pb.Dimension, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	originalDimension, err := s.findDimensionByNameOrId(ctx, request.Target)
	if err != nil {
		return nil, err
	}

	editedDimension, err := s.server.GamebackendService.EditDimension(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Name change requires delete and recreate
	if request.OptionalName != nil {
		for _, m := range originalDimension.Maps {
			err := s.deleteGameServers(ctx, originalDimension, m)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "deleting old game server for %s - %s failed: %s",
					originalDimension.Name,
					m.Name,
					err.Error(),
				)
			}
		}

		// Only recreate here if maps didn't change
		if !request.EditMaps {
			err = s.setupNewDimension(ctx, editedDimension)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "setting up dimension gameservers: %s", err.Error())
			}
		}
	} else {
		// Everything wasn't deleted so we can safely change one-by-one
		if request.EditMaps {
			// Create a map of new maps based off of their Id
			newMaps := make(map[*uuid.UUID]*model.Map, len(request.MapIds))
			for _, m := range editedDimension.Maps {
				newMaps[m.Id] = m
			}

			// Loop original maps
			for _, m := range originalDimension.Maps {
				if _, ok := newMaps[m.Id]; !ok {
					// original had the map, but not the new one
					err := s.deleteGameServers(ctx, editedDimension, m)
					if err != nil {
						return nil, status.Errorf(codes.Internal, "unable to delete old gameserver world %s: %s", m.Name, err.Error())
					}
				} else {
					// original dimension has this map, so it's not new.
					delete(newMaps, m.Id)

					// UPdate the version
					if request.OptionalVersion != nil {
						err := s.updateGameServers(ctx, editedDimension, m)
						if err != nil {
							return nil, status.Errorf(codes.Internal, "unable to update gameserver world %s: %s", m.Name, err.Error())
						}
					}
				}
			}

			// newMaps now only contains map that weren't in the original
			for _, newMap := range newMaps {
				err = s.createGameServers(ctx, editedDimension, newMap)
				if err != nil {
					return nil, status.Errorf(codes.Internal, "unable to delete old gameserver world %s: %s", newMap.Name, err.Error())
				}
			}
		} else if request.OptionalVersion != nil {
			// Maps weren't changed, so update the versions
			for _, m := range editedDimension.Maps {
				err = s.createGameServers(ctx, editedDimension, m)
				if err != nil {
					return nil, status.Errorf(codes.Internal, "unable to update gameserver world %s: %s", m.Name, err.Error())
				}
			}
		}

	}

	return editedDimension.ToPb(), nil
}

// EditMap implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) EditMap(
	ctx context.Context,
	request *pb.EditMapRequest,
) (*pb.Map, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	originalMap, err := s.findMapByNameOrId(ctx, request.Target)
	if err != nil {
		return nil, err
	}

	editedMap, err := s.server.GamebackendService.EditMap(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	dimensions, err := s.server.GamebackendService.FindDimensionsWithMapIds(ctx, []*uuid.UUID{editedMap.Id})
	if request.OptionalName != nil {
		// Need to delete and recreate
		for _, dimension := range dimensions {
			err = s.deleteGameServers(ctx, dimension, originalMap)
			if err != nil {
				return nil,
					status.Errorf(codes.Internal, "deleting game servers for dimension %s: %s", dimension.Name, err.Error())
			}

			err = s.createGameServers(ctx, dimension, editedMap)
			if err != nil {
				return nil,
					status.Errorf(codes.Internal, "creating game servers for dimension %s: %s", dimension.Name, err.Error())
			}
		}
	} else {
		if request.OptionalPath != nil {
			for _, dimension := range dimensions {
				err = s.updateGameServers(ctx, dimension, editedMap)
				if err != nil {
					return nil,
						status.Errorf(codes.Internal, "updating game servers for dimension %s: %s", dimension.Name, err.Error())
				}
			}
		}
	}

	return editedMap.ToPb(), nil
}

// GetAllDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetAllDimension(
	ctx context.Context,
	request *emptypb.Empty,
) (*pb.Dimensions, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	out, err := s.server.GamebackendService.FindAllDimensions(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Dimensions{Dimensions: out.ToPb()}, nil
}

// GetAllMaps implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetAllMaps(
	ctx context.Context,
	request *emptypb.Empty,
) (*pb.Maps, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	out, err := s.server.GamebackendService.FindAllMaps(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Maps{Maps: out.ToPb()}, nil
}

// GetDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetDimension(
	ctx context.Context,
	request *pb.DimensionTarget,
) (*pb.Dimension, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	var out *model.Dimension
	switch target := request.FindBy.(type) {
	case *pb.DimensionTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, err
		}
		out, err = s.server.GamebackendService.FindDimensionById(ctx, &id)

	case *pb.DimensionTarget_Name:
		out, err = s.server.GamebackendService.FindDimensionByName(ctx, target.Name)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if out == nil {
		return nil, model.ErrDoesNotExist
	}

	return out.ToPb(), nil
}

// GetMap implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetMap(
	ctx context.Context,
	request *pb.MapTarget,
) (*pb.Map, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	var out *model.Map
	switch target := request.FindBy.(type) {
	case *pb.MapTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, err
		}
		out, err = s.server.GamebackendService.FindMapById(ctx, &id)

	case *pb.MapTarget_Name:
		out, err = s.server.GamebackendService.FindMapByName(ctx, target.Name)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if out == nil {
		return nil, model.ErrDoesNotExist
	}

	return out.ToPb(), nil
}

func NewServerManagerServiceServer(
	ctx context.Context,
	server *gamebackend.GameBackendServerContext,
) (pb.ServerManagerServiceServer, error) {
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

		return &serverManagerServiceServer{
			server: server,
			agones: agones,
		}, nil
	}

	return &serverManagerServiceServer{
		server: server,
	}, nil
}

func (s serverManagerServiceServer) serverContext(ctx context.Context) (context.Context, error) {
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

func (s serverManagerServiceServer) hasServerManagerRole(ctx context.Context) error {
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.GameBackend.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleServerManager, model.GamebackendClientId) {
		return model.ErrUnauthorized
	}

	return nil
}

func (s serverManagerServiceServer) createGameServers(
	ctx context.Context,
	dimension *model.Dimension,
	m *model.Map,
) error {
	if s.agones == nil {
		if s.server.GlobalConfig.GameBackend.Mode != config.LocalMode {
			return ErrNoAgonesConnect
		}

		log.Logger.WithContext(ctx).Infof("Local Mode: Not creating game server %s-%s", dimension.Name, m.Name)
		return nil
	}

	// Create the fleet
	fleet, err := s.agones.AgonesV1().Fleets(s.server.GlobalConfig.Agones.Namespace).Create(
		ctx,
		buildFleet(dimension, m, s.server.GlobalConfig.Agones.Namespace),
		metav1.CreateOptions{},
	)
	if err != nil {
		return fmt.Errorf("creating fleet: %w", err)
	}

	// Create autoscaler
	_, err = s.agones.AutoscalingV1().FleetAutoscalers(fleet.Namespace).Create(
		ctx,
		buildAutoscalingFleet(dimension, m, fleet.Namespace),
		metav1.CreateOptions{},
	)
	if err != nil {
		return fmt.Errorf("creating fleet autoscaler: %w", err)
	}

	return nil
}

func (s serverManagerServiceServer) deleteGameServers(
	ctx context.Context,
	dimension *model.Dimension,
	m *model.Map,
) error {
	if s.agones == nil {
		if s.server.GlobalConfig.GameBackend.Mode != config.LocalMode {
			return ErrNoAgonesConnect
		}

		log.Logger.WithContext(ctx).Infof("Local Mode: Not creating game server %s-%s", dimension.Name, m.Name)
		return nil
	}

	namespace := s.server.GlobalConfig.Agones.Namespace

	// Delete autoscaler
	err := s.agones.AutoscalingV1().FleetAutoscalers(namespace).Delete(
		ctx,
		getFleetAutoscalerName(dimension, m),
		metav1.DeleteOptions{},
	)

	if err != nil &&
		!errors.IsNotFound(err) {
		return fmt.Errorf("deleting fleet autoscaler: %w", err)
	}

	// Delete autoscaler
	err = s.agones.AgonesV1().Fleets(namespace).Delete(
		ctx,
		getFleetName(dimension, m),
		metav1.DeleteOptions{},
	)
	if err != nil &&
		!errors.IsNotFound(err) {
		return fmt.Errorf("deleting fleet: %w", err)
	}

	return nil
}

func (s serverManagerServiceServer) updateGameServers(
	ctx context.Context,
	dimension *model.Dimension,
	m *model.Map,
) error {
	if s.agones == nil {
		if s.server.GlobalConfig.GameBackend.Mode != config.LocalMode {
			return ErrNoAgonesConnect
		}

		log.Logger.WithContext(ctx).Infof("Local Mode: Not creating game server %s-%s", dimension.Name, m.Name)
		return nil
	}

	// Update the fleet
	fleet, err := s.agones.AgonesV1().Fleets(s.server.GlobalConfig.Agones.Namespace).Update(
		ctx,
		buildFleet(dimension, m, s.server.GlobalConfig.Agones.Namespace),
		metav1.UpdateOptions{},
	)
	if err != nil {
		return fmt.Errorf("creating fleet: %w", err)
	}

	// Update autoscalergones not setup, not connected in local mode
	_, err = s.agones.AutoscalingV1().FleetAutoscalers(fleet.Namespace).Update(
		ctx,
		buildAutoscalingFleet(dimension, m, fleet.Namespace),
		metav1.UpdateOptions{},
	)
	if err != nil {
		return fmt.Errorf("creating fleet autoscaler: %w", err)
	}

	return nil
}

func buildAutoscalingFleet(dimension *model.Dimension, m *model.Map, namespace string) *autoscalingv1.FleetAutoscaler {
	return &autoscalingv1.FleetAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getFleetAutoscalerName(dimension, m),
			Namespace: namespace,
		},
		Spec: autoscalingv1.FleetAutoscalerSpec{
			FleetName: getFleetName(dimension, m),
			Policy: autoscalingv1.FleetAutoscalerPolicy{
				Type: autoscalingv1.BufferPolicyType,
				Buffer: &autoscalingv1.BufferPolicy{
					MaxReplicas: 10,
					MinReplicas: 2,
					BufferSize: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 2,
					},
				},
			},
		},
	}
}

func buildFleet(dimension *model.Dimension, m *model.Map, namespace string) *v1.Fleet {
	// Create the starting arguments
	startArgs := make([]string, 2)

	// Map name to load
	startArgs[0] = m.Path

	// Enable logging
	startArgs[1] = "-log"

	return &v1.Fleet{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getFleetName(dimension, m),
			Namespace: namespace,
		},
		Spec: v1.FleetSpec{
			Replicas:   1,
			Scheduling: "",
			Template: v1.GameServerTemplateSpec{
				Spec: v1.GameServerSpec{
					Container: "",
					Ports: []v1.GameServerPort{
						{
							Name:          "default",
							PortPolicy:    "Dynamic",
							ContainerPort: 7777,
						},
					},
					Health: v1.Health{
						Disabled:            false,
						PeriodSeconds:       10,
						FailureThreshold:    3,
						InitialDelaySeconds: 300,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name:      getFleetName(dimension, m),
							Namespace: namespace,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            "gameserver",
									Image:           dimension.GetImageName(),
									Args:            startArgs,
									ImagePullPolicy: "Always",
								},
							},
							ImagePullSecrets: []corev1.LocalObjectReference{
								{
									Name: "regcred",
								},
							},
						},
					},
				},
			},
		},
	}
}

func getBaseFleetName(dimension *model.Dimension, m *model.Map) string {
	return fmt.Sprintf("%s-%s",
		strings.ToLower(dimension.Name),
		strings.ToLower(m.Name),
	)
}

func getFleetName(dimension *model.Dimension, m *model.Map) string {
	return fmt.Sprintf("fleet-%s", getBaseFleetName(dimension, m))
}

func getFleetAutoscalerName(dimension *model.Dimension, m *model.Map) string {
	return fmt.Sprintf("fleet-autoscaler-%s", getBaseFleetName(dimension, m))
}

func (s serverManagerServiceServer) setupNewDimension(ctx context.Context, dimension *model.Dimension) error {
	var err error
	for _, m := range dimension.Maps {
		err = s.createGameServers(ctx, dimension, m)
		if err != nil {
			return err
		}
	}

	if s.server.GlobalConfig.GameBackend.Mode == config.LocalMode {
		log.Logger.WithContext(ctx).Infof("agones not setup, not connected in local mode")
	}

	return nil
}

func (s serverManagerServiceServer) findMapByNameOrId(
	ctx context.Context,
	requestTarget *pb.MapTarget,
) (out *model.Map, err error) {
	switch target := requestTarget.FindBy.(type) {
	case *pb.MapTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, err
		}
		out, err = s.server.GamebackendService.FindMapById(ctx, &id)

	case *pb.MapTarget_Name:
		out, err = s.server.GamebackendService.FindMapByName(ctx, target.Name)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if out == nil {
		return nil, model.ErrDoesNotExist
	}

	return out, nil
}

func (s serverManagerServiceServer) findDimensionByNameOrId(
	ctx context.Context,
	requestTarget *pb.DimensionTarget,
) (out *model.Dimension, err error) {
	switch target := requestTarget.FindBy.(type) {
	case *pb.DimensionTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, err
		}
		out, err = s.server.GamebackendService.FindDimensionById(ctx, &id)

	case *pb.DimensionTarget_Name:
		out, err = s.server.GamebackendService.FindDimensionByName(ctx, target.Name)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if out == nil {
		return nil, model.ErrDoesNotExist
	}

	return out, nil
}

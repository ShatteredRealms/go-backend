package srv

import (
	context "context"
	"errors"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"

	v1 "agones.dev/agones/pkg/apis/agones/v1"
	autoscalingv1 "agones.dev/agones/pkg/apis/autoscaling/v1"
	"github.com/Nerzal/gocloak/v13"
	gamebackend "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	corev1 "k8s.io/api/core/v1"
	k8errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

	if s.server.GlobalConfig.GameBackend.Mode != config.LocalMode {
		err = s.setupNewDimension(ctx, dimension)
	} else {
		log.Logger.WithContext(ctx).Infof("Running Local Mode - Skipping dimension setup")
	}

	return dimension.ToPb(), err
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

	dimension, err := s.server.GamebackendService.FindDimension(ctx, request)
	if err != nil {
		return nil, err
	}
	if dimension == nil {
		return nil, model.ErrDoesNotExist.Err()
	}

	err = s.server.GamebackendService.DeleteDimensionById(ctx, dimension.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "deleteing dimension %s: %s", dimension.Name, err.Error())
	}
	log.Logger.WithContext(ctx).Infof("Deleted dimension %s (%s)", dimension.Name, dimension.Id.String())

	if s.server.GlobalConfig.GameBackend.Mode != config.LocalMode {
		for _, m := range dimension.Maps {
			log.Logger.Infof("Deleting gameserver dimension %s map %s", dimension.Name, m.Name)
			err = s.deleteGameServers(ctx, dimension, m)
			if err != nil {
				err = errors.Join(err, status.Errorf(codes.Internal, "deleting game servers for dimension %s, map %s: %s", dimension.Name, m.Name, err.Error()))
			}
		}
	} else {
		log.Logger.WithContext(ctx).Info("Local Mode: Not deleting game servers")
	}

	return &emptypb.Empty{}, err
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
		err = joinErrorAndLog(nil, "deleteing map %v: %w", request, err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	log.Logger.WithContext(ctx).Infof("Deleted map %s (%s)", m.Name, m.Id.String())

	var errs error
	if s.server.GlobalConfig.GameBackend.Mode != config.LocalMode {
		dimensions, err := s.server.GamebackendService.FindDimensionsWithMapIds(ctx, []*uuid.UUID{m.Id})
		for _, dimension := range dimensions {
			log.Logger.Infof("Deleting gameserver dimension %s map %s", dimension.Name, m.Name)
			err = s.deleteGameServers(ctx, dimension, m)
			if err != nil {
				errs = joinErrorAndLog(errs, "deleting game servers for dimension %s, map %s: %s", dimension.Name, m.Name, err.Error())
			}
		}
	} else {
		log.Logger.WithContext(ctx).Info("Local Mode: Not deleting game servers")
	}

	return &emptypb.Empty{}, errs
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

	newDimension, err := s.server.GamebackendService.DuplicateDimension(ctx, request.Target, request.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if newDimension == nil {
		return nil, status.Error(codes.Internal, "failed to create new dimension")
	}

	if s.server.GlobalConfig.GameBackend.Mode != config.LocalMode {
		err = s.setupNewDimension(ctx, newDimension)
		if err != nil {
			err = status.Errorf(codes.Internal, "setup dimension: %s", err.Error())
		}
	} else {
		log.Logger.WithContext(ctx).Info("Local Mode: Not creating game servers")
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

	originalDimension, err := s.server.GamebackendService.FindDimension(ctx, request.Target)
	if err != nil {
		return nil, err
	}
	if originalDimension == nil {
		return nil, model.ErrDoesNotExist.Err()
	}

	newDimension, err := s.server.GamebackendService.EditDimension(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var errs error
	if s.server.GlobalConfig.GameBackend.Mode != config.LocalMode {
		if originalDimension.Name != newDimension.Name {
			// Name change requires delete and recreate
			for _, m := range originalDimension.Maps {
				err := s.deleteGameServers(ctx, originalDimension, m)
				if err != nil {
					errs = joinErrorAndLog(errs, "deleting old game server for dimension %s, Map %s failed: %w",
						originalDimension.Name, m.Name, err)
				}
			}

			err = s.setupNewDimension(ctx, newDimension)
			if err != nil {
				errs = joinErrorAndLog(errs, "setting up dimension %s gameservers: %w", newDimension.Name, err)
			}
		} else {
			if request.EditMaps {
				currentMaps := make(map[*uuid.UUID]*model.Map, len(request.MapIds))
				for _, m := range newDimension.Maps {
					currentMaps[m.Id] = m
				}

				for _, m := range originalDimension.Maps {
					if _, ok := currentMaps[m.Id]; !ok {
						err := s.deleteGameServers(ctx, newDimension, m)
						if err != nil {
							errs = joinErrorAndLog(errs, "unable to delete old gameserver world %s: %w", m.Name, err)
						}
					} else {
						delete(currentMaps, m.Id)

						if request.OptionalVersion != nil {
							err := s.updateGameServers(ctx, newDimension, m)
							if err != nil {
								errs = joinErrorAndLog(errs, "unable to update gameserver world %s: %w", m.Name, err)
							}
						}
					}
				}

				// newMaps now only contains map that weren't in the original
				for _, newMap := range currentMaps {
					err = s.createGameServers(ctx, newDimension, newMap)
					if err != nil {
						errs = joinErrorAndLog(errs, "unable to delete old gameserver world %s: %w", newMap.Name, err)
					}
				}
			} else if request.OptionalVersion != nil {
				// Maps weren't changed, so update the versions
				for _, m := range newDimension.Maps {
					err := s.updateGameServers(ctx, newDimension, m)
					if err != nil {
						errs = joinErrorAndLog(errs, "unable to update gameserver world %s: %w", m.Name, err)
					}
				}
			}
		}
	} else {
		log.Logger.WithContext(ctx).Info("Local Mode: Not creating game servers")
	}

	return newDimension.ToPb(), errs
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

	originalMap, err := s.server.GamebackendService.FindMap(ctx, request.Target)
	if err != nil {
		return nil, err
	}
	if originalMap == nil {
		return nil, model.ErrDoesNotExist.Err()
	}

	newMap, err := s.server.GamebackendService.EditMap(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var errs error
	if s.server.GlobalConfig.GameBackend.Mode != config.LocalMode {
		if originalMap.Name != newMap.Name {
			// Need to delete and recreate
			for _, dimension := range originalMap.Dimensions {
				err = s.deleteGameServers(ctx, dimension, originalMap)
				if err != nil {
					errs = joinErrorAndLog(errs, "deleting old game server for dimension %s, Map %s failed: %w",
						dimension.Name, originalMap.Name, err)
				}

				err = s.createGameServers(ctx, dimension, newMap)
				if err != nil {
					errs = joinErrorAndLog(errs, "setting up dimension %s gameservers: %w", dimension.Name, err)
				}
			}
		} else {
			// Update dimension
			for _, dimension := range newMap.Dimensions {
				err = s.updateGameServers(ctx, dimension, newMap)
				if err != nil {
					errs = joinErrorAndLog(errs, "unable to update gameserver world %s: %w", newMap.Name, err)
				}
			}
		}
	} else {
		log.Logger.WithContext(ctx).Info("Local Mode: Not creating game servers")
	}

	return newMap.ToPb(), errs
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

	return out.ToPb(), nil
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

	dimension, err := s.server.GamebackendService.FindDimension(ctx, request)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if dimension == nil {
		return nil, model.ErrDoesNotExist.Err()
	}

	return dimension.ToPb(), nil
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

	m, err := s.server.GamebackendService.FindMap(ctx, request)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if m == nil {
		return nil, model.ErrDoesNotExist.Err()
	}

	return m.ToPb(), nil
}

func NewServerManagerServiceServer(
	ctx context.Context,
	server *gamebackend.GameBackendServerContext,
) (pb.ServerManagerServiceServer, error) {
	token, err := server.KeycloakClient.LoginClient(
		ctx,
		server.GlobalConfig.GameBackend.Keycloak.ClientId,
		server.GlobalConfig.GameBackend.Keycloak.ClientSecret,
		server.GlobalConfig.Keycloak.Realm,
	)
	if err != nil {
		return nil, fmt.Errorf("login keycloak: %v", err)
	}

	err = createRoles(ctx,
		server.KeycloakClient,
		token.AccessToken,
		server.GlobalConfig.Keycloak.Realm,
		server.GlobalConfig.GameBackend.Keycloak.Id,
		&ConnectionRoles,
	)
	if err != nil {
		return nil, err
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
		s.server.GlobalConfig.Keycloak.Realm,
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
	_, claims, err := helpers.VerifyClaims(ctx, s.server.KeycloakClient, s.server.GlobalConfig.Keycloak.Realm)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("verify claims: %v", err)
		return model.ErrUnauthorized.Err()
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleServerManager, model.GamebackendClientId) {
		return model.ErrUnauthorized.Err()
	}

	return nil
}

func (s serverManagerServiceServer) createGameServers(
	ctx context.Context,
	dimension *model.Dimension,
	m *model.Map,
) error {
	if s.server.AgonesClient == nil {
		return ErrNoAgonesConnect
	}

	// Create the fleet
	fleet, err := s.server.AgonesClient.AgonesV1().Fleets(s.server.GlobalConfig.Agones.Namespace).Create(
		ctx,
		buildFleet(dimension, m, s.server.GlobalConfig.Agones.Namespace),
		metav1.CreateOptions{},
	)
	if err != nil {
		return fmt.Errorf("creating fleet: %w", err)
	}

	// Create autoscaler
	_, err = s.server.AgonesClient.AutoscalingV1().FleetAutoscalers(fleet.Namespace).Create(
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
	if s.server.AgonesClient == nil {
		return ErrNoAgonesConnect
	}

	namespace := s.server.GlobalConfig.Agones.Namespace

	// Delete autoscaler
	err := s.server.AgonesClient.AutoscalingV1().FleetAutoscalers(namespace).Delete(
		ctx,
		getFleetAutoscalerName(dimension, m),
		metav1.DeleteOptions{},
	)

	if err != nil &&
		!k8errors.IsNotFound(err) {
		return fmt.Errorf("deleting fleet autoscaler: %w", err)
	}

	// Delete autoscaler
	err = s.server.AgonesClient.AgonesV1().Fleets(namespace).Delete(
		ctx,
		getFleetName(dimension, m),
		metav1.DeleteOptions{},
	)
	if err != nil &&
		!k8errors.IsNotFound(err) {
		return fmt.Errorf("deleting fleet: %w", err)
	}

	return nil
}

func (s serverManagerServiceServer) updateGameServers(
	ctx context.Context,
	dimension *model.Dimension,
	m *model.Map,
) error {
	if s.server.AgonesClient == nil {
		return ErrNoAgonesConnect
	}

	fleet, err := s.server.AgonesClient.AgonesV1().Fleets(s.server.GlobalConfig.Agones.Namespace).Update(
		ctx,
		buildFleet(dimension, m, s.server.GlobalConfig.Agones.Namespace),
		metav1.UpdateOptions{},
	)
	if err != nil {
		return fmt.Errorf("creating fleet: %w", err)
	}

	_, err = s.server.AgonesClient.AutoscalingV1().FleetAutoscalers(fleet.Namespace).Update(
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
	return &v1.Fleet{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getFleetName(dimension, m),
			Namespace: namespace,
		},
		Spec: v1.FleetSpec{
			Replicas:   1,
			Scheduling: "",
			Strategy: appsv1.DeploymentStrategy{
				Type: "RollingUpdate",
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 25,
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 25,
					},
				},
			},
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
									Name:  "gameserver",
									Image: dimension.GetImageName(),
									Args: []string{
										m.Path,
										"-log",
									},
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

	var errs error
	for _, m := range dimension.Maps {
		err := s.createGameServers(ctx, dimension, m)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
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
		log.Logger.WithContext(ctx).Errorf("target type unknown: %v", requestTarget)
		return nil, model.ErrHandleRequest.Err()
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if out == nil {
		return nil, model.ErrDoesNotExist.Err()
	}

	return out, nil
}

func joinErrorAndLog(err error, str string, args ...any) error {
	newErr := fmt.Errorf(str, args...)
	log.Logger.Error(err)
	return errors.Join(err, newErr)
}

package srv

import (
	context "context"
	"fmt"
	"reflect"

	"agones.dev/agones/pkg/client/clientset/versioned"
	"github.com/Nerzal/gocloak/v13"
	gamebackend "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"k8s.io/client-go/rest"
)

var (
	ServerManagerRoles = make([]*gocloak.Role, 0)

	RoleServerManager = registerServerManagerRole(&gocloak.Role{
		Name:        gocloak.StringP("server_manager"),
		Description: gocloak.StringP("Allow the use of the server manager service"),
	})
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

// CreateChatTemplate implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) CreateChatTemplate(
	ctx context.Context,
	request *pb.CreateChatTemplateRequest,
) (*pb.ChatTemplate, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	chatTemplate, err := s.server.GamebackendService.CreateChatTemplate(ctx, request.Name)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "creating: %w", err)
	}

	return chatTemplate.ToPb(), nil
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

	chatTemplateIds, err := helpers.ParseUUIDs(request.ChatTemplateIds)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dimension, err := s.server.GamebackendService.CreateDimension(
		ctx,
		request.Name,
		request.Location,
		request.Version,
		mapIds,
		chatTemplateIds,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "creating: %w", err)
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
		return nil, status.Errorf(codes.Internal, "creating: %w", err)
	}

	return m.ToPb(), nil
}

// DeleteChatTemplate implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) DeleteChatTemplate(
	ctx context.Context,
	request *pb.ChatTemplateTarget,
) (*emptypb.Empty, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	switch target := request.FindBy.(type) {
	case *pb.ChatTemplateTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, err
		}
		err = s.server.GamebackendService.DeleteChatTemplateById(ctx, &id)
		return &emptypb.Empty{}, err

	case *pb.ChatTemplateTarget_Name:
		err = s.server.GamebackendService.DeleteChatTemplateByName(ctx, target.Name)
		return &emptypb.Empty{}, err

	default:
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}
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

	switch target := request.FindBy.(type) {
	case *pb.DimensionTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, err
		}
		err = s.server.GamebackendService.DeleteDimensionById(ctx, &id)
		return &emptypb.Empty{}, err

	case *pb.DimensionTarget_Name:
		err = s.server.GamebackendService.DeleteDimensionByName(ctx, target.Name)
		return &emptypb.Empty{}, err

	default:
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}
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

	switch target := request.FindBy.(type) {
	case *pb.MapTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, err
		}
		err = s.server.GamebackendService.DeleteMapById(ctx, &id)
		return &emptypb.Empty{}, err

	case *pb.MapTarget_Name:
		err = s.server.GamebackendService.DeleteMapByName(ctx, target.Name)
		return &emptypb.Empty{}, err

	default:
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}
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
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if newDimension == nil {
		return nil, status.Error(codes.Internal, "failed to create new dimension")
	}

	return newDimension.ToPb(), err
}

// EditChatTemplate implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) EditChatTemplate(
	ctx context.Context,
	request *pb.EditChatTemplateRequest,
) (*pb.ChatTemplate, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	out, err := s.server.GamebackendService.EditChatTemplate(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return out.ToPb(), nil
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

	out, err := s.server.GamebackendService.EditDimension(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return out.ToPb(), nil
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

	out, err := s.server.GamebackendService.EditMap(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return out.ToPb(), nil
}

// GetAllChatTemplates implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetAllChatTemplates(
	ctx context.Context,
	request *emptypb.Empty,
) (*pb.ChatTemplates, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	out, err := s.server.GamebackendService.FindAllChatTemplates(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ChatTemplates{ChatTemplates: out.ToPb()}, nil
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

// GetChatTemplate implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetChatTemplate(
	ctx context.Context,
	request *pb.ChatTemplateTarget,
) (*pb.ChatTemplate, error) {
	err := s.hasServerManagerRole(ctx)
	if err != nil {
		return nil, err
	}

	var out *model.ChatTemplate
	switch target := request.FindBy.(type) {
	case *pb.ChatTemplateTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, err
		}
		out, err = s.server.GamebackendService.FindChatTemplateById(ctx, &id)

	case *pb.ChatTemplateTarget_Name:
		out, err = s.server.GamebackendService.FindChatTemplateByName(ctx, target.Name)

	default:
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return out.ToPb(), nil
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
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
			log.WithContext(ctx).Errorf("creating config: %v", err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		agones, err := versioned.NewForConfig(conf)
		if err != nil {
			log.WithContext(ctx).Errorf("creating agones connection: %v", err)
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
	claims, err := helpers.ExtractClaims(ctx)
	if err != nil {
		log.WithContext(ctx).Infof("extract claims failed: %v", err)
		return model.ErrUnauthorized
	}

	// Validate requester has correct permission
	if !claims.HasResourceRole(RoleServerManager, model.GamebackendClientId) {
		return model.ErrUnauthorized
	}

	return nil
}

package srv

import (
	context "context"
	"fmt"

	"agones.dev/agones/pkg/client/clientset/versioned"
	"github.com/Nerzal/gocloak/v13"
	gamebackend "github.com/ShatteredRealms/go-backend/cmd/gamebackend/app"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
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
	panic("unimplemented")
}

// CreateDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) CreateDimension(
	ctx context.Context,
	request *pb.CreateDimensionRequest,
) (*pb.Dimension, error) {
	panic("unimplemented")
}

// CreateMap implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) CreateMap(
	ctx context.Context,
	request *pb.CreateMapRequest,
) (*pb.Map, error) {
	panic("unimplemented")
}

// DeleteChatTemplate implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) DeleteChatTemplate(
	ctx context.Context,
	request *pb.ChatTemplateTarget,
) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// DeleteDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) DeleteDimension(
	ctx context.Context,
	request *pb.DimensionTarget,
) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// DeleteMap implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) DeleteMap(
	ctx context.Context,
	request *pb.MapTarget,
) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// DuplicateDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) DuplicateDimension(
	ctx context.Context,
	request *pb.DuplicateDimensionRequest,
) (*pb.Dimension, error) {
	panic("unimplemented")
}

// EditChatTemplate implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) EditChatTemplate(
	ctx context.Context,
	request *pb.EditChatTemplateRequest,
) (*pb.ChatTemplate, error) {
	panic("unimplemented")
}

// EditDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) EditDimension(
	ctx context.Context,
	request *pb.EditDimensionRequest,
) (*pb.Dimension, error) {
	panic("unimplemented")
}

// EditMap implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) EditMap(
	ctx context.Context,
	request *pb.EditMapRequest,
) (*pb.Map, error) {
	panic("unimplemented")
}

// GetAllChatTemplates implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetAllChatTemplates(
	ctx context.Context,
	request *emptypb.Empty,
) (*pb.ChatTemplates, error) {
	panic("unimplemented")
}

// GetAllDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetAllDimension(
	ctx context.Context,
	request *emptypb.Empty,
) (*pb.Dimensions, error) {
	panic("unimplemented")
}

// GetAllMaps implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetAllMaps(
	ctx context.Context,
	request *emptypb.Empty,
) (*pb.Maps, error) {
	panic("unimplemented")
}

// GetChatTemplate implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetChatTemplate(
	ctx context.Context,
	request *pb.ChatTemplateTarget,
) (*pb.ChatTemplate, error) {
	panic("unimplemented")
}

// GetDimension implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetDimension(
	ctx context.Context,
	request *pb.DimensionTarget,
) (*pb.Dimension, error) {
	panic("unimplemented")
}

// GetMap implements pb.ServerManagerServiceServer.
func (s *serverManagerServiceServer) GetMap(
	ctx context.Context,
	request *pb.MapTarget,
) (*pb.Map, error) {
	panic("unimplemented")
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

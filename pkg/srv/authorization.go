package srv

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type AuthorizationServiceServer struct {
	pb.UnimplementedAuthorizationServiceServer
	UserService       service.UserService
	PermissionService service.PermissionService
	roleService       service.RoleService
	AllPermissions    *pb.UserPermissions
	authInterceptor   *interceptor.AuthInterceptor
	userUpdates       chan string
	roleUpdates       chan string
}

func NewAuthorizationServiceServer(
	u service.UserService,
	permissionService service.PermissionService,
	roleService service.RoleService,
) *AuthorizationServiceServer {
	return &AuthorizationServiceServer{
		UserService:       u,
		PermissionService: permissionService,
		roleService:       roleService,
		userUpdates:       make(chan string),
		roleUpdates:       make(chan string),
	}
}

func (s *AuthorizationServiceServer) AddAuthInterceptor(interceptor *interceptor.AuthInterceptor) {
	s.authInterceptor = interceptor
}

func (s *AuthorizationServiceServer) GetAuthorization(
	ctx context.Context,
	message *pb.Username,
) (*pb.AuthorizationMessage, error) {
	user := s.UserService.FindByUsername(ctx, message.Username)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	permissions := s.PermissionService.FindPermissionsForUsername(ctx, user.Username).ToPB().Permissions
	roles := user.Roles.ToPB().Roles
	for i, role := range roles {
		roles[i].Permissions = s.PermissionService.FindPermissionsForRole(ctx, role.Name).ToPB().Permissions
	}

	resp := &pb.AuthorizationMessage{
		Username:    message.Username,
		Roles:       roles,
		Permissions: permissions,
	}

	return resp, nil
}

func (s *AuthorizationServiceServer) AddAuthorization(ctx context.Context, message *pb.AuthorizationMessage) (*emptypb.Empty, error) {
	user := s.UserService.FindByUsername(ctx, message.Username)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	for _, v := range message.Permissions {
		err := s.PermissionService.AddPermissionForUser(
			ctx,
			&model.UserPermission{
				Username:   user.Username,
				Permission: v.Permission.Value,
				Other:      v.Other,
			},
		)

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	for _, v := range message.Roles {
		err := s.UserService.AddToRole(
			ctx,
			user,
			&model.Role{
				Name: v.Name,
			})

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	s.userUpdates <- message.Username
	err := s.authInterceptor.ClearUserCache(message.Username)
	return &emptypb.Empty{}, err
}

func (s *AuthorizationServiceServer) RemoveAuthorization(ctx context.Context, message *pb.AuthorizationMessage) (*emptypb.Empty, error) {
	user := s.UserService.FindByUsername(ctx, message.Username)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	for _, v := range message.Permissions {
		err := s.PermissionService.RemPermissionForUser(
			ctx,
			&model.UserPermission{
				Username:   user.Username,
				Permission: v.Permission.Value,
				Other:      v.Other,
			},
		)

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	for _, v := range message.Roles {
		err := s.UserService.RemFromRole(
			ctx,
			user,
			&model.Role{
				Name: v.Name,
			})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	s.userUpdates <- message.Username
	err := s.authInterceptor.ClearUserCache(message.Username)
	return &emptypb.Empty{}, err
}

func (s *AuthorizationServiceServer) GetRoles(ctx context.Context, message *emptypb.Empty) (*pb.UserRoles, error) {
	resp := &pb.UserRoles{
		Roles: s.roleService.FindAll(ctx).ToPB().Roles,
	}
	for i, v := range resp.Roles {
		resp.Roles[i].Permissions = s.PermissionService.FindPermissionsForRole(ctx, v.Name).ToPB().Permissions
	}

	return resp, nil
}

func (s *AuthorizationServiceServer) GetRole(ctx context.Context, message *pb.RoleName) (*pb.UserRole, error) {
	resp := s.roleService.FindByName(ctx, message.Name).ToPB()
	resp.Permissions = s.PermissionService.FindPermissionsForRole(ctx, message.Name).ToPB().Permissions

	return resp, nil
}

func (s *AuthorizationServiceServer) CreateRole(ctx context.Context, message *pb.UserRole) (*emptypb.Empty, error) {
	_, err := s.roleService.Create(
		ctx,
		&model.Role{
			Name: message.Name,
		})

	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthorizationServiceServer) EditRole(ctx context.Context, message *pb.RequestEditUserRole) (*emptypb.Empty, error) {

	if message.Permissions != nil {
		newPermissions := make([]*model.RolePermission, len(message.Permissions))
		for i, permission := range message.Permissions {
			newPermissions[i] = &model.RolePermission{
				Role:       message.Name,
				Permission: permission.Permission.Value,
				Other:      permission.Other,
			}
		}

		err := s.PermissionService.ResetPermissionsForRole(ctx, message.Name, newPermissions)
		if err != nil {
			return nil, err
		}
	}

	if message.NewName != nil {
		err := s.roleService.NewName(
			ctx,
			message.Name,
			message.NewName.Value,
		)
		if err != nil {
			return nil, err
		}

		err = s.PermissionService.UpdateRoleName(ctx, message.Name, message.NewName.Value)
		if err != nil {
			// Update permissions failed, rollback name change
			// @TODO: Use trx with commits
			innerErr := s.roleService.NewName(
				ctx,
				message.NewName.Value,
				message.Name,
			)
			if innerErr != nil {
				return nil, innerErr
			}

			return nil, err
		}
	}

	s.roleUpdates <- message.Name
	err := s.authInterceptor.ClearCache()
	return &emptypb.Empty{}, err
}

func (s *AuthorizationServiceServer) DeleteRole(ctx context.Context, message *pb.UserRole) (*emptypb.Empty, error) {
	err := s.roleService.Delete(
		ctx,
		&model.Role{
			Name: message.Name,
		},
	)

	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	s.roleUpdates <- message.Name
	err = s.authInterceptor.ClearCache()
	return &emptypb.Empty{}, err
}

func (s *AuthorizationServiceServer) GetAllPermissions(ctx context.Context, message *emptypb.Empty) (*pb.UserPermissions, error) {
	return s.AllPermissions, nil
}

func (s *AuthorizationServiceServer) SetupAllPermissions(accountsServicesInfo map[string]grpc.ServiceInfo) {
	accountsPermissions := setupPermissions(accountsServicesInfo)
	charactersPermissions := setupPermissions(getCharactersServiceInfo())
	gameBackendPermissions := setupPermissions(getGameBackendServiceInfo())
	chatPermissions := setupPermissions(getChatServiceInfo())

	s.AllPermissions = &pb.UserPermissions{Permissions: accountsPermissions}
	s.AllPermissions.Permissions = append(s.AllPermissions.Permissions, charactersPermissions...)
	s.AllPermissions.Permissions = append(s.AllPermissions.Permissions, gameBackendPermissions...)
	s.AllPermissions.Permissions = append(s.AllPermissions.Permissions, chatPermissions...)
}

func (s *AuthorizationServiceServer) SubscribeUserUpdates(message *emptypb.Empty, stream pb.AuthorizationService_SubscribeUserUpdatesServer) error {
	for {
		select {
		case <-stream.Context().Done():
			log.Debug("User subscribe context closed")
			return nil
		case username := <-s.userUpdates:
			log.Debug("Sending update")
			err := stream.Send(&pb.Username{Username: username})
			log.Debug("Broadcast role update")
			if err != nil {
				return err
			}
		}
	}
}
func (s *AuthorizationServiceServer) SubscribeRoleUpdates(message *emptypb.Empty, stream pb.AuthorizationService_SubscribeRoleUpdatesServer) error {
	for {
		select {
		case <-stream.Context().Done():
			log.Debug("Role subscribe context closed")
			return nil
		case roleName := <-s.roleUpdates:
			log.Debug("Sending update")
			err := stream.Send(&pb.RoleName{Name: roleName})
			log.Debug("Broadcast role update")
			if err != nil {
				return err
			}
		}
	}
}

func getChatServiceInfo() map[string]grpc.ServiceInfo {
	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, pb.UnimplementedChatServiceServer{})
	return grpcServer.GetServiceInfo()
}

func getCharactersServiceInfo() map[string]grpc.ServiceInfo {
	grpcServer := grpc.NewServer()
	pb.RegisterCharactersServiceServer(grpcServer, pb.UnimplementedCharactersServiceServer{})
	return grpcServer.GetServiceInfo()
}

func getGameBackendServiceInfo() map[string]grpc.ServiceInfo {
	grpcServer := grpc.NewServer()
	pb.RegisterConnectionServiceServer(grpcServer, pb.UnimplementedConnectionServiceServer{})
	return grpcServer.GetServiceInfo()
}

func setupPermissions(serviceInfos map[string]grpc.ServiceInfo) []*pb.UserPermission {

	count := 0
	for _, serviceInfo := range serviceInfos {
		count += len(serviceInfo.Methods)
	}

	methods := make([]*pb.UserPermission, count)
	index := 0
	for serviceName, serviceInfo := range serviceInfos {
		for _, method := range serviceInfo.Methods {
			methods[index] = &pb.UserPermission{
				Permission: &wrapperspb.StringValue{Value: fmt.Sprintf("/%s/%s", serviceName, method.Name)},
			}
			index++
		}
	}

	return methods
}

package srv

import (
	"context"
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

type AuthorizationServiceServer struct {
	pb.UnimplementedAuthorizationServiceServer
	UserService       service.UserService
	PermissionService service.PermissionService
	roleService       service.RoleService
	AllPermissions    *pb.UserPermissions
	authInterceptor   *interceptor.AuthInterceptor
	userUpdates       chan uint64
	roleUpdates       chan uint64
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
		userUpdates:       make(chan uint64),
		roleUpdates:       make(chan uint64),
	}
}

func (s *AuthorizationServiceServer) AddAuthInterceptor(interceptor *interceptor.AuthInterceptor) {
	s.authInterceptor = interceptor
}

func (s *AuthorizationServiceServer) GetAuthorization(
	ctx context.Context,
	message *pb.IDMessage,
) (*pb.AuthorizationMessage, error) {
	user := s.UserService.FindById(ctx, uint(message.Id))
	if user == nil || !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	permissions := s.PermissionService.FindPermissionsForUserID(ctx, user.ID).ToPB().Permissions
	roles := user.Roles.ToPB().Roles
	for i, role := range roles {
		roles[i].Permissions = s.PermissionService.FindPermissionsForRoleID(ctx, uint(role.Id)).ToPB().Permissions
	}

	resp := &pb.AuthorizationMessage{
		UserId:      message.Id,
		Roles:       roles,
		Permissions: permissions,
	}

	return resp, nil
}

func (s *AuthorizationServiceServer) AddAuthorization(ctx context.Context, message *pb.AuthorizationMessage) (*emptypb.Empty, error) {
	user := s.UserService.FindById(ctx, uint(message.UserId))
	if user == nil || !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	for _, v := range message.Permissions {
		err := s.PermissionService.AddPermissionForUser(
			ctx,
			&model.UserPermission{
				UserID:     user.ID,
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
				Model: gorm.Model{
					ID: uint(v.Id),
				},
			})

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	s.userUpdates <- message.UserId
	err := s.authInterceptor.ClearUserCache(uint(message.UserId))
	return &emptypb.Empty{}, err
}

func (s *AuthorizationServiceServer) RemoveAuthorization(ctx context.Context, message *pb.AuthorizationMessage) (*emptypb.Empty, error) {
	user := s.UserService.FindById(ctx, uint(message.UserId))
	if user == nil || !user.Exists() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	for _, v := range message.Permissions {
		err := s.PermissionService.RemPermissionForUser(
			ctx,
			&model.UserPermission{
				UserID:     user.ID,
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
				Model: gorm.Model{
					ID: uint(v.Id),
				},
			})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	s.userUpdates <- message.UserId
	err := s.authInterceptor.ClearUserCache(uint(message.UserId))
	return &emptypb.Empty{}, err
}

func (s *AuthorizationServiceServer) GetRoles(ctx context.Context, message *emptypb.Empty) (*pb.UserRoles, error) {
	resp := &pb.UserRoles{
		Roles: s.roleService.FindAll(ctx).ToPB().Roles,
	}
	for i, v := range resp.Roles {
		resp.Roles[i].Permissions = s.PermissionService.FindPermissionsForRoleID(ctx, uint(v.Id)).ToPB().Permissions
	}

	return resp, nil
}

func (s *AuthorizationServiceServer) GetRole(ctx context.Context, message *pb.IDMessage) (*pb.UserRole, error) {
	resp := s.roleService.FindById(ctx, uint(message.Id)).ToPB()
	resp.Permissions = s.PermissionService.FindPermissionsForRoleID(ctx, uint(message.Id)).ToPB().Permissions

	return resp, nil
}

func (s *AuthorizationServiceServer) CreateRole(ctx context.Context, message *pb.UserRole) (*emptypb.Empty, error) {
	_, err := s.roleService.Create(
		ctx,
		&model.Role{
			Name: message.Name.Value,
		})

	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthorizationServiceServer) EditRole(ctx context.Context, message *pb.UserRole) (*emptypb.Empty, error) {
	if message.Name != nil {
		err := s.roleService.Update(
			ctx,
			&model.Role{
				Model: gorm.Model{
					ID: uint(message.Id),
				},
				Name: message.Name.Value,
			})
		if err != nil {
			return nil, err
		}
	}
	if message.Permissions != nil {
		newPermissions := make([]*model.RolePermission, len(message.Permissions))
		for i, permission := range message.Permissions {
			newPermissions[i] = &model.RolePermission{
				RoleID:     uint(message.Id),
				Permission: permission.Permission.Value,
				Other:      permission.Other,
			}
		}

		err := s.PermissionService.ResetPermissionsForRole(ctx, uint(message.Id), newPermissions)
		if err != nil {
			return nil, err
		}
	}

	s.roleUpdates <- message.Id
	err := s.authInterceptor.ClearCache()
	return &emptypb.Empty{}, err
}

func (s *AuthorizationServiceServer) DeleteRole(ctx context.Context, message *pb.UserRole) (*emptypb.Empty, error) {
	err := s.roleService.Delete(
		ctx,
		&model.Role{
			Model: gorm.Model{
				ID: uint(message.Id),
			},
		},
	)

	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	s.roleUpdates <- message.Id
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
		case userId := <-s.userUpdates:
			log.Debug("Sending update")
			err := stream.Send(&pb.IDMessage{Id: userId})
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
		case userId := <-s.roleUpdates:
			log.Debug("Sending update")
			err := stream.Send(&pb.IDMessage{Id: userId})
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

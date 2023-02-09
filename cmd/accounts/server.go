package main

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	accountService "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/srv"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(
	u accountService.UserService,
	p accountService.PermissionService,
	r accountService.RoleService,
	jwt service.JWTService,
	ctx context.Context,
) (*grpc.Server, *runtime.ServeMux, error) {
	publicRPCs := map[string]struct{}{
		"/sro.accounts.HealthService/Health":           {},
		"/sro.accounts.AuthenticationService/Login":    {},
		"/sro.accounts.AuthenticationService/Register": {},
	}

	authorizationServiceServer := srv.NewAuthorizationServiceServer(u, p, r)
	authInterceptor := interceptor.NewAuthInterceptor(jwt, publicRPCs, getPermissions(ctx, authorizationServiceServer))
	authorizationServiceServer.AddAuthInterceptor(authInterceptor)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInterceptor.Unary(), helpers.UnaryLogRequest()),
		grpc.ChainStreamInterceptor(authInterceptor.Stream(), helpers.StreamLogRequest()),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	authenticationServiceServer := srv.NewAuthenticationServiceServer(u, jwt, p)
	pb.RegisterAuthenticationServiceServer(grpcServer, authenticationServiceServer)
	err := pb.RegisterAuthenticationServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Accounts.Local.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	userServiceServer := srv.NewUserServiceServer(u, p, jwt)
	pb.RegisterUserServiceServer(grpcServer, userServiceServer)
	err = pb.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Accounts.Local.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	healthServiceServer := srv.NewHealthServiceServer()
	pb.RegisterHealthServiceServer(grpcServer, healthServiceServer)
	err = pb.RegisterHealthServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Accounts.Local.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	pb.RegisterAuthorizationServiceServer(grpcServer, authorizationServiceServer)
	err = pb.RegisterAuthorizationServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		conf.Accounts.Local.Address(),
		opts,
	)
	if err != nil {
		return nil, nil, err
	}

	// Compute the AllPermissions method once and save in memory
	authorizationServiceServer.SetupAllPermissions(grpcServer.GetServiceInfo())

	return grpcServer, gwmux, nil
}

func getPermissions(
	ctx context.Context,
	server *srv.AuthorizationServiceServer,
) func(userID uint) map[string]bool {
	return func(userID uint) map[string]bool {
		// UserID 0 is for server communication
		if userID == 0 {
			allPerms := make(map[string]bool, len(server.AllPermissions.Permissions))
			for _, perm := range server.AllPermissions.Permissions {
				allPerms[perm.Permission.Value] = true
			}

			return allPerms
		}

		user := server.UserService.FindById(ctx, userID)
		if user == nil || !user.Exists() {
			return map[string]bool{}
		}

		resp := make(map[string]bool)

		for _, role := range user.Roles {
			for _, rolePermission := range server.PermissionService.FindPermissionsForRoleID(ctx, role.ID) {
				resp[rolePermission.Permission] = resp[rolePermission.Permission] || rolePermission.Other
			}
		}

		for _, userPermission := range server.PermissionService.FindPermissionsForUserID(ctx, userID) {
			resp[userPermission.Permission] = resp[userPermission.Permission] || userPermission.Other
		}

		return resp
	}
}

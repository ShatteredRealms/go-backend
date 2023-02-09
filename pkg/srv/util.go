package srv

import (
	"context"
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"time"
)

const (
	retryAfter = time.Second * 10
)

func GetPermissions(
	ctx context.Context,
	authorizationService pb.AuthorizationServiceClient,
	jwtService service.JWTService,
	requestingHost string,
) func(userID uint) map[string]bool {
	return func(userID uint) map[string]bool {
		md := metadata.New(
			map[string]string{
				"authorization": fmt.Sprintf(
					"Bearer %s", generateTemporaryServerToken(ctx, jwtService, requestingHost),
				),
			},
		)
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		authorizations, err := authorizationService.GetAuthorization(ctx, &pb.IDMessage{Id: uint64(userID)})

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return map[string]bool{}
		}

		resp := make(map[string]bool)

		for _, role := range authorizations.Roles {
			for _, rolePermission := range role.Permissions {
				resp[rolePermission.Permission.Value] = resp[rolePermission.Permission.Value] || rolePermission.Other
			}
		}

		for _, userPermission := range authorizations.Permissions {
			resp[userPermission.Permission.Value] = resp[userPermission.Permission.Value] || userPermission.Other
		}

		return resp
	}
}

func generateTemporaryServerToken(ctx context.Context, jwtService service.JWTService, requestingHost string) string {
	out, _ := jwtService.Create(ctx, time.Second, requestingHost, jwt.MapClaims{"sub": 0})
	return out
}
func ProcessUserUpdates(
	ctx context.Context,
	authorizationClient pb.AuthorizationServiceClient,
	interceptor *interceptor.AuthInterceptor,
	jwtService service.JWTService,
	serviceAuthName string,
) {
	userUpdatesClient, err := authorizationClient.SubscribeUserUpdates(serverAuthContext(ctx, jwtService, serviceAuthName), &emptypb.Empty{})
	if err != nil {
		log.WithContext(ctx).Errorf("Unable to subscribe to user updates from authorization client. Retrying in %d seconds", retryAfter/time.Second)
		time.Sleep(retryAfter)
		ProcessUserUpdates(ctx, authorizationClient, interceptor, jwtService, serviceAuthName)
		return
	}
	log.Info("Successfully subscribed to user updates from authorization server.")
	for {
		msg, err := userUpdatesClient.Recv()

		if err == nil {
			log.Debugf("Update to user %d permissions. Clearing permissions cache for that user.", msg.Id)
			err = interceptor.ClearUserCache(uint(msg.Id))
			if err != nil {
				log.WithContext(ctx).Warning("Clearing cache: %v", err)
			}
		} else if err == io.EOF {
			log.Infof("User updates stream ended. Retrying in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			ProcessUserUpdates(ctx, authorizationClient, interceptor, jwtService, serviceAuthName)
			return
		} else {
			log.WithContext(ctx).Errorf("User updates: %v.", err)
			log.Infof("Retrying connection in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			ProcessUserUpdates(ctx, authorizationClient, interceptor, jwtService, serviceAuthName)
			return
		}
	}
}

func ProcessRoleUpdates(
	ctx context.Context,
	authorizationClient pb.AuthorizationServiceClient,
	interceptor *interceptor.AuthInterceptor,
	jwtService service.JWTService,
	serviceAuthName string,
) {
	roleUpdatesClient, err := authorizationClient.SubscribeRoleUpdates(serverAuthContext(ctx, jwtService, serviceAuthName), &emptypb.Empty{})
	if err != nil {
		log.WithContext(ctx).Errorf("Unable to subscribe to role updates from authorization client. Retrying in %d seconds", retryAfter/time.Second)
		time.Sleep(retryAfter)
		ProcessRoleUpdates(ctx, authorizationClient, interceptor, jwtService, serviceAuthName)
		return
	}
	log.Info("Successfully subscribed to role updates from authorization server.")
	for {
		msg, err := roleUpdatesClient.Recv()
		if err == nil {
			log.Debugf("Update to role %d permissions. Clearing permissions cache for all users.", msg.Id)
			err = interceptor.ClearCache()
			if err != nil {
				log.WithContext(ctx).Warning("Clearing cache: %v", err)
			}
		} else if err == io.EOF {
			log.Infof("Role updates stream ended. Retrying in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			ProcessRoleUpdates(ctx, authorizationClient, interceptor, jwtService, serviceAuthName)
			return
		} else {
			log.WithContext(ctx).Errorf("Role Updates: %v.", err)
			log.Infof("Retrying connection in %d seconds", retryAfter/time.Second)
			time.Sleep(retryAfter)
			ProcessRoleUpdates(ctx, authorizationClient, interceptor, jwtService, serviceAuthName)
			return
		}
	}
}

func serverAuthContext(ctx context.Context, jwtService service.JWTService, authorizer string) context.Context {
	md := metadata.New(
		map[string]string{
			"authorization": fmt.Sprintf(
				"Bearer %s", generateTemporaryServerToken(ctx, jwtService, authorizer),
			),
		},
	)
	return metadata.NewOutgoingContext(context.Background(), md)
}

func DialOtelGrpc(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, InsecureOtelGrpcDialOpts()...)
}

func OtelGrpcDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	}
}

func InsecureOtelGrpcDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	}
}

func CreateGrpcServerWithAuth(ctx context.Context, jwt service.JWTService, accountsAddress string, serviceName string, publicRPCs map[string]struct{}) (*grpc.Server, *runtime.ServeMux, []grpc.DialOption, error) {
	if publicRPCs == nil {
		publicRPCs = make(map[string]struct{})
	}

	authName := fmt.Sprintf("sro.com/%s/v1", serviceName)

	conn, err := grpc.Dial(
		accountsAddress,
		InsecureOtelGrpcDialOpts()...,
	)

	if err != nil {
		return nil, nil, nil, err
	}

	authInterceptor := interceptor.NewAuthInterceptor(
		jwt,
		publicRPCs,
		GetPermissions(ctx, pb.NewAuthorizationServiceClient(conn), jwt, authName),
	)

	go ProcessRoleUpdates(ctx, pb.NewAuthorizationServiceClient(conn), authInterceptor, jwt, authName)
	go ProcessUserUpdates(ctx, pb.NewAuthorizationServiceClient(conn), authInterceptor, jwt, authName)

	return grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				authInterceptor.Unary(),
				helpers.UnaryLogRequest(),
				otelgrpc.UnaryServerInterceptor(),
			),
			grpc.ChainStreamInterceptor(
				authInterceptor.Stream(),
				helpers.StreamLogRequest(),
				otelgrpc.StreamServerInterceptor(),
			),
		),
		runtime.NewServeMux(),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		nil
}

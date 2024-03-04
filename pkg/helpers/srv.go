package helpers

import (
	"context"
	"strings"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	AuthorizationHeader = "authorization"
	AuthorizationScheme = "Bearer "
)

type BasicClaims struct {
	Subject string `json:"sub"`
}

var (
	jwtParser = jwt.NewParser()
)

func UnaryLogRequest() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Logger.WithContext(ctx).Info(info.FullMethod)
		return handler(ctx, req)
	}
}
func StreamLogRequest() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Logger.WithContext(stream.Context()).Info(info.FullMethod)
		return handler(srv, stream)
	}
}

func GrpcDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

func GrpcClientWithOtel(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, GrpcDialOpts()...)
}

func InitServerDefaults() (*grpc.Server, *runtime.ServeMux) {
	return grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				UnaryLogRequest(),
				otelgrpc.UnaryServerInterceptor(),
			),
			grpc.ChainStreamInterceptor(
				StreamLogRequest(),
				otelgrpc.StreamServerInterceptor(),
			)),
		runtime.NewServeMux()
}

func ExtractToken(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", status.Errorf(codes.Internal, "context is missing")
	}

	val := metautils.ExtractIncoming(ctx).Get(AuthorizationHeader)
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "request missing authorization")
	}

	if !strings.HasPrefix(val, AuthorizationScheme) {
		return "", status.Errorf(codes.Unauthenticated, "invalid authorization scheme. Expected %s.", AuthorizationScheme)
	}

	return val[len(AuthorizationScheme):], nil
}

func ExtractClaims(ctx context.Context) (*model.SROClaims, error) {
	if ctx == nil {
		return nil, status.Errorf(codes.Internal, "context is missing")
	}

	token, err := ExtractToken(ctx)
	if err != nil {
		return nil, err
	}

	jwtToken, _, err := jwtParser.ParseUnverified(token, &model.SROClaims{})
	if err != nil {
		log.Logger.WithContext(ctx).Infof("invalid token: %s", token)
		return nil, err
	}

	claims := jwtToken.Claims.(*model.SROClaims)
	return claims, nil
}

func VerifyClaims(ctx context.Context, client model.KeycloakClient, realm string) (*jwt.Token, *model.SROClaims, error) {
	if ctx == nil {
		return nil, nil, status.Errorf(codes.Internal, "context is missing")
	}

	if client == nil {
		return nil, nil, status.Errorf(codes.Internal, "gocloak is missing")
	}

	tokenString, err := ExtractToken(ctx)
	if err != nil {
		return nil, nil, err
	}

	var claims model.SROClaims
	token, err := client.DecodeAccessTokenCustomClaims(
		ctx,
		tokenString,
		realm,
		&claims,
	)

	if err != nil {
		log.Logger.WithContext(ctx).Errorf("extract claims: %v", err)
		return nil, nil, model.ErrUnauthorized
	}

	if !token.Valid {
		return nil, nil, model.ErrUnauthorized
	}

	return token, &claims, nil
}

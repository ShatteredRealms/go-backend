package helpers

import (
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"strings"
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
		log.WithContext(ctx).Info(info.FullMethod)
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
		log.WithContext(stream.Context()).Info(info.FullMethod)
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
				UnaryAddRolesToCtx(),
				otelgrpc.UnaryServerInterceptor(),
			),
			grpc.ChainStreamInterceptor(
				StreamLogRequest(),
				StreamAddRolesToCtx(),
				otelgrpc.StreamServerInterceptor(),
			)),
		runtime.NewServeMux()
}

func ExtractToken(ctx context.Context) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get(AuthorizationHeader)
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "request missing authorization")
	}

	if !strings.HasPrefix(val, AuthorizationScheme) {
		return "", status.Errorf(codes.Unauthenticated, "invalid authorization scheme. Expected %s.", AuthorizationScheme)
	}

	return strings.TrimPrefix(val, AuthorizationScheme), nil
}

func ExtractClaims(ctx context.Context) (*model.SROClaims, error) {
	token, err := ExtractToken(ctx)
	jwtToken, _, err := jwtParser.ParseUnverified(token, model.SROClaims{})
	if err != nil {
		return nil, err
	}

	claims := jwtToken.Claims.(model.SROClaims)
	return &claims, nil
}
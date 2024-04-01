package helpers

import (
	"context"
	"strings"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/srospan"
	"github.com/golang-jwt/jwt/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
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

func GrpcDialOpts() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

func GrpcClientWithOtel(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, GrpcDialOpts()...)
}

func InitServerDefaults() (*grpc.Server, *runtime.ServeMux) {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall),
		logging.WithCodes(logging.DefaultErrorToCode),
	}

	return grpc.NewServer(
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
			grpc.ChainUnaryInterceptor(
				logging.UnaryServerInterceptor(interceptorLogger(log.Logger), opts...),
			),
			grpc.ChainStreamInterceptor(
				logging.StreamServerInterceptor(interceptorLogger(log.Logger), opts...),
			)),
		runtime.NewServeMux()
}

func ExtractToken(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", model.ErrMissingContext
	}

	val := metautils.ExtractIncoming(ctx).Get(AuthorizationHeader)
	if val == "" {
		return "", model.ErrMissingAuthorization
	}

	if !strings.HasPrefix(val, AuthorizationScheme) {
		return "", model.ErrInvalidAuthorization
	}

	return val[len(AuthorizationScheme):], nil
}

func ExtractClaims(ctx context.Context) (*model.SROClaims, error) {
	if ctx == nil {
		return nil, model.ErrMissingContext
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
	if client == nil {
		return nil, nil, status.Errorf(codes.Internal, model.ErrMissingGocloak.Error())
	}

	tokenString, err := ExtractToken(ctx)
	if err != nil {
		return nil, nil, status.Errorf(codes.Unauthenticated, "invalid authentication")
	}

	var claims model.SROClaims
	token, err := client.DecodeAccessTokenCustomClaims(
		ctx,
		tokenString,
		realm,
		&claims,
	)

	if err != nil {
		log.Logger.WithContext(ctx).Infof("issues extracting claims: %v", err)
		return nil, nil, model.ErrUnauthorized.Err()
	}

	if !token.Valid {
		return nil, nil, model.ErrUnauthorized.Err()
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.SourceOwnerId(claims.Subject),
		srospan.SourceOwnerUsername(claims.Username),
	)

	return token, &claims, nil
}

// InterceptorLogger adapts logrus logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorLogger(l logrus.FieldLogger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make(map[string]any, len(fields)/2)
		i := logging.Fields(fields).Iterator()
		for i.Next() {
			k, v := i.At()
			f[k] = v
		}
		l := l.WithFields(f)

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg)
		case logging.LevelInfo:
			l.Info(msg)
		case logging.LevelWarn:
			l.Warn(msg)
		case logging.LevelError:
			l.Error(msg)
		default:
			l.Fatalf("unknown level %v", lvl)
		}
	})
}

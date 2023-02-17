package interceptor

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/allegro/bigcache/v3"
	"github.com/golang-jwt/jwt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

const (
	AuthorizationHeader = "authorization"
	AuthorizationScheme = "Bearer "
	AuthorizedOtherKey  = "sro:authOther"
)

var (
	ErrorUnauthorized = status.Errorf(codes.Unauthenticated, "Invalid permissions")
	tracer            = otel.Tracer("auth")
)

type AuthInterceptor struct {
	// The JWT service to use for verifying JWTs
	jwtService service.JWTService

	// publicRPCs is a map of all public gRPC functions that do not require permissions to be called
	publicRPCs map[string]struct{}

	// userPermissionsCache contains keys of usersnames and values of an array of their permissions they have access to
	userPermissionsCache *bigcache.BigCache

	// getUserPermissions function called when a users permissions are not in the cache. Should get the current
	// permissions for the user and put them in a map, the value of the "Other" field for the permission.
	getCurrentUserPermissions func(username string) map[string]bool

	tracer trace.Tracer
}

func NewAuthInterceptor(
	jwtService service.JWTService,
	publicRPCs map[string]struct{},
	getCurrentUserPermissions func(username string) map[string]bool,
) *AuthInterceptor {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		panic(err)
	}

	return &AuthInterceptor{
		jwtService:                jwtService,
		publicRPCs:                publicRPCs,
		userPermissionsCache:      cache,
		getCurrentUserPermissions: getCurrentUserPermissions,
		tracer:                    otel.Tracer("auth-interceptor"),
	}
}

func (interceptor *AuthInterceptor) updateUserPermissionsCache(username string, permissions map[string]bool) error {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(permissions)
	if err != nil {
		return err
	}

	return interceptor.userPermissionsCache.Set(username, buf.Bytes())
}

func (interceptor *AuthInterceptor) getCachedUserPermissions(username string) (map[string]bool, error) {
	raw, err := interceptor.userPermissionsCache.Get(username)
	if err != nil {
		return nil, err
	}

	var permissions map[string]bool
	buf := bytes.NewReader(raw)
	return permissions, gob.NewDecoder(buf).Decode(&permissions)
}

func (interceptor *AuthInterceptor) getUserPermissions(username string) map[string]bool {
	permissions, err := interceptor.getCachedUserPermissions(username)

	if err != nil {
		permissions = interceptor.getCurrentUserPermissions(username)
		_ = interceptor.updateUserPermissionsCache(username, permissions)
	}

	return permissions
}

func (interceptor *AuthInterceptor) ClearUserCache(username string) error {
	return interceptor.userPermissionsCache.Delete(username)
}

func (interceptor *AuthInterceptor) ClearCache() error {
	return interceptor.userPermissionsCache.Reset()
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		other, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(context.WithValue(ctx, AuthorizedOtherKey, other), req)
	}
}

func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		other, err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		ctx := context.WithValue(stream.Context(), AuthorizedOtherKey, other)

		return handler(srv, &grpc_middleware.WrappedServerStream{
			ServerStream:   stream,
			WrappedContext: ctx,
		})
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (bool, error) {
	ctx, span := interceptor.tracer.Start(ctx, "authorize")
	span.SetAttributes(attribute.String("method", method))

	if _, ok := interceptor.publicRPCs[method]; ok {
		return false, nil
	}

	// Get the token from the request
	token, err := ExtractAuthToken(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, "extract token")
		return false, err
	}

	// Get the username from the claim
	username, err := ExtractSubFromToken(ctx, token, interceptor.jwtService)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, "no sub")
		return false, status.Error(codes.Unauthenticated, err.Error())
	}

	// server sending request
	if username == "sro.com" {
		return true, nil
	}

	// Check the permission
	permissions := interceptor.getUserPermissions(username)
	if other, ok := permissions[method]; ok {
		return other, nil
	}

	span.RecordError(ErrorUnauthorized)
	span.SetStatus(otelcodes.Error, "unauthorized")
	return false, ErrorUnauthorized
}

func ExtractAuthToken(ctx context.Context) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get(AuthorizationHeader)
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "Request missing authorization")
	}

	if !strings.HasPrefix(val, AuthorizationScheme) {
		return "", status.Errorf(codes.Unauthenticated, "Invalid authorization scheme. Expected %s.", AuthorizationScheme)
	}

	return strings.TrimPrefix(val, AuthorizationScheme), nil
}

func ExtractSubFromToken(ctx context.Context, token string, jwtService service.JWTService) (string, error) {
	claims, err := jwtService.Validate(ctx, token)
	if err != nil {
		return "", fmt.Errorf("invalid authentication token")
	}

	if claims["sub"] == nil {
		return "", fmt.Errorf("token missing subject")
	}

	// Need to cast to float64 since that is JSON default for all numbers
	// SEE https://github.com/dgrijalva/jwt-go/issues/287
	username, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("unable to cast sub to string")
	}

	return username, nil
}

func AuthorizedForOther(ctx context.Context) bool {
	return ctx.Value(AuthorizedOtherKey).(bool)
}

// AuthorizedForTarget Checks the context for the jwt sub (account id) and checks if it matches the targetId.
// if it does match then it's authorized. Otherwise, checks if the ctx has been marked as authorized for
// other and returns true if it is. Should only be called after the interceptor.
func AuthorizedForTarget(ctx context.Context, jwtService service.JWTService, targetUsername string) (bool, error) {
	token, err := ExtractAuthToken(ctx)
	if err != nil {
		return false, err
	}

	subUsername, err := ExtractSubFromToken(ctx, token, jwtService)
	if err != nil {
		return false, err
	}

	if subUsername == targetUsername {
		return true, nil
	}

	return AuthorizedForOther(ctx), nil
}

func ExtractSubject(ctx context.Context, jwtService service.JWTService) (string, error) {
	token, err := ExtractAuthToken(ctx)
	if err != nil {
		return "", err
	}

	subUsername, err := ExtractSubFromToken(ctx, token, jwtService)
	if err != nil {
		return "", err
	}

	return subUsername, nil
}

func ExtractCtxClaims(ctx context.Context, jwtService service.JWTService) (jwt.MapClaims, error) {
	ctx, span := tracer.Start(ctx, "extract claims")
	token, err := ExtractAuthToken(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := jwtService.Validate(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid authentication token")
	}

	if claims["sub"] == nil {
		return nil, fmt.Errorf("token missing subject")
	}

	span.End()

	return claims, nil
}

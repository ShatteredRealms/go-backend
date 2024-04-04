package auth

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/WilSimpson/gocloak/v13"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
)

type claimContextKeyType int8

var (
	publicMethods = make(map[string]struct{}, 10)
)

func init() {
	publicMethods["/sro.HealthService/Health"] = struct{}{}
}

func RegisterPublicServiceMethods(methods ...string) {
	for _, method := range methods {
		publicMethods[method] = struct{}{}
	}
}

func AuthFunc(kcClient gocloak.KeycloakClient, realm string) auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		_, claims, err := verifyClaims(ctx, kcClient, realm)

		if err != nil {
			log.Logger.WithContext(ctx).Infof("Verifying claims failed: %s", err)
			return nil, common.ErrUnauthorized.Err()
		}

		return insertClaims(ctx, claims), nil
	}
}

func NotPublicServiceMatcher(_ context.Context, callMeta interceptors.CallMeta) bool {
	_, ok := publicMethods[callMeta.FullMethod()]
	log.Logger.Debugf("Verify Auth (%t): %s", !ok, callMeta.FullMethod())
	return !ok
}

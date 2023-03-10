package srv

import (
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotAuthorized = status.Error(codes.Unauthenticated, "not authorized")
)

func CtxHasRole(ctx context.Context, role string) error {
	var roles map[string]struct{}
	roles = ctx.Value(helpers.RolesCtxKey).(map[string]struct{})

	if _, ok := roles[role]; !ok {
		return ErrNotAuthorized
	}

	return nil
}

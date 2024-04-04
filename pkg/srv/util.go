package srv

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/WilSimpson/gocloak/v13"
	"go.opentelemetry.io/otel/attribute"
)

func createRoles(
	ctx context.Context,
	srvCtx *config.ServerContext,
	roles *[]*gocloak.Role,
) error {
	ctx, span := srvCtx.Tracer.Start(ctx, "roles.create")
	defer span.End()
	jwtToken, err := srvCtx.GetJWT(ctx)
	if err != nil {
		return err
	}

	var errs error
	for _, role := range *roles {
		_, err := srvCtx.KeycloakClient.CreateClientRole(
			ctx,
			jwtToken.AccessToken,
			srvCtx.GlobalConfig.Keycloak.Realm,
			srvCtx.RefSROServer.Keycloak.Id,
			*role,
		)

		// Code 409 is conflict
		if err != nil {
			if err.(*gocloak.APIError).Code != 409 {
				span.SetAttributes(attribute.String("role."+*role.Name, "error"))
				errs = errors.Join(errs, fmt.Errorf("creating role %s: %v", *role.Name, err))
			} else {
				span.SetAttributes(attribute.String("role."+*role.Name, "exists"))
			}
		} else {
			span.SetAttributes(attribute.String("role."+*role.Name, "created"))
		}
	}

	return errs
}

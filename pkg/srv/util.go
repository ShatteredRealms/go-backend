package srv

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotAuthorized = status.Error(codes.Unauthenticated, "not authorized")
)

func createRoles(
	ctx context.Context,
	token string,
	realm string,
	id string,
	roles *[]*gocloak.Role
) error {
	for _, role := range *roles {
		_, err = server.KeycloakClient.CreateClientRole(ctx, token, realm, id, *role)

		// Code 409 is conflict
		if err != nil && err.(*gocloak.APIError).Code != 409 {
			return fmt.Errorf("creating role %s: %v", role.Name, err)
		}
	}

	return nil
}

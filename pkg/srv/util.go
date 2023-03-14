package srv

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

var ()

func createRoles(
	ctx context.Context,
	client *gocloak.GoCloak,
	token string,
	realm string,
	id string,
	roles *[]*gocloak.Role,
) error {
	for _, role := range *roles {
		_, err := client.CreateClientRole(ctx, token, realm, id, *role)

		// Code 409 is conflict
		if err != nil && err.(*gocloak.APIError).Code != 409 {
			return fmt.Errorf("creating role %s: %v", *role.Name, err)
		}
	}

	return nil
}

package service

import "fmt"

// BaseUrl get the base url for the client
func (c *keycloakClient) baseUrl() string {
	return fmt.Sprintf("%s/realms/%s", c.Conf.BaseURL, c.Conf.Realm)
}

// BaseUrl get the base admin url for the client
func (c *keycloakClient) baseAdminUrl() string {
	return fmt.Sprintf("%s/admin/realms/%s", c.Conf.BaseURL, c.Conf.Realm)
}

// IntrospectUrl get the introspection url
func (c *keycloakClient) introspectUrl() string {
	return fmt.Sprintf(
		"%s/protocol/openid-connect/token/introspect",
		c.baseUrl(),
	)
}

// CreateClientRoleUrl get the url used for creating a client
func (c *keycloakClient) createClientRoleUrl() string {
	return fmt.Sprintf(
		"%s/clients/%s/roles",
		c.baseAdminUrl(),
		c.Conf.Id,
	)
}

func (c *keycloakClient) getAllClientRolesUrl() string {
	return fmt.Sprintf(
		"%s/clients/%s/roles",
		c.baseAdminUrl(),
		c.Conf.Id,
	)
}

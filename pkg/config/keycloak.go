package config

// KeycloakClientConfig oidc client for keycloak
type KeycloakClientConfig struct {
	Id           string `yaml:"id"`
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
}

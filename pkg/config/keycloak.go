package config

// KeycloakClientConfig oidc client for keycloak
type KeycloakClientConfig struct {
	Id           string `yaml:"id"`
	Realm        string `yaml:"realm"`
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	BaseURL      string `yaml:"baseURL"`
}

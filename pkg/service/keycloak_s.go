package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"io"
	"net/http"
	"strings"
)

const contentType = "application/json"

type keycloakClient struct {
	Conf   config.KeycloakClientConfig
	Client http.Client
}

func NewKeycloakClient(conf config.KeycloakClientConfig) *keycloakClient {
	return &keycloakClient{
		Conf:   conf,
		Client: http.Client{},
	}
}

func (c *keycloakClient) CreateRole(role *model.RoleRepresentation) error {
	body, err := structToJson(role)
	if err != nil {
		return fmt.Errorf("encoding role data: %v", err)
	}

	req, err := http.NewRequest("POST", c.createClientRoleUrl(), body)
	if err != nil {
		return fmt.Errorf("new request: %v", err)
	}
	c.addAuthorizationHeader(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("post new client roles status %d: %v", resp.StatusCode, err)
	}

	if resp.StatusCode != 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading response body: %v", err)
		}

		return fmt.Errorf("server responded with: %v", string(b))
	}

	return nil
}

func (c *keycloakClient) GetRoles() ([]model.RoleRepresentation, error) {
	req, err := http.NewRequest("GET", c.getAllClientRolesUrl(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}
	c.addAuthorizationHeader(req)
	fmt.Printf("auth header: %s\n", req.Header.Get("Authorization"))

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get roles status %d: %v", resp.StatusCode, err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading req resp: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("req resp: %s", string(bytes))
	}

	var roles []model.RoleRepresentation
	err = json.NewDecoder(strings.NewReader(string(bytes))).Decode(&roles)
	if err != nil {
		return nil, fmt.Errorf("decoding roles: %v", err)
	}

	return roles, nil
}

// CreateRolesNotExist creates all the given roles if none exist with the given name
func (c *keycloakClient) CreateRolesNotExist(roles ...*model.RoleRepresentation) error {
	currentRoles, err := c.GetRoles()
	if err != nil {
		return fmt.Errorf("get roles: %v", err)
	}

	currentRolesMap := make(map[string]struct{}, len(currentRoles))
	for _, r := range currentRoles {
		currentRolesMap[r.Name] = struct{}{}
	}

	for _, r := range roles {
		if _, ok := currentRolesMap[r.Name]; !ok {
			r.ClientRole = true
			r.ContainerId = c.Conf.Id
			r.Id = ""
			if err := c.CreateRole(r); err != nil {
				return fmt.Errorf("creating role: %v", err)
			}
		}
	}

	return nil
}

func (c *keycloakClient) addAuthorizationHeader(r *http.Request) {
	unencoded := fmt.Sprintf("%s:%s", c.Conf.ClientId, c.Conf.ClientSecret)
	r.Header.Set(
		"Authorization",
		fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(unencoded))),
	)
}

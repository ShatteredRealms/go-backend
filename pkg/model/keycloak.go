package model

// RoleRepresentation verbose representation of a role in keycloak. Used for creation of roles.
type RoleRepresentation struct {
	// Id unchangeable value used for unique identification
	Id string `json:"id,omitempty"`

	// Name used for referencing
	Name string `json:"name,omitempty"`

	// Description what the role is used for
	Description string `json:"description,omitempty"`

	// ClientRole whether the role is tied to a client
	ClientRole bool `json:"clientRole,omitempty"`

	// ContainerId internal client keycloak id
	ContainerId string `json:"containerId,omitempty"`

	// ClientId client id for the role. Not used in keycloak
	ClientId string `json:"clientId,omitempty"`
}

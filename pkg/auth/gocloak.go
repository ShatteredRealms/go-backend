package auth

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
)

// KeycloakClient holds all methods a client should fulfill
type KeycloakClient interface {
	RestyClient() *resty.Client
	SetRestyClient(restyClient *resty.Client)

	GetToken(ctx context.Context, realm string, options gocloak.TokenOptions) (*gocloak.JWT, error)
	GetRequestingPartyToken(ctx context.Context, token, realm string, options gocloak.RequestingPartyTokenOptions) (*gocloak.JWT, error)
	GetRequestingPartyPermissions(ctx context.Context, token, realm string, options gocloak.RequestingPartyTokenOptions) (*[]gocloak.RequestingPartyPermission, error)
	GetRequestingPartyPermissionDecision(ctx context.Context, token, realm string, options gocloak.RequestingPartyTokenOptions) (*gocloak.RequestingPartyPermissionDecision, error)

	Login(ctx context.Context, clientID, clientSecret, realm, username, password string) (*gocloak.JWT, error)
	LoginOtp(ctx context.Context, clientID, clientSecret, realm, username, password, totp string) (*gocloak.JWT, error)
	Logout(ctx context.Context, clientID, clientSecret, realm, refreshToken string) error
	LogoutPublicClient(ctx context.Context, clientID, realm, accessToken, refreshToken string) error
	LogoutAllSessions(ctx context.Context, accessToken, realm, userID string) error
	RevokeUserConsents(ctx context.Context, accessToken, realm, userID, clientID string) error
	LogoutUserSession(ctx context.Context, accessToken, realm, session string) error
	LoginClient(ctx context.Context, clientID, clientSecret, realm string) (*gocloak.JWT, error)
	LoginClientSignedJWT(ctx context.Context, clientID, realm string, key interface{}, signedMethod jwt.SigningMethod, expiresAt *jwt.NumericDate) (*gocloak.JWT, error)
	LoginAdmin(ctx context.Context, username, password, realm string) (*gocloak.JWT, error)
	RefreshToken(ctx context.Context, refreshToken, clientID, clientSecret, realm string) (*gocloak.JWT, error)
	DecodeAccessToken(ctx context.Context, accessToken, realm string) (*jwt.Token, *jwt.MapClaims, error)
	DecodeAccessTokenCustomClaims(ctx context.Context, accessToken, realm string, claims jwt.Claims) (*jwt.Token, error)
	RetrospectToken(ctx context.Context, accessToken, clientID, clientSecret, realm string) (*gocloak.IntroSpectTokenResult, error)
	GetIssuer(ctx context.Context, realm string) (*gocloak.IssuerResponse, error)
	GetCerts(ctx context.Context, realm string) (*gocloak.CertResponse, error)
	GetServerInfo(ctx context.Context, accessToken string) ([]*gocloak.ServerInfoRepresentation, error)
	GetUserInfo(ctx context.Context, accessToken, realm string) (*gocloak.UserInfo, error)
	GetRawUserInfo(ctx context.Context, accessToken, realm string) (map[string]interface{}, error)
	SetPassword(ctx context.Context, token, userID, realm, password string, temporary bool) error
	ExecuteActionsEmail(ctx context.Context, token, realm string, params gocloak.ExecuteActionsEmail) error

	CreateUser(ctx context.Context, token, realm string, user gocloak.User) (string, error)
	CreateGroup(ctx context.Context, accessToken, realm string, group gocloak.Group) (string, error)
	CreateChildGroup(ctx context.Context, token, realm, groupID string, group gocloak.Group) (string, error)
	CreateClientRole(ctx context.Context, accessToken, realm, idOfClient string, role gocloak.Role) (string, error)
	CreateClient(ctx context.Context, accessToken, realm string, newClient gocloak.Client) (string, error)
	CreateClientScope(ctx context.Context, accessToken, realm string, scope gocloak.ClientScope) (string, error)
	CreateComponent(ctx context.Context, accessToken, realm string, component gocloak.Component) (string, error)
	CreateClientScopeMappingsRealmRoles(ctx context.Context, token, realm, idOfClient string, roles []gocloak.Role) error
	CreateClientScopeMappingsClientRoles(ctx context.Context, token, realm, idOfClient, idOfSelectedClient string, roles []gocloak.Role) error
	CreateClientScopesScopeMappingsRealmRoles(ctx context.Context, token, realm, idOfCLientScope string, roles []gocloak.Role) error
	CreateClientScopesScopeMappingsClientRoles(ctx context.Context, token, realm, idOfClientScope, idOfClient string, roles []gocloak.Role) error

	UpdateUser(ctx context.Context, accessToken, realm string, user gocloak.User) error
	UpdateGroup(ctx context.Context, accessToken, realm string, updatedGroup gocloak.Group) error
	UpdateRole(ctx context.Context, accessToken, realm, idOfClient string, role gocloak.Role) error
	UpdateClient(ctx context.Context, accessToken, realm string, updatedClient gocloak.Client) error
	UpdateClientScope(ctx context.Context, accessToken, realm string, scope gocloak.ClientScope) error

	DeleteUser(ctx context.Context, accessToken, realm, userID string) error
	DeleteComponent(ctx context.Context, accessToken, realm, componentID string) error
	DeleteGroup(ctx context.Context, accessToken, realm, groupID string) error
	DeleteClientRole(ctx context.Context, accessToken, realm, idOfClient, roleName string) error
	DeleteClientRoleFromUser(ctx context.Context, token, realm, idOfClient, userID string, roles []gocloak.Role) error
	DeleteClient(ctx context.Context, accessToken, realm, idOfClient string) error
	DeleteClientScope(ctx context.Context, accessToken, realm, scopeID string) error
	DeleteClientScopeMappingsRealmRoles(ctx context.Context, token, realm, idOfClient string, roles []gocloak.Role) error
	DeleteClientScopeMappingsClientRoles(ctx context.Context, token, realm, idOfClient, idOfSelectedClient string, roles []gocloak.Role) error
	DeleteClientScopesScopeMappingsRealmRoles(ctx context.Context, token, realm, idOfCLientScope string, roles []gocloak.Role) error
	DeleteClientScopesScopeMappingsClientRoles(ctx context.Context, token, realm, idOfClientScope, ifOfClient string, roles []gocloak.Role) error

	GetClient(ctx context.Context, accessToken, realm, idOfClient string) (*gocloak.Client, error)
	GetClientsDefaultScopes(ctx context.Context, token, realm, idOfClient string) ([]*gocloak.ClientScope, error)
	AddDefaultScopeToClient(ctx context.Context, token, realm, idOfClient, scopeID string) error
	RemoveDefaultScopeFromClient(ctx context.Context, token, realm, idOfClient, scopeID string) error
	GetClientsOptionalScopes(ctx context.Context, token, realm, idOfClient string) ([]*gocloak.ClientScope, error)
	AddOptionalScopeToClient(ctx context.Context, token, realm, idOfClient, scopeID string) error
	RemoveOptionalScopeFromClient(ctx context.Context, token, realm, idOfClient, scopeID string) error
	GetDefaultOptionalClientScopes(ctx context.Context, token, realm string) ([]*gocloak.ClientScope, error)
	GetDefaultDefaultClientScopes(ctx context.Context, token, realm string) ([]*gocloak.ClientScope, error)
	GetClientScope(ctx context.Context, token, realm, scopeID string) (*gocloak.ClientScope, error)
	GetClientScopes(ctx context.Context, token, realm string) ([]*gocloak.ClientScope, error)
	GetClientScopeMappings(ctx context.Context, token, realm, idOfClient string) (*gocloak.MappingsRepresentation, error)
	GetClientScopeMappingsRealmRoles(ctx context.Context, token, realm, idOfClient string) ([]*gocloak.Role, error)
	GetClientScopeMappingsRealmRolesAvailable(ctx context.Context, token, realm, idOfClient string) ([]*gocloak.Role, error)
	GetClientScopesScopeMappingsRealmRolesAvailable(ctx context.Context, token, realm, idOfClientScope string) ([]*gocloak.Role, error)
	GetClientScopesScopeMappingsClientRolesAvailable(ctx context.Context, token, realm, idOfClientScope, idOfClient string) ([]*gocloak.Role, error)
	GetClientScopeMappingsClientRoles(ctx context.Context, token, realm, idOfClient, idOfSelectedClient string) ([]*gocloak.Role, error)
	GetClientScopesScopeMappingsRealmRoles(ctx context.Context, token, realm, idOfClientScope string) ([]*gocloak.Role, error)
	GetClientScopesScopeMappingsClientRoles(ctx context.Context, token, realm, idOfClientScope, idOfClient string) ([]*gocloak.Role, error)
	GetClientScopeMappingsClientRolesAvailable(ctx context.Context, token, realm, idOfClient, idOfSelectedClient string) ([]*gocloak.Role, error)
	GetClientSecret(ctx context.Context, token, realm, idOfClient string) (*gocloak.CredentialRepresentation, error)
	GetClientServiceAccount(ctx context.Context, token, realm, idOfClient string) (*gocloak.User, error)
	RegenerateClientSecret(ctx context.Context, token, realm, idOfClient string) (*gocloak.CredentialRepresentation, error)
	GetKeyStoreConfig(ctx context.Context, accessToken, realm string) (*gocloak.KeyStoreConfig, error)
	GetUserByID(ctx context.Context, accessToken, realm, userID string) (*gocloak.User, error)
	GetUserCount(ctx context.Context, accessToken, realm string, params gocloak.GetUsersParams) (int, error)
	GetUsers(ctx context.Context, accessToken, realm string, params gocloak.GetUsersParams) ([]*gocloak.User, error)
	GetUserGroups(ctx context.Context, token, realm, userID string, params gocloak.GetGroupsParams) ([]*gocloak.Group, error)
	AddUserToGroup(ctx context.Context, token, realm, userID, groupID string) error
	DeleteUserFromGroup(ctx context.Context, token, realm, userID, groupID string) error
	GetComponents(ctx context.Context, accessToken, realm string) ([]*gocloak.Component, error)
	GetGroups(ctx context.Context, accessToken, realm string, params gocloak.GetGroupsParams) ([]*gocloak.Group, error)
	GetGroupsCount(ctx context.Context, token, realm string, params gocloak.GetGroupsParams) (int, error)
	GetGroup(ctx context.Context, accessToken, realm, groupID string) (*gocloak.Group, error)
	GetDefaultGroups(ctx context.Context, accessToken, realm string) ([]*gocloak.Group, error)
	AddDefaultGroup(ctx context.Context, accessToken, realm, groupID string) error
	RemoveDefaultGroup(ctx context.Context, accessToken, realm, groupID string) error
	GetGroupMembers(ctx context.Context, accessToken, realm, groupID string, params gocloak.GetGroupsParams) ([]*gocloak.User, error)
	GetRoleMappingByGroupID(ctx context.Context, accessToken, realm, groupID string) (*gocloak.MappingsRepresentation, error)
	GetRoleMappingByUserID(ctx context.Context, accessToken, realm, userID string) (*gocloak.MappingsRepresentation, error)
	GetClientRoles(ctx context.Context, accessToken, realm, idOfClient string, params gocloak.GetRoleParams) ([]*gocloak.Role, error)
	GetClientRole(ctx context.Context, token, realm, idOfClient, roleName string) (*gocloak.Role, error)
	GetClientRoleByID(ctx context.Context, accessToken, realm, roleID string) (*gocloak.Role, error)
	GetClients(ctx context.Context, accessToken, realm string, params gocloak.GetClientsParams) ([]*gocloak.Client, error)
	AddClientRoleComposite(ctx context.Context, token, realm, roleID string, roles []gocloak.Role) error
	DeleteClientRoleComposite(ctx context.Context, token, realm, roleID string, roles []gocloak.Role) error
	GetUsersByRoleName(ctx context.Context, token, realm, roleName string, params gocloak.GetUsersByRoleParams) ([]*gocloak.User, error)
	GetUsersByClientRoleName(ctx context.Context, token, realm, idOfClient, roleName string, params gocloak.GetUsersByRoleParams) ([]*gocloak.User, error)
	CreateClientProtocolMapper(ctx context.Context, token, realm, idOfClient string, mapper gocloak.ProtocolMapperRepresentation) (string, error)
	UpdateClientProtocolMapper(ctx context.Context, token, realm, idOfClient, mapperID string, mapper gocloak.ProtocolMapperRepresentation) error
	DeleteClientProtocolMapper(ctx context.Context, token, realm, idOfClient, mapperID string) error

	// *** Realm Roles ***

	CreateRealmRole(ctx context.Context, token, realm string, role gocloak.Role) (string, error)
	GetRealmRole(ctx context.Context, token, realm, roleName string) (*gocloak.Role, error)
	GetRealmRoles(ctx context.Context, accessToken, realm string, params gocloak.GetRoleParams) ([]*gocloak.Role, error)
	GetRealmRoleByID(ctx context.Context, token, realm, roleID string) (*gocloak.Role, error)
	GetRealmRolesByUserID(ctx context.Context, accessToken, realm, userID string) ([]*gocloak.Role, error)
	GetRealmRolesByGroupID(ctx context.Context, accessToken, realm, groupID string) ([]*gocloak.Role, error)
	UpdateRealmRole(ctx context.Context, token, realm, roleName string, role gocloak.Role) error
	UpdateRealmRoleByID(ctx context.Context, token, realm, roleID string, role gocloak.Role) error
	DeleteRealmRole(ctx context.Context, token, realm, roleName string) error
	AddRealmRoleToUser(ctx context.Context, token, realm, userID string, roles []gocloak.Role) error
	DeleteRealmRoleFromUser(ctx context.Context, token, realm, userID string, roles []gocloak.Role) error
	AddRealmRoleToGroup(ctx context.Context, token, realm, groupID string, roles []gocloak.Role) error
	DeleteRealmRoleFromGroup(ctx context.Context, token, realm, groupID string, roles []gocloak.Role) error
	AddRealmRoleComposite(ctx context.Context, token, realm, roleName string, roles []gocloak.Role) error
	DeleteRealmRoleComposite(ctx context.Context, token, realm, roleName string, roles []gocloak.Role) error
	GetCompositeRealmRoles(ctx context.Context, token, realm, roleName string) ([]*gocloak.Role, error)
	GetCompositeRealmRolesByRoleID(ctx context.Context, token, realm, roleID string) ([]*gocloak.Role, error)
	GetCompositeRealmRolesByUserID(ctx context.Context, token, realm, userID string) ([]*gocloak.Role, error)
	GetCompositeRealmRolesByGroupID(ctx context.Context, token, realm, groupID string) ([]*gocloak.Role, error)
	GetAvailableRealmRolesByUserID(ctx context.Context, token, realm, userID string) ([]*gocloak.Role, error)
	GetAvailableRealmRolesByGroupID(ctx context.Context, token, realm, groupID string) ([]*gocloak.Role, error)

	// *** Client Roles ***

	AddClientRoleToUser(ctx context.Context, token, realm, idOfClient, userID string, roles []gocloak.Role) error
	AddClientRoleToGroup(ctx context.Context, token, realm, idOfClient, groupID string, roles []gocloak.Role) error
	DeleteClientRoleFromGroup(ctx context.Context, token, realm, idOfClient, groupID string, roles []gocloak.Role) error
	GetCompositeClientRolesByRoleID(ctx context.Context, token, realm, idOfClient, roleID string) ([]*gocloak.Role, error)
	GetClientRolesByUserID(ctx context.Context, token, realm, idOfClient, userID string) ([]*gocloak.Role, error)
	GetClientRolesByGroupID(ctx context.Context, token, realm, idOfClient, groupID string) ([]*gocloak.Role, error)
	GetCompositeClientRolesByUserID(ctx context.Context, token, realm, idOfClient, userID string) ([]*gocloak.Role, error)
	GetCompositeClientRolesByGroupID(ctx context.Context, token, realm, idOfClient, groupID string) ([]*gocloak.Role, error)
	GetAvailableClientRolesByUserID(ctx context.Context, token, realm, idOfClient, userID string) ([]*gocloak.Role, error)
	GetAvailableClientRolesByGroupID(ctx context.Context, token, realm, idOfClient, groupID string) ([]*gocloak.Role, error)

	// *** Realm ***

	GetRealm(ctx context.Context, token, realm string) (*gocloak.RealmRepresentation, error)
	GetRealms(ctx context.Context, token string) ([]*gocloak.RealmRepresentation, error)
	CreateRealm(ctx context.Context, token string, realm gocloak.RealmRepresentation) (string, error)
	UpdateRealm(ctx context.Context, token string, realm gocloak.RealmRepresentation) error
	DeleteRealm(ctx context.Context, token, realm string) error
	ClearRealmCache(ctx context.Context, token, realm string) error
	ClearUserCache(ctx context.Context, token, realm string) error
	ClearKeysCache(ctx context.Context, token, realm string) error

	GetClientUserSessions(ctx context.Context, token, realm, idOfClient string) ([]*gocloak.UserSessionRepresentation, error)
	GetClientOfflineSessions(ctx context.Context, token, realm, idOfClient string) ([]*gocloak.UserSessionRepresentation, error)
	GetUserSessions(ctx context.Context, token, realm, userID string) ([]*gocloak.UserSessionRepresentation, error)
	GetUserOfflineSessionsForClient(ctx context.Context, token, realm, userID, idOfClient string) ([]*gocloak.UserSessionRepresentation, error)

	// *** Protection API ***
	GetResource(ctx context.Context, token, realm, idOfClient, resourceID string) (*gocloak.ResourceRepresentation, error)
	GetResources(ctx context.Context, token, realm, idOfClient string, params gocloak.GetResourceParams) ([]*gocloak.ResourceRepresentation, error)
	CreateResource(ctx context.Context, token, realm, idOfClient string, resource gocloak.ResourceRepresentation) (*gocloak.ResourceRepresentation, error)
	UpdateResource(ctx context.Context, token, realm, idOfClient string, resource gocloak.ResourceRepresentation) error
	DeleteResource(ctx context.Context, token, realm, idOfClient, resourceID string) error

	GetResourceClient(ctx context.Context, token, realm, resourceID string) (*gocloak.ResourceRepresentation, error)
	GetResourcesClient(ctx context.Context, token, realm string, params gocloak.GetResourceParams) ([]*gocloak.ResourceRepresentation, error)
	CreateResourceClient(ctx context.Context, token, realm string, resource gocloak.ResourceRepresentation) (*gocloak.ResourceRepresentation, error)
	UpdateResourceClient(ctx context.Context, token, realm string, resource gocloak.ResourceRepresentation) error
	DeleteResourceClient(ctx context.Context, token, realm, resourceID string) error

	GetScope(ctx context.Context, token, realm, idOfClient, scopeID string) (*gocloak.ScopeRepresentation, error)
	GetScopes(ctx context.Context, token, realm, idOfClient string, params gocloak.GetScopeParams) ([]*gocloak.ScopeRepresentation, error)
	CreateScope(ctx context.Context, token, realm, idOfClient string, scope gocloak.ScopeRepresentation) (*gocloak.ScopeRepresentation, error)
	UpdateScope(ctx context.Context, token, realm, idOfClient string, resource gocloak.ScopeRepresentation) error
	DeleteScope(ctx context.Context, token, realm, idOfClient, scopeID string) error

	GetPolicy(ctx context.Context, token, realm, idOfClient, policyID string) (*gocloak.PolicyRepresentation, error)
	GetPolicies(ctx context.Context, token, realm, idOfClient string, params gocloak.GetPolicyParams) ([]*gocloak.PolicyRepresentation, error)
	CreatePolicy(ctx context.Context, token, realm, idOfClient string, policy gocloak.PolicyRepresentation) (*gocloak.PolicyRepresentation, error)
	UpdatePolicy(ctx context.Context, token, realm, idOfClient string, policy gocloak.PolicyRepresentation) error
	DeletePolicy(ctx context.Context, token, realm, idOfClient, policyID string) error

	GetResourcePolicy(ctx context.Context, token, realm, permissionID string) (*gocloak.ResourcePolicyRepresentation, error)
	GetResourcePolicies(ctx context.Context, token, realm string, params gocloak.GetResourcePoliciesParams) ([]*gocloak.ResourcePolicyRepresentation, error)
	CreateResourcePolicy(ctx context.Context, token, realm, resourceID string, policy gocloak.ResourcePolicyRepresentation) (*gocloak.ResourcePolicyRepresentation, error)
	UpdateResourcePolicy(ctx context.Context, token, realm, permissionID string, policy gocloak.ResourcePolicyRepresentation) error
	DeleteResourcePolicy(ctx context.Context, token, realm, permissionID string) error

	GetPermission(ctx context.Context, token, realm, idOfClient, permissionID string) (*gocloak.PermissionRepresentation, error)
	GetPermissions(ctx context.Context, token, realm, idOfClient string, params gocloak.GetPermissionParams) ([]*gocloak.PermissionRepresentation, error)
	GetPermissionResources(ctx context.Context, token, realm, idOfClient, permissionID string) ([]*gocloak.PermissionResource, error)
	GetPermissionScopes(ctx context.Context, token, realm, idOfClient, permissionID string) ([]*gocloak.PermissionScope, error)
	GetDependentPermissions(ctx context.Context, token, realm, idOfClient, policyID string) ([]*gocloak.PermissionRepresentation, error)
	CreatePermission(ctx context.Context, token, realm, idOfClient string, permission gocloak.PermissionRepresentation) (*gocloak.PermissionRepresentation, error)
	UpdatePermission(ctx context.Context, token, realm, idOfClient string, permission gocloak.PermissionRepresentation) error
	DeletePermission(ctx context.Context, token, realm, idOfClient, permissionID string) error

	CreatePermissionTicket(ctx context.Context, token, realm string, permissions []gocloak.CreatePermissionTicketParams) (*gocloak.PermissionTicketResponseRepresentation, error)
	GrantUserPermission(ctx context.Context, token, realm string, permission gocloak.PermissionGrantParams) (*gocloak.PermissionGrantResponseRepresentation, error)
	UpdateUserPermission(ctx context.Context, token, realm string, permission gocloak.PermissionGrantParams) (*gocloak.PermissionGrantResponseRepresentation, error)
	GetUserPermissions(ctx context.Context, token, realm string, params gocloak.GetUserPermissionParams) ([]*gocloak.PermissionGrantResponseRepresentation, error)
	DeleteUserPermission(ctx context.Context, token, realm, ticketID string) error

	// *** Credentials API ***

	GetCredentialRegistrators(ctx context.Context, token, realm string) ([]string, error)
	GetConfiguredUserStorageCredentialTypes(ctx context.Context, token, realm, userID string) ([]string, error)
	GetCredentials(ctx context.Context, token, realm, UserID string) ([]*gocloak.CredentialRepresentation, error)
	DeleteCredentials(ctx context.Context, token, realm, UserID, CredentialID string) error
	UpdateCredentialUserLabel(ctx context.Context, token, realm, userID, credentialID, userLabel string) error
	DisableAllCredentialsByType(ctx context.Context, token, realm, userID string, types []string) error
	MoveCredentialBehind(ctx context.Context, token, realm, userID, credentialID, newPreviousCredentialID string) error
	MoveCredentialToFirst(ctx context.Context, token, realm, userID, credentialID string) error

	// *** Authentication Flows ***
	GetAuthenticationFlows(ctx context.Context, token, realm string) ([]*gocloak.AuthenticationFlowRepresentation, error)
	GetAuthenticationFlow(ctx context.Context, token, realm string, authenticationFlowID string) (*gocloak.AuthenticationFlowRepresentation, error)
	CreateAuthenticationFlow(ctx context.Context, token, realm string, flow gocloak.AuthenticationFlowRepresentation) error
	UpdateAuthenticationFlow(ctx context.Context, token, realm string, flow gocloak.AuthenticationFlowRepresentation, authenticationFlowID string) (*gocloak.AuthenticationFlowRepresentation, error)
	DeleteAuthenticationFlow(ctx context.Context, token, realm, flowID string) error

	// *** Identity Providers ***

	CreateIdentityProvider(ctx context.Context, token, realm string, providerRep gocloak.IdentityProviderRepresentation) (string, error)
	GetIdentityProvider(ctx context.Context, token, realm, alias string) (*gocloak.IdentityProviderRepresentation, error)
	GetIdentityProviders(ctx context.Context, token, realm string) ([]*gocloak.IdentityProviderRepresentation, error)
	UpdateIdentityProvider(ctx context.Context, token, realm, alias string, providerRep gocloak.IdentityProviderRepresentation) error
	DeleteIdentityProvider(ctx context.Context, token, realm, alias string) error

	CreateIdentityProviderMapper(ctx context.Context, token, realm, alias string, mapper gocloak.IdentityProviderMapper) (string, error)
	GetIdentityProviderMapper(ctx context.Context, token string, realm string, alias string, mapperID string) (*gocloak.IdentityProviderMapper, error)
	CreateUserFederatedIdentity(ctx context.Context, token, realm, userID, providerID string, federatedIdentityRep gocloak.FederatedIdentityRepresentation) error
	GetUserFederatedIdentities(ctx context.Context, token, realm, userID string) ([]*gocloak.FederatedIdentityRepresentation, error)
	DeleteUserFederatedIdentity(ctx context.Context, token, realm, userID, providerID string) error

	// *** Events API ***
	GetEvents(ctx context.Context, token string, realm string, params gocloak.GetEventsParams) ([]*gocloak.EventRepresentation, error)
}

package internal

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/WilSimpson/gocloak/v13"
)

type KeycloakManager struct {
	Config     *config.GlobalConfig
	Admins     []*gocloak.User
	Users      []*gocloak.User
	GuestCount int

	keycloak gocloak.KeycloakClient
	token    *gocloak.JWT
}

func NewKeycloakManager(
	ctx context.Context,
	conf *config.GlobalConfig,
	numAdmins, numUsers, numGuests int,
) *KeycloakManager {
	return &KeycloakManager{
		Config:     conf,
		Admins:     make([]*gocloak.User, numAdmins),
		Users:      make([]*gocloak.User, numUsers),
		GuestCount: numGuests,
		keycloak:   gocloak.NewClient(conf.Keycloak.BaseURL),
	}
}

func (km *KeycloakManager) Setup(ctx context.Context) error {
	err := km.newToken(ctx)
	if err != nil {
		return fmt.Errorf("login keycloak: %v", err)
	}

	saRole, err := km.keycloak.GetRealmRole(context.Background(), km.token.AccessToken, km.Config.Keycloak.Realm, "super admin")
	if err != nil {
		return fmt.Errorf("getting sa role: %w", err)
	}

	userRole, err := km.keycloak.GetRealmRole(context.Background(), km.token.AccessToken, km.Config.Keycloak.Realm, "user")
	if err != nil {
		return fmt.Errorf("getting user role: %w", err)
	}

	for idx := range km.Admins {
		km.Admins[idx], err = km.createUserWithRole(ctx, idx, saRole)
		if err != nil {
			return fmt.Errorf("create admin: %w", err)
		}
	}

	for idx := range km.Users {
		km.Users[idx], err = km.createUserWithRole(ctx, idx, userRole)
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}
	}

	return nil
}

func (km *KeycloakManager) createUserWithRole(ctx context.Context, count int, role *gocloak.Role) (user *gocloak.User, err error) {
	name := "st" + strings.ReplaceAll(*role.Name, " ", "") + strconv.Itoa(count)
	user = &gocloak.User{
		ID:            new(string),
		Username:      gocloak.StringP(name),
		Enabled:       gocloak.BoolP(true),
		Totp:          gocloak.BoolP(true),
		EmailVerified: gocloak.BoolP(true),
		FirstName:     gocloak.StringP(name),
		LastName:      gocloak.StringP(""),
		Email:         gocloak.StringP("stresstest@shatteredrealms.online"),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP("Password1!"),
			},
		},
	}

	*user.ID, err = km.keycloak.CreateUser(ctx, km.token.AccessToken, km.Config.Keycloak.Realm, *user)
	if err != nil {
		return nil, err
	}

	return user, km.keycloak.AddRealmRoleToUser(
		ctx,
		km.token.AccessToken,
		km.Config.Keycloak.Realm,
		*user.ID,
		[]gocloak.Role{*role},
	)
}

func (km *KeycloakManager) newToken(ctx context.Context) (err error) {
	km.token, err = km.keycloak.LoginClient(
		ctx,
		km.Config.Character.Keycloak.ClientId,
		km.Config.Character.Keycloak.ClientSecret,
		km.Config.Keycloak.Realm,
	)

	return err
}

func (km *KeycloakManager) Shutdown(ctx context.Context) error {
	err := km.newToken(ctx)
	if err != nil {
		return fmt.Errorf("login keycloak: %v", err)
	}

	fmt.Println("Waiting 1 second for users to syncronize")
	time.Sleep(time.Second)

	// Need to search since there could be a disconnect if cancelation happens during setup
	users, err := km.keycloak.GetUsers(ctx, km.token.AccessToken, km.Config.Keycloak.Realm, gocloak.GetUsersParams{
		Email: gocloak.StringP("stresstest@shatteredrealms.online"),
		Max:   gocloak.IntP(len(km.Admins) + len(km.Users)),
	})
	if err != nil {
		return fmt.Errorf("unable to get stress test users: %w", err)
	}

	for _, user := range users {
		if user != nil {
			err = errors.Join(err, km.keycloak.DeleteUser(ctx, km.token.AccessToken, km.Config.Keycloak.Realm, *user.ID))
		}
	}

	return err
}

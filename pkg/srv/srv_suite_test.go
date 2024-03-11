package srv_test

import (
	"context"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
)

const ()

var (
	keycloak *gocloak.GoCloak
	conf     *config.GlobalConfig
	admin    = gocloak.User{
		ID:            new(string),
		Username:      gocloak.StringP("testadmin"),
		Enabled:       gocloak.BoolP(true),
		Totp:          gocloak.BoolP(true),
		EmailVerified: gocloak.BoolP(true),
		FirstName:     gocloak.StringP("adminfirstname"),
		LastName:      gocloak.StringP("adminlastname"),
		Email:         gocloak.StringP("admin@example.com"),
		RealmRoles:    &[]string{"super admin"},
		Credentials: &[]gocloak.CredentialRepresentation{
			gocloak.CredentialRepresentation{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP("Password1!"),
			},
		},
	}
	player = gocloak.User{
		ID:            new(string),
		Username:      gocloak.StringP("testplayer"),
		Enabled:       gocloak.BoolP(true),
		Totp:          gocloak.BoolP(true),
		EmailVerified: gocloak.BoolP(true),
		FirstName:     gocloak.StringP("playerfirstname"),
		LastName:      gocloak.StringP("playerlastname"),
		Email:         gocloak.StringP("player@example.com"),
		RealmRoles:    &[]string{"user", "public"},
		Credentials: &[]gocloak.CredentialRepresentation{
			gocloak.CredentialRepresentation{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP("Password1!"),
			},
		},
	}

	adminToken  *gocloak.JWT
	playerToken *gocloak.JWT
	clientToken *gocloak.JWT

	incAdminCtx  context.Context
	incPlayerCtx context.Context
	incClientCtx context.Context
)

func TestSrv(t *testing.T) {
	var closeFunc func()
	BeforeSuite(func() {
		closeFunc, keycloak = testdb.SetupKeycloakWithDocker()
		Expect(keycloak).NotTo(BeNil())
		conf = config.NewGlobalConfig(context.Background())
		token, err := keycloak.LoginClient(
			context.Background(),
			conf.Character.Keycloak.ClientId,
			conf.Character.Keycloak.ClientSecret,
			conf.Character.Keycloak.Realm,
		)
		Expect(err).NotTo(HaveOccurred())

		*admin.ID, err = keycloak.CreateUser(context.Background(), token.AccessToken, conf.Character.Keycloak.Realm, admin)
		Expect(err).NotTo(HaveOccurred())
		*player.ID, err = keycloak.CreateUser(context.Background(), token.AccessToken, conf.Character.Keycloak.Realm, player)
		Expect(err).NotTo(HaveOccurred())

		adminToken, err = keycloak.GetToken(context.Background(), conf.Character.Keycloak.Realm, gocloak.TokenOptions{
			ClientID:     &conf.Character.Keycloak.ClientId,
			ClientSecret: &conf.Character.Keycloak.ClientSecret,
			GrantType:    gocloak.StringP("password"),
			Username:     admin.Username,
			Password:     gocloak.StringP("Password1!"),
		})
		Expect(err).NotTo(HaveOccurred())
		playerToken, err = keycloak.GetToken(context.Background(), conf.Character.Keycloak.Realm, gocloak.TokenOptions{
			ClientID:     &conf.Character.Keycloak.ClientId,
			ClientSecret: &conf.Character.Keycloak.ClientSecret,
			GrantType:    gocloak.StringP("password"),
			Username:     player.Username,
			Password:     gocloak.StringP("Password1!"),
		})
		Expect(err).NotTo(HaveOccurred())
		clientToken, err = keycloak.GetToken(context.Background(), conf.Character.Keycloak.Realm, gocloak.TokenOptions{
			ClientID:     &conf.Character.Keycloak.ClientId,
			ClientSecret: &conf.Character.Keycloak.ClientSecret,
			GrantType:    gocloak.StringP("client_credentials"),
		})
		Expect(err).NotTo(HaveOccurred())

		md := metadata.New(
			map[string]string{
				"authorization": adminToken.AccessToken,
			},
		)
		incAdminCtx = metadata.NewIncomingContext(context.Background(), md)
		md = metadata.New(
			map[string]string{
				"authorization": playerToken.AccessToken,
			},
		)
		incPlayerCtx = metadata.NewIncomingContext(context.Background(), md)
		md = metadata.New(
			map[string]string{
				"authorization": clientToken.AccessToken,
			},
		)
		incClientCtx = metadata.NewIncomingContext(context.Background(), md)
	})

	AfterSuite(func() {
		closeFunc()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Srv Suite")
}

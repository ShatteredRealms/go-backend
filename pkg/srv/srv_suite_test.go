package srv_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
)

const (
	topicName = "testtopic"
)

var (
	// Global
	conf    *config.GlobalConfig
	fakeErr = fmt.Errorf("error")

	// Keycloak
	keycloak *gocloak.GoCloak
	admin    = gocloak.User{
		ID:            new(string),
		Username:      gocloak.StringP("testadmin"),
		Enabled:       gocloak.BoolP(true),
		Totp:          gocloak.BoolP(true),
		EmailVerified: gocloak.BoolP(true),
		FirstName:     gocloak.StringP("adminfirstname"),
		LastName:      gocloak.StringP("adminlastname"),
		Email:         gocloak.StringP("admin@example.com"),
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
		Credentials: &[]gocloak.CredentialRepresentation{
			gocloak.CredentialRepresentation{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP("Password1!"),
			},
		},
	}
	guest = gocloak.User{
		ID:            new(string),
		Username:      gocloak.StringP("testguest"),
		Enabled:       gocloak.BoolP(true),
		Totp:          gocloak.BoolP(true),
		EmailVerified: gocloak.BoolP(true),
		FirstName:     gocloak.StringP("guestfirstname"),
		LastName:      gocloak.StringP("guestlastname"),
		Email:         gocloak.StringP("guest@example.com"),
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
	guestToken  *gocloak.JWT

	incAdminCtx  context.Context
	incPlayerCtx context.Context
	incClientCtx context.Context
	incGuestCtx  context.Context

	// Kafka
	kafkaPort uint
)

func TestSrv(t *testing.T) {
	var keycloakCloseFunc func()
	var kafkaCloseFunc func()
	var err error

	SynchronizedBeforeSuite(func() []byte {
		var host string
		keycloakCloseFunc, host = testdb.SetupKeycloakWithDocker()
		Expect(host).NotTo(BeNil())

		keycloak = gocloak.NewClient(string(host))

		conf = config.NewGlobalConfig(context.Background())
		clientToken, err = keycloak.LoginClient(
			context.Background(),
			conf.Character.Keycloak.ClientId,
			conf.Character.Keycloak.ClientSecret,
			conf.Character.Keycloak.Realm,
		)
		Expect(err).NotTo(HaveOccurred())

		// Eventually(func() error {
		*admin.ID, err = keycloak.CreateUser(context.Background(), clientToken.AccessToken, conf.Character.Keycloak.Realm, admin)
		Expect(err).NotTo(HaveOccurred())
		// }).Within(time.Minute).ProbeEvery(time.Second).ShouldNot(HaveOccurred())
		*player.ID, err = keycloak.CreateUser(context.Background(), clientToken.AccessToken, conf.Character.Keycloak.Realm, player)
		Expect(err).NotTo(HaveOccurred())
		*guest.ID, err = keycloak.CreateUser(context.Background(), clientToken.AccessToken, conf.Character.Keycloak.Realm, guest)
		Expect(err).NotTo(HaveOccurred())

		saRole, err := keycloak.GetRealmRole(context.Background(), clientToken.AccessToken, conf.Character.Keycloak.Realm, "super admin")
		Expect(err).NotTo(HaveOccurred())
		userRole, err := keycloak.GetRealmRole(context.Background(), clientToken.AccessToken, conf.Character.Keycloak.Realm, "user")
		Expect(err).NotTo(HaveOccurred())

		err = keycloak.AddRealmRoleToUser(
			context.Background(),
			clientToken.AccessToken,
			conf.Character.Keycloak.Realm,
			*admin.ID,
			[]gocloak.Role{*saRole},
		)
		Expect(err).NotTo(HaveOccurred())
		err = keycloak.AddRealmRoleToUser(
			context.Background(),
			clientToken.AccessToken,
			conf.Character.Keycloak.Realm,
			*player.ID,
			[]gocloak.Role{*userRole},
		)
		Expect(err).NotTo(HaveOccurred())

		var kafkaPort uint
		kafkaCloseFunc, kafkaPort = testdb.SetupKafkaWithDocker()

		out := fmt.Sprintf("%s\n%d", host, kafkaPort)

		return []byte(out)
	}, func(data []byte) {
		splitData := strings.Split(string(data), "\n")
		Expect(splitData).To(HaveLen(2))

		host := splitData[0]
		kafkaPort64, err := strconv.ParseUint(splitData[1], 10, 32)
		Expect(err).NotTo(HaveOccurred())
		kafkaPort = uint(kafkaPort64)

		keycloak = gocloak.NewClient(string(host))
		conf = config.NewGlobalConfig(context.Background())

		clientToken, err = keycloak.LoginClient(
			context.Background(),
			conf.Character.Keycloak.ClientId,
			conf.Character.Keycloak.ClientSecret,
			conf.Character.Keycloak.Realm,
		)
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
		guestToken, err = keycloak.GetToken(context.Background(), conf.Character.Keycloak.Realm, gocloak.TokenOptions{
			ClientID:     &conf.Character.Keycloak.ClientId,
			ClientSecret: &conf.Character.Keycloak.ClientSecret,
			GrantType:    gocloak.StringP("password"),
			Username:     guest.Username,
			Password:     gocloak.StringP("Password1!"),
		})
		Expect(err).NotTo(HaveOccurred())

		admins, err := keycloak.GetUsers(
			context.Background(),
			clientToken.AccessToken,
			conf.Character.Keycloak.Realm,
			gocloak.GetUsersParams{Username: admin.Username},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(admins).To(HaveLen(1))
		admin = *admins[0]

		players, err := keycloak.GetUsers(
			context.Background(),
			clientToken.AccessToken,
			conf.Character.Keycloak.Realm,
			gocloak.GetUsersParams{Username: player.Username},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(players).To(HaveLen(1))
		player = *players[0]

		guests, err := keycloak.GetUsers(
			context.Background(),
			clientToken.AccessToken,
			conf.Character.Keycloak.Realm,
			gocloak.GetUsersParams{Username: guest.Username},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(guests).To(HaveLen(1))
		guest = *guests[0]

		md := metadata.New(
			map[string]string{
				"authorization": "Bearer " + adminToken.AccessToken,
			},
		)
		incAdminCtx = metadata.NewIncomingContext(context.Background(), md)
		md = metadata.New(
			map[string]string{
				"authorization": "Bearer " + playerToken.AccessToken,
			},
		)
		incPlayerCtx = metadata.NewIncomingContext(context.Background(), md)
		md = metadata.New(
			map[string]string{
				"authorization": "Bearer " + clientToken.AccessToken,
			},
		)
		incClientCtx = metadata.NewIncomingContext(context.Background(), md)
		md = metadata.New(
			map[string]string{
				"authorization": "Bearer " + guestToken.AccessToken,
			},
		)
		incGuestCtx = metadata.NewIncomingContext(context.Background(), md)
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		if keycloakCloseFunc != nil {
			keycloakCloseFunc()
		}
		if kafkaCloseFunc != nil {
			kafkaCloseFunc()
		}
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Srv Suite")
}

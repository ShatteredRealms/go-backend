package srv_test

import (
	"testing"

	"github.com/Nerzal/gocloak/v13"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	keycloak *gocloak.GoCloak
)

func TestSrv(t *testing.T) {
	var closeFunc func()
	BeforeSuite(func() {
		closeFunc, keycloak = testdb.SetupKeycloakWithDocker()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Srv Suite")

	AfterSuite(func() {
		closeFunc()
	})
}

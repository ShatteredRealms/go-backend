package repository_test

import (
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}

var (
	DB            *gorm.DB
	cleanupDocker func()
)

var _ = BeforeSuite(func() {
	DB, cleanupDocker = helpers.SetupGormWithDocker()
})

var _ = AfterSuite(func() {
	cleanupDocker()
})

var _ = BeforeEach(func() {
	// clear db tables before each test
	Ω(DB.Exec(`DROP SCHEMA public CASCADE;CREATE SCHEMA public;`).Error).To(Succeed())
})

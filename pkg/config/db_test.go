package config_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/config"
)

var _ = Describe("Db config", func() {
	var (
		config *config.DBConfig
		dsn    string
	)

	BeforeEach(func() {
		Expect(faker.FakeData(&config)).To(Succeed())
		dsn = ""
	})

	Describe("MySQLDSN", func() {
		It("should work", func() {
			dsn = config.MySQLDSN()
			Expect(dsn).To(ContainSubstring(config.Name))
		})
	})

	Describe("PostgresDSN", func() {
		It("should work", func() {
			dsn = config.PostgresDSN()
			Expect(dsn).To(ContainSubstring(config.Name))
		})
	})

	Describe("MongoDSN", func() {
		It("should work", func() {
			dsn = config.MongoDSN()
		})
	})

	AfterEach(func() {
		Expect(dsn).To(ContainSubstring(config.Username))
		Expect(dsn).To(ContainSubstring(config.Password))
		Expect(dsn).To(ContainSubstring(config.ServerAddress.Host))
		Expect(dsn).To(ContainSubstring(config.ServerAddress.Port))
	})
})

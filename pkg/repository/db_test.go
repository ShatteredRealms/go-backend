package repository_test

import (
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
)

var _ = Describe("Db repository", func() {
	var (
		pool config.DBPoolConfig
	)
	BeforeEach(func() {
		log.Logger, _ = test.NewNullLogger()
		pool = config.DBPoolConfig{
			Master: data.GormConfig,
			Slaves: []config.DBConfig{data.GormConfig},
		}
	})
	Describe("ConnectDb", func() {
		When("given valid input", func() {
			It("should work", func() {
				out, err := repository.ConnectDB(pool, data.RedisConfig)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error", func() {
				pool.Master.Host = "a"
				out, err := repository.ConnectDB(pool, data.RedisConfig)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})
})

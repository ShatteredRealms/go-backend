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
			Master: config.DBConfig{
				Host:     "127.0.0.1",
				Port:     "5432",
				Name:     "postgres",
				Username: "postgres",
				Password: "password",
			},
			Slaves: []config.DBConfig{{
				Host:     "127.0.0.1",
				Port:     "5432",
				Name:     "postgres",
				Username: "postgres",
				Password: "password",
			}},
		}
	})
	Describe("ConnectDb", func() {
		When("given valid input", func() {
			It("should work", func() {
				out, err := repository.ConnectDB(pool)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error", func() {
				pool.Master.Host = "a"
				out, err := repository.ConnectDB(pool)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})
})

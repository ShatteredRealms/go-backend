package repository_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
)

var _ = Describe("Kafka repository", func() {
	Describe("ConnectKafka", func() {
		var cleanupFunc func()
		Context("valid input", func() {
			It("should work", func() {
				cleanupFunc = testdb.SetupKafkaWithDocker()
				kafka, err := repository.ConnectKafka(config.ServerAddress{
					Port: 29092,
					Host: "localhost",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(kafka).NotTo(BeNil())
			})
		})

		AfterEach(func() {
			if cleanupFunc != nil {
				cleanupFunc()
			}
		})
	})
})

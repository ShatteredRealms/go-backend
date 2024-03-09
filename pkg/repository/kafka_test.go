package repository_test

import (
	. "github.com/onsi/ginkgo/v2"
	// . "github.com/onsi/gomega"
	//
	// "github.com/ShatteredRealms/go-backend/pkg/config"
	// "github.com/ShatteredRealms/go-backend/pkg/repository"
	// testdb "github.com/ShatteredRealms/go-backend/test/db"
)

var _ = Describe("Kafka repository", func() {
	Describe("ConnectKafka", func() {
		Context("valid input", func() {
			It("should work", func() {
				// cleanupFunc, port := testdb.SetupKafkaWithDocker()
				// defer cleanupFunc()
				// kafka, err := repository.ConnectKafka(config.ServerAddress{
				// 	Port: port,
				// 	Host: "localhost",
				// })
				// Expect(err).NotTo(HaveOccurred())
				// Expect(kafka).NotTo(BeNil())
			})
		})
	})
})

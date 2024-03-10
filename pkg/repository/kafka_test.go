package repository_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/segmentio/kafka-go"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
)

var _ = Describe("Kafka repository", func() {
	Describe("ConnectKafka", func() {
		Context("valid input", func() {
			It("should work", func() {
				cleanupFunc, port := testdb.SetupKafkaWithDocker()
				var kafka *kafka.Conn
				var err error
				defer cleanupFunc()
				Eventually(func(g Gomega) error {
					kafka, err = repository.ConnectKafka(config.ServerAddress{
						Port: port,
						Host: "localhost",
					})
					return err
				}).Within(time.Minute).Should(Succeed())
				Expect(kafka).NotTo(BeNil())
			})
		})
	})
})

package repository_test

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
)

var _ = Describe("Kafka repository", func() {
	Describe("ConnectKafka", func() {
		Context("valid input", func() {
			It("should work", func() {
				cleanupFunc, _, kResource := testdb.SetupKafkaWithDocker()
				sPort := kResource.GetPort("29092/tcp")
				port, err := strconv.ParseUint(sPort, 10, 64)
				Expect(err).NotTo(HaveOccurred())
				kafka, err := repository.ConnectKafka(config.ServerAddress{
					Port: uint(port),
					Host: "127.0.0.1",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(kafka).NotTo(BeNil())
				cleanupFunc()
			})
		})
	})
})

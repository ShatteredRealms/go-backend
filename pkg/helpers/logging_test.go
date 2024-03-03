package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
)

var _ = Describe("Logging helpers", func() {
	Describe("SetupLogger", func() {
		It("Should setup loging", func() {
			helpers.SetupLogger()
			Expect(log.GetLevel()).To(Equal(log.TraceLevel))
		})
	})
})

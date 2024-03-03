package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
)

var _ = Describe("Srv helpers", func() {

	log.StandardLogger().ExitFunc = func(int) { wasFatal = true }

	BeforeEach(func() {
		wasFatal = false
	})

	Describe("UnaryLogRequest", func() {
		It("should handle requests", func() {
			helpers.UnaryLogRequest()
			Expect(wasFatal).To(BeFalse())
		})
	})

	Describe("StreamLogRequest", func() {
		It("should handle requests", func() {
			helpers.StreamLogRequest()
			Expect(wasFatal).To(BeFalse())
		})
	})
})

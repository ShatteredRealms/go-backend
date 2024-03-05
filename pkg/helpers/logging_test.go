package helpers_test

import (
	"github.com/ShatteredRealms/go-backend/pkg/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
)

var _ = Describe("Logging helpers", func() {
	Describe("SetupLogger", func() {
		It("Should setup loging", func() {
			helpers.SetupLogger("tests")
			Expect(log.Logger.GetLevel()).To(Equal(logrus.TraceLevel))
		})
	})
})

package helpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var (
	wasFatal bool
)

func TestHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Helpers Suite")

	BeforeEach(func() {
		log.StandardLogger().ExitFunc = func(int) { wasFatal = true }
		wasFatal = false
	})
}

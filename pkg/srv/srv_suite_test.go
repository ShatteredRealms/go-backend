package srv_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSrv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Srv Suite")
}

package gamebackend_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGamebackend(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gamebackend Suite")
}

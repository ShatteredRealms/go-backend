package character_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCharacter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Character Suite")
}

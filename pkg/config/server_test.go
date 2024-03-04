package config_test

import (
	"context"
	"fmt"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/config"
)

var _ = Describe("Server config", func() {
	Describe("ServerAddress Address", func() {
		It("should generate an address", func() {
			config := config.NewGlobalConfig(context.TODO())
			Expect(config).NotTo(BeNil())
			Expect(config.Character.Local.Address()).To(Equal(fmt.Sprintf("%s:%d", config.Character.Local.Host, config.Character.Local.Port)))
			config.Character.Local.Host = faker.Username()
			Expect(config.Character.Local.Address()).To(Equal(fmt.Sprintf("%s:%d", config.Character.Local.Host, config.Character.Local.Port)))
		})
	})
})

package config_test

import (
	"context"
	"os"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	"github.com/ShatteredRealms/go-backend/pkg/config"
)

var _ = Describe("Global config", func() {
	Describe("NewGlobalConfig", func() {
		It("should generate a valid config", func() {
			Expect(config.NewGlobalConfig(nil)).NotTo(BeNil())
			Expect(config.NewGlobalConfig(context.Background())).NotTo(BeNil())
		})

		It("should read in env variables", func() {
			key := "SRO_CHARACTER_LOCAL_HOST"
			val := faker.Username()
			GinkgoT().Setenv(key, val)
			Expect(os.Getenv(key)).To(Equal(val))
			conf := config.NewGlobalConfig(nil)
			Expect(conf).NotTo(BeNil())
			Expect(conf.Character.Local.Host).To(Equal(val))
		})

		It("should read in from config file", func() {
			confFromFile := &config.GlobalConfig{}
			yamlFile, err := os.ReadFile("../../test/config.yaml")
			err = yaml.Unmarshal(yamlFile, confFromFile)
			Expect(err).To(BeNil())
			conf := config.NewGlobalConfig(nil)
			Expect(conf).NotTo(BeNil())
			Expect(conf.Character.Keycloak.ClientSecret).To(Equal(confFromFile.Character.Keycloak.ClientSecret))
		})
	})
})

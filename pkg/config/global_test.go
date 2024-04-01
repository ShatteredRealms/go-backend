package config_test

import (
	"context"
	"os"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"gopkg.in/yaml.v2"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
)

var _ = Describe("Global config", func() {
	BeforeEach(func() {
		log.Logger, _ = test.NewNullLogger()
	})
	Describe("NewGlobalConfig", func() {
		It("should generate a valid config", func() {
			Expect(config.NewGlobalConfig(context.TODO())).NotTo(BeNil())
			Expect(config.NewGlobalConfig(context.Background())).NotTo(BeNil())
		})

		It("should read in env variables", func() {
			key := "SRO_VERSION"
			val := faker.Username()
			GinkgoT().Setenv(key, val)
			Expect(os.Getenv(key)).To(Equal(val))
			conf := config.NewGlobalConfig(context.TODO())
			Expect(conf).NotTo(BeNil())
			Expect(conf.Version).To(Equal(val))

			key = "SRO_CHARACTER_LOCAL_HOST"
			val = faker.Username()
			GinkgoT().Setenv(key, val)
			Expect(os.Getenv(key)).To(Equal(val))
			conf = config.NewGlobalConfig(context.TODO())
			Expect(conf).NotTo(BeNil())
			Expect(conf.Character.Local.Host).To(Equal(val))
		})

		It("should read in from config file", func() {
			confFromFile := &config.GlobalConfig{}
			yamlFile, err := os.ReadFile("../../test/config.yaml")
			Expect(err).To(BeNil())
			err = yaml.Unmarshal(yamlFile, confFromFile)
			Expect(err).To(BeNil())
			conf := config.NewGlobalConfig(context.TODO())
			Expect(conf).NotTo(BeNil())
			Expect(conf.Character.Keycloak.ClientSecret).To(Equal(confFromFile.Character.Keycloak.ClientSecret))
		})
	})
})

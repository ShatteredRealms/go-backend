package config_test

import (
	"context"
	"os"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
)

const (
	envKey         = "SRO_FOO"
	nestedEnvKey   = "SRO_BAR_BAZ"
	embeddedEnvKey = "SRO_BUZ"
)

type TestStruct struct {
	EmbeddedStruct `yaml:",inline" mapstructure:",squash"`

	Foo string
	Bar TestInnerStruct
}

type EmbeddedStruct struct {
	Buz string
}

type TestInnerStruct struct {
	Baz string
}

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

		// It("should read in from config file", func() {
		// 	confFromFile := &config.GlobalConfig{}
		// 	yamlFile, err := os.ReadFile("../../test/config.yaml")
		// 	Expect(err).To(BeNil())
		// 	yaml.Unmarshal(yamlFile, confFromFile)
		// 	conf := config.NewGlobalConfig(context.TODO())
		// 	Expect(conf).NotTo(BeNil())
		// 	Expect(conf.Character.Keycloak.ClientSecret).To(Equal(confFromFile.Character.Keycloak.ClientSecret))
		// })
	})
	Describe("Config helpers", func() {
		var (
			testStruct  *TestStruct
			originalVal string
			envVal      string
		)

		BeforeEach(func() {
			testStruct = &TestStruct{
				Foo: faker.Name(),
				Bar: TestInnerStruct{
					Baz: faker.Name(),
				},
			}
		})

		Describe("BindEnvsToStruct", func() {
			Context("nested structs", func() {
				BeforeEach(func() {
					originalVal = testStruct.Bar.Baz
					envVal = originalVal + faker.Username()
					GinkgoT().Setenv(nestedEnvKey, envVal)
					Expect(os.Getenv(nestedEnvKey)).To(Equal(envVal))
				})

				It("should bind to pointers to structs", func() {
					config.BindEnvsToStruct(testStruct)
				})

				It("should bind to structs", func() {
					config.BindEnvsToStruct(*testStruct)
				})

				AfterEach(func() {
					Expect(originalVal).To(Equal(testStruct.Bar.Baz))
					Expect(viper.Unmarshal(&testStruct)).Should(Succeed())
					Expect(testStruct.Bar.Baz).To(Equal(envVal))
				})
			})

			Context("non-nested structs", func() {
				BeforeEach(func() {
					originalVal = testStruct.Foo
					envVal = faker.Name()
					GinkgoT().Setenv(envKey, envVal)
					Expect(os.Getenv(envKey)).To(Equal(envVal))
				})

				It("should bind to pointers to structs", func() {
					config.BindEnvsToStruct(testStruct)
				})

				It("should bind to structs", func() {
					config.BindEnvsToStruct(*testStruct)
				})

				AfterEach(func() {
					Expect(originalVal).To(Equal(testStruct.Foo))
					Expect(viper.Unmarshal(&testStruct)).Should(Succeed())
					Expect(testStruct.Foo).To(Equal(envVal))
				})
			})

			Context("embedded structs", func() {
				BeforeEach(func() {
					originalVal = testStruct.Buz
					envVal = faker.Name()
					GinkgoT().Setenv(embeddedEnvKey, envVal)
					Expect(os.Getenv(embeddedEnvKey)).To(Equal(envVal))
				})

				It("should bind to pointers to structs", func() {
					config.BindEnvsToStruct(testStruct)
				})

				It("should bind to structs", func() {
					config.BindEnvsToStruct(*testStruct)
				})

				AfterEach(func() {
					Expect(testStruct.Buz).To(Equal(originalVal))
					Expect(viper.Unmarshal(&testStruct)).Should(Succeed())
					Expect(testStruct.Buz).To(Equal(envVal))
				})
			})
		})
	})
})

package helpers_test

import (
	"os"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
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

var _ = Describe("Config helpers", func() {
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
				helpers.BindEnvsToStruct(testStruct)
			})

			It("should bind to structs", func() {
				helpers.BindEnvsToStruct(*testStruct)
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
				helpers.BindEnvsToStruct(testStruct)
			})

			It("should bind to structs", func() {
				helpers.BindEnvsToStruct(*testStruct)
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
				helpers.BindEnvsToStruct(testStruct)
			})

			It("should bind to structs", func() {
				helpers.BindEnvsToStruct(*testStruct)
			})

			AfterEach(func() {
				Expect(testStruct.Buz).To(Equal(originalVal))
				Expect(viper.Unmarshal(&testStruct)).Should(Succeed())
				Expect(testStruct.Buz).To(Equal(envVal))
			})
		})
	})
})

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
	envKey       = "SRO_FOO"
	nestedEnvKey = "SRO_BAR_BAZ"
)

type TestStruct struct {
	Foo string
	Bar TestInnerStruct
}

type TestInnerStruct struct {
	Baz string
}

var _ = Describe("Config helpers", func() {
	var (
		testStruct  *TestStruct
		originalBaz string
		envVal      string
	)

	BeforeEach(func() {
		viper.SetEnvPrefix("SRO")
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
				originalBaz = testStruct.Bar.Baz
				envVal = faker.Name()
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
				Expect(originalBaz).To(Equal(testStruct.Bar.Baz))
				Expect(viper.Unmarshal(&testStruct)).Should(Succeed())
				Expect(testStruct.Bar.Baz).To(Equal(envVal))
			})
		})

		Context("non-nested structs", func() {
			BeforeEach(func() {
				originalBaz = testStruct.Foo
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
				Expect(originalBaz).To(Equal(testStruct.Foo))
				Expect(viper.Unmarshal(&testStruct)).Should(Succeed())
				Expect(testStruct.Foo).To(Equal(envVal))
			})
		})
	})
})

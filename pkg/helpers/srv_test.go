package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
)

var _ = Describe("Srv helpers", func() {
	var (
		hook *test.Hook
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
		hook.Reset()
	})

	Describe("GrpcDialOpts", func() {
		It("should create DialOptions", func() {
			Expect(helpers.GrpcDialOpts()).NotTo(BeEmpty())
		})
	})

	Describe("GrpcClientWithOtel", func() {
		It("should dial address without error", func() {
			client, err := helpers.GrpcClientWithOtel("127.0.0.1:9999")
			Expect(client).NotTo(BeNil())
			Expect(err).To(Succeed())
		})
	})

	Describe("InitServerDefaults", func() {
		It("should create default server and mux", func() {
			server, mux := helpers.InitServerDefaults(nil, "")
			Expect(server).NotTo(BeNil())
			Expect(mux).NotTo(BeNil())
		})
	})
})

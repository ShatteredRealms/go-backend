package helpers_test

import (
	"context"
	"encoding/base64"

	"github.com/bxcodec/faker/v4"
	"github.com/golang-jwt/jwt/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/model"
)

var _ = Describe("Srv helpers", func() {

	log.StandardLogger().ExitFunc = func(int) { wasFatal = true }

	BeforeEach(func() {
		wasFatal = false
	})

	Describe("UnaryLogRequest", func() {
		It("should handle requests", func() {
			helpers.UnaryLogRequest()
			Expect(wasFatal).To(BeFalse())
		})
	})

	Describe("StreamLogRequest", func() {
		It("should handle requests", func() {
			helpers.StreamLogRequest()
			Expect(wasFatal).To(BeFalse())
		})
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
			server, mux := helpers.InitServerDefaults()
			Expect(server).NotTo(BeNil())
			Expect(mux).NotTo(BeNil())
		})
	})

	Describe("ExtractToken", func() {
		Context("invalid input", func() {
			It("should err on nil ctx", func() {
				token, err := helpers.ExtractToken(nil)
				Expect(token).To(BeEmpty())
				Expect(err).NotTo(Succeed())
			})

			It("should err on ctx with no auth", func() {
				token, err := helpers.ExtractToken(context.Background())
				Expect(token).To(BeEmpty())
				Expect(err).NotTo(Succeed())
			})

			It("should err on incorrect auth scheme", func() {
				md := metadata.New(
					map[string]string{
						"authorization": faker.Name(),
					},
				)
				token, err := helpers.ExtractToken(metadata.NewIncomingContext(context.Background(), md))
				Expect(token).To(BeEmpty())
				Expect(err).NotTo(Succeed())
			})
		})
		Context("correct input", func() {
			It("should extract auth token", func() {
				inToken := faker.Name()
				token, err := helpers.ExtractToken(createIncomingAuthToken(inToken))
				Expect(token).To(Equal(inToken))
				Expect(err).To(Succeed())
			})
		})
	})

	Describe("ExtractClaims", func() {
		Context("invalid input", func() {
			It("should error on nil context", func() {
				claims, err := helpers.ExtractClaims(nil)
				Expect(claims).To(BeNil())
				Expect(err).NotTo(BeNil())
			})

			It("should error on context with no auth token", func() {
				claims, err := helpers.ExtractClaims(context.Background())
				Expect(claims).To(BeNil())
				Expect(err).NotTo(BeNil())
			})

			It("should error on garbage auth token", func() {
				claims, err := helpers.ExtractClaims(createIncomingAuthToken("token"))
				Expect(claims).To(BeNil())
				Expect(err).NotTo(BeNil())
			})

			It("should error on invalid auth token", func() {
				claims, err := helpers.ExtractClaims(createIncomingAuthToken(faker.Jwt()))
				Expect(claims).To(BeNil())
				Expect(err).NotTo(BeNil())
			})
		})

		Context("valid input", func() {
			It("should extract claims", func() {
				var claims *model.SROClaims
				faker.FakeData(&claims)

				bytes, err := base64.StdEncoding.DecodeString("gEQCe2i8oOWWmyerVdL3KZik4FdyGUGGls/dIewSkVo=")
				Expect(err).To(BeNil())

				token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(bytes)
				Expect(err).To(BeNil())

				outClaims, err := helpers.ExtractClaims(createIncomingAuthToken(token))
				Expect(err).To(BeNil())
				Expect(outClaims.Issuer).To(Equal(claims.Issuer))
			})
		})
	})

	Describe("VerifyClaims", func() {
		Context("invalid input", func() {
		})

		Context("valid input", func() {
		})
	})
})

func createIncomingAuthToken(jwt string) context.Context {
	md := metadata.New(
		map[string]string{
			"authorization": "Bearer " + jwt,
		},
	)
	return metadata.NewIncomingContext(context.Background(), md)
}

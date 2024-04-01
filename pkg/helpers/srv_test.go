package helpers_test

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/bxcodec/faker/v4"
	"github.com/golang-jwt/jwt/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
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
		var (
			mockCtrl     *gomock.Controller
			mockKeycloak *mocks.MockKeycloakClient
		)
		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			mockKeycloak = mocks.NewMockKeycloakClient(mockCtrl)
		})
		Context("invalid input", func() {
			It("should error on nil inputs", func() {
				jwtToken, outClaims, err := helpers.VerifyClaims(nil, nil, faker.Name())
				Expect(jwtToken).To(BeNil())
				Expect(outClaims).To(BeNil())
				Expect(err).NotTo(BeNil())

				jwtToken, outClaims, err = helpers.VerifyClaims(context.Background(), nil, faker.Name())
				Expect(jwtToken).To(BeNil())
				Expect(outClaims).To(BeNil())
				Expect(err).NotTo(BeNil())

				jwtToken, outClaims, err = helpers.VerifyClaims(nil, &gocloak.GoCloak{}, faker.Name())
				Expect(jwtToken).To(BeNil())
				Expect(outClaims).To(BeNil())
				Expect(err).NotTo(BeNil())
			})

			It("should require ctx with valid incoming token", func() {
				jwtToken, outClaims, err := helpers.VerifyClaims(context.Background(), &gocloak.GoCloak{}, "")
				Expect(jwtToken).To(BeNil())
				Expect(outClaims).To(BeNil())
				Expect(err).NotTo(BeNil())
			})

			It("should error on keycloak error", func() {
				claims := model.SROClaims{
					Username: faker.Username(),
				}
				bytes, err := base64.StdEncoding.DecodeString("gEQCe2i8oOWWmyerVdL3KZik4FdyGUGGls/dIewSkVo=")
				Expect(err).To(BeNil())

				token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(bytes)
				Expect(err).To(BeNil())

				ctx := createIncomingAuthToken(token)
				realm := faker.Name()
				err = fmt.Errorf(faker.Username())

				mockKeycloak.EXPECT().DecodeAccessTokenCustomClaims(gomock.Eq(ctx), gomock.Any(), gomock.Eq(realm), gomock.Any()).Return(nil, err)

				jwtToken, outClaims, err := helpers.VerifyClaims(ctx, mockKeycloak, realm)
				Expect(jwtToken).To(BeNil())
				Expect(outClaims).To(BeNil())
				Expect(err).To(MatchError(err))
			})

			It("should error on invalid token", func() {
				claims := model.SROClaims{
					Username: faker.Username(),
				}
				bytes, err := base64.StdEncoding.DecodeString("gEQCe2i8oOWWmyerVdL3KZik4FdyGUGGls/dIewSkVo=")
				Expect(err).To(BeNil())

				originalToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				token, err := originalToken.SignedString(bytes)
				Expect(err).To(BeNil())

				ctx := createIncomingAuthToken(token)
				realm := faker.Name()
				err = fmt.Errorf(faker.Username())

				mockKeycloak.EXPECT().
					DecodeAccessTokenCustomClaims(
						gomock.Eq(ctx),
						gomock.Eq(token),
						gomock.Eq(realm),
						gomock.Any(),
					).Return(originalToken, nil)

				jwtToken, outClaims, err := helpers.VerifyClaims(ctx, mockKeycloak, realm)
				Expect(jwtToken).To(BeNil())
				Expect(outClaims).To(BeNil())
				Expect(err).To(MatchError(model.ErrUnauthorized.Err()))
			})
		})

		Context("valid input", func() {
			It("should work", func() {
				claims := model.SROClaims{
					Username: faker.Username(),
				}
				bytes, err := base64.StdEncoding.DecodeString("gEQCe2i8oOWWmyerVdL3KZik4FdyGUGGls/dIewSkVo=")
				Expect(err).To(BeNil())

				originalToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				originalToken.Valid = true
				token, err := originalToken.SignedString(bytes)
				Expect(err).To(BeNil())

				ctx := createIncomingAuthToken(token)
				realm := faker.Name()
				err = fmt.Errorf(faker.Username())

				mockKeycloak.EXPECT().
					DecodeAccessTokenCustomClaims(
						gomock.Eq(ctx),
						gomock.Eq(token),
						gomock.Eq(realm),
						gomock.Any(),
					).Return(originalToken, nil)

				// Note: claims are set by pointer pass to DecodeAccessTokenCustomClaims so testing is irrelevant
				jwtToken, _, err := helpers.VerifyClaims(ctx, mockKeycloak, realm)
				Expect(jwtToken).NotTo(BeNil(), "valid token")
				Expect(*jwtToken).To(Equal(*originalToken), "token should resolve correctly")
				Expect(err).To(BeNil(), "should not error")
			})
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

type mockHandler struct {
	WasCalled bool
}

func (m *mockHandler) StreamHandler(srv interface{}, stream grpc.ServerStream) error {
	m.WasCalled = true
	return nil
}

func (m *mockHandler) UnaryHandler(ctx context.Context, req interface{}) (interface{}, error) {
	m.WasCalled = true
	return nil, nil
}

type mockServerStream struct{}

func (m mockServerStream) Context() context.Context {
	return context.Background()
}

func (m mockServerStream) SetHeader(metadata.MD) error {
	return nil
}

func (m mockServerStream) SendHeader(metadata.MD) error {
	return nil
}

func (m mockServerStream) SetTrailer(metadata.MD) {

}

func (m mockServerStream) SendMsg(i interface{}) error {
	return nil
}

func (m mockServerStream) RecvMsg(i interface{}) error {
	return nil
}

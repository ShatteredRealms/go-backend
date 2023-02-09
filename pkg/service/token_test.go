package service_test

import (
	"encoding/json"
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"io/fs"
	"reflect"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt"
)

var _ = Describe("Token", func() {

	keyDir := "../../test/auth"

	var jwtService service.JWTService
	var ttl time.Duration
	var expectedError error
	var claims jwt.MapClaims

	BeforeEach(func() {
		expectedError = nil
	})

	Context("Service creation", func() {
		var dir string

		It("should fail if there is no private key", func() {
			dir = keyDir + "/invalid/nopriv"
			expectedError = &fs.PathError{
				Op:   "open",
				Path: dir + "/key",
				Err:  syscall.Errno(0x2),
			}
			expectedError = fmt.Errorf("jwt: private key file: %w", expectedError)
		})

		It("should fail if there is no public key", func() {
			dir = keyDir + "/invalid/nopub"
			expectedError = &fs.PathError{
				Op:   "open",
				Path: dir + "/key.pub",
				Err:  syscall.Errno(0x2),
			}
			expectedError = fmt.Errorf("jwt: public key file: %w", expectedError)
		})

		It("should succeed both keys are present", func() {
			dir = keyDir
		})

		AfterEach(func() {
			jwtService, err := service.NewJWTService(dir)
			if expectedError == nil {
				Expect(err).To(BeNil())
				Expect(jwtService).NotTo(BeNil())
			} else {
				Expect(err).To(Equal(expectedError))
				Expect(jwtService).To(BeNil())
			}
		})
	})

	Context("Existing JWT service", func() {
		BeforeEach(func() {
			claims = make(jwt.MapClaims)
			ttl = time.Hour
			jwtService, _ = service.NewJWTService(keyDir)
		})

		Context("Creation", func() {
			It("should fail with an invalid jwt private key", func() {
				jwtService, _ = service.NewJWTService(keyDir + "/invalid")
				expectedError = fmt.Errorf("Invalid Key: Key must be a PEM encoded PKCS1 or PKCS8 key")
				expectedError = fmt.Errorf("create: parse key: %w", expectedError)
			})

			It("should fail with duplicate claims", func() {
				key := "iss"
				claims[key] = helpers.RandString(10)
				expectedError = fmt.Errorf("claim value already set: %s", key)
			})

			It("should fail with invalid claims", func() {
				key := helpers.RandString(4)
				claims[key] = make(chan int)
				expectedError = &json.UnsupportedTypeError{
					Type: reflect.TypeOf(claims[key]),
				}
				expectedError = fmt.Errorf("create: sign token: %w", expectedError)
			})

			It("should succeed if claims and service are valid", func() {
				key := helpers.RandString(4)
				claims[key] = helpers.RandString(10)
			})

			AfterEach(func() {
				token, err := jwtService.Create(ttl, "test.domain.tld", claims)
				if expectedError == nil {
					Expect(err).To(BeNil())
					Expect(token).NotTo(BeEmpty())
				} else {
					Expect(err).To(Equal(expectedError))
					Expect(token).To(BeEmpty())
				}
			})
		})

		Context("Validation", func() {
			var token string
			var claims jwt.MapClaims

			BeforeEach(func() {
				var err error
				jwtService, err = service.NewJWTService(keyDir)
				Expect(err).To(BeNil())

				claims = make(jwt.MapClaims)
				claims[helpers.RandString(4)] = helpers.RandString(10)
				token, _ = jwtService.Create(time.Hour, "test.domain.tld", nil)
				Expect(err).To(BeNil())
			})

			It("should fail with an invalid jwt public key", func() {
				jwtService, _ = service.NewJWTService(keyDir + "/invalid")
				expectedError = fmt.Errorf("Invalid Key: Key must be a PEM encoded PKCS1 or PKCS8 key")
				expectedError = fmt.Errorf("validates: parse key: %w", expectedError)
			})

			It("should success with valid jwt service and token", func() {
			})

			AfterEach(func() {
				claims, err := jwtService.Validate(token)
				if expectedError == nil {
					Expect(err).To(BeNil())
					Expect(claims).NotTo(BeNil())
				} else {
					Expect(err).To(Equal(expectedError))
					Expect(claims).To(BeNil())
				}
			})
		})
	})
})

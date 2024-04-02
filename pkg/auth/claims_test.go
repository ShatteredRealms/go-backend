package auth_test

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/ShatteredRealms/go-backend/pkg/auth"
	"github.com/bxcodec/faker/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"
)

var _ = Describe("Claims auth", func() {
	Context("contexts", func() {
		var ctx context.Context
		var randVal, randWord string
		type ctxKey int
		var key ctxKey

		BeforeEach(func() {
			key = 0
			randVal = faker.Word()
			randWord = faker.Word()
			ctx = context.WithValue(context.Background(), key, randVal)
		})

		Describe("ContextAddClientToken", func() {
			It("should add bearer token to outgoing context", func() {
				initialIncomingAuth := metautils.ExtractIncoming(ctx).Get("authorization")
				newCtx := auth.AddOutgoingToken(ctx, randWord)
				Expect(metautils.ExtractOutgoing(newCtx).Get("authorization")).To(Equal("Bearer "+randWord), "bearer token missing")
				Expect(metautils.ExtractIncoming(newCtx).Get("authorization")).To(Equal(initialIncomingAuth), "shouldn't change incoming context")
				Expect(ctx.Value(key)).To(Equal(randVal), "shouldn't lose existing ctx values")
			})
		})

		Describe("PassAuthContext", func() {
			It("should add incoming authorization to outgoing", func() {
				md := metadata.New(
					map[string]string{
						"authorization": randWord,
					},
				)
				ctx = metadata.NewIncomingContext(ctx, md)
				initialIncomingAuth := metautils.ExtractIncoming(ctx).Get("authorization")
				ctx = auth.PassOutgoing(ctx)
				Expect(metautils.ExtractOutgoing(ctx).Get("authorization")).To(Equal(randWord), "authorization missing")
				Expect(metautils.ExtractIncoming(ctx).Get("authorization")).To(Equal(initialIncomingAuth), "shouldn't change incoming context")
				Expect(ctx.Value(key)).To(Equal(randVal), "shouldn't lose existing ctx values")
			})
		})
	})

	Describe("HasRole", func() {
		It("should work", func() {
			role := &gocloak.Role{
				Name: gocloak.StringP(faker.Username()),
			}
			role2 := &gocloak.Role{
				Name: gocloak.StringP(faker.Email()),
			}
			claims := &auth.SROClaims{
				RealmRoles: auth.ClaimRoles{
					Roles: []string{faker.Username(), *role.Name},
				},
			}

			Expect(claims.HasRole(role)).To(BeTrue())
			Expect(claims.HasRole(role2)).To(BeFalse())
		})
	})

	Describe("HasResourceRole", func() {
		It("should work", func() {
			role := &gocloak.Role{
				Name: gocloak.StringP(faker.Username()),
			}

			role2 := &gocloak.Role{
				Name: gocloak.StringP(faker.Email()),
			}

			clientId := faker.Username()
			claims := &auth.SROClaims{
				ResourceAccess: auth.ClaimResourceAccess{
					clientId: auth.ClaimRoles{
						Roles: []string{faker.Username(), *role.Name},
					},
				},
			}

			Expect(claims.HasResourceRole(role, clientId)).To(BeTrue())
			Expect(claims.HasResourceRole(role2, clientId)).To(BeFalse())
			Expect(claims.HasResourceRole(role, clientId+"a")).To(BeFalse())
		})
	})
})

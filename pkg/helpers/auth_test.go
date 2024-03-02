package helpers_test

import (
	"context"

	"github.com/bxcodec/faker/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/metadata"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
)

var _ = Describe("Auth helpers", func() {
	var ctx context.Context
	var randVal, randWord string

	BeforeEach(func() {
		randVal = faker.Word()
		randWord = faker.Word()
		ctx = context.WithValue(context.Background(), "key", randVal)
	})

	Describe("ContextAddClientToken", func() {
		It("should add bearer token to outgoing context", func() {
			initialIncomingAuth := metautils.ExtractIncoming(ctx).Get("authorization")
			newCtx := helpers.ContextAddClientToken(ctx, randWord)
			Expect(metautils.ExtractOutgoing(newCtx).Get("authorization")).To(Equal("Bearer "+randWord), "bearer token missing")
			Expect(metautils.ExtractIncoming(newCtx).Get("authorization")).To(Equal(initialIncomingAuth), "shouldn't change incoming context")
			Expect(ctx.Value("key")).To(Equal(randVal), "shouldn't lose existing ctx values")
		})
	})

	Describe("ContextAddClientBearerToken", func() {
		It("should add authorization metadata to outgoing context", func() {
			initialIncomingAuth := metautils.ExtractIncoming(ctx).Get("authorization")
			ctx = helpers.ContextAddClientBearerToken(ctx, randWord)
			Expect(metautils.ExtractOutgoing(ctx).Get("authorization")).To(Equal(randWord), "authorization missing")
			Expect(metautils.ExtractIncoming(ctx).Get("authorization")).To(Equal(initialIncomingAuth), "shouldn't change incoming context")
			Expect(ctx.Value("key")).To(Equal(randVal), "shouldn't lose existing ctx values")
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
			ctx = helpers.PassAuthContext(ctx)
			Expect(metautils.ExtractOutgoing(ctx).Get("authorization")).To(Equal(randWord), "authorization missing")
			Expect(metautils.ExtractIncoming(ctx).Get("authorization")).To(Equal(initialIncomingAuth), "shouldn't change incoming context")
			Expect(ctx.Value("key")).To(Equal(randVal), "shouldn't lose existing ctx values")
		})
	})
})

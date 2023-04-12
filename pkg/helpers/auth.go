package helpers

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/metadata"
)

const (
	RolesCtxKey = "roles"
)

func ContextAddClientToken(ctx context.Context, token string) context.Context {
	return ContextAddClientBearerToken(ctx, "Bearer "+token)
}

func ContextAddClientBearerToken(ctx context.Context, token string) context.Context {
	md := metadata.New(
		map[string]string{
			"authorization": token,
		},
	)

	return metadata.AppendToOutgoingContext(metadata.NewOutgoingContext(ctx, md))
}

func PassAuthContext(ctx context.Context) context.Context {
	return ContextAddClientBearerToken(
		ctx,
		metautils.ExtractIncoming(ctx).Get("authorization"),
	)
}

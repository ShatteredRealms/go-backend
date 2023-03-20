package helpers

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
)

const (
	RolesCtxKey = "roles"
)

func ContextAddClientToken(ctx context.Context, token string) context.Context {
	md := metadata.New(
		map[string]string{
			"authorization": fmt.Sprintf("Bearer %s", token),
		},
	)
	return metadata.AppendToOutgoingContext(metadata.NewOutgoingContext(ctx, md))
}

package helpers

import (
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/grpc/metadata"
)

const (
	RolesCtxKey = "roles"
)

func ContextAddClientAuth(ctx context.Context, clientId string, clientSecret string) context.Context {
	unencoded := fmt.Sprintf("%s:%s", clientId, clientSecret)
	md := metadata.New(
		map[string]string{
			"authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(unencoded))),
		},
	)
	return metadata.AppendToOutgoingContext(metadata.NewOutgoingContext(ctx, md))
}

package helpers

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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

func UnaryAddRolesToCtx() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		newCtx, err := newCtxWithRoles(ctx)

		if err != nil {
			log.WithContext(ctx).Error(err)
			return handler(ctx, req)
		}

		return handler(newCtx, req)
	}
}

func StreamAddRolesToCtx() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		newCtx, err := newCtxWithRoles(stream.Context())

		if err != nil {
			log.WithContext(stream.Context()).Error(err)
			return handler(srv, stream)
		}

		return handler(srv, &grpc_middleware.WrappedServerStream{
			ServerStream:   stream,
			WrappedContext: newCtx,
		})
	}
}

func newCtxWithRoles(ctx context.Context) (context.Context, error) {
	token, err := ExtractToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("extracting token: %v", err)
	}

	jwtToken, _, err := jwtParser.ParseUnverified(token, model.SROClaims{})
	if err != nil {
		return nil, fmt.Errorf("parsing token: %v", err)
	}

	claims := jwtToken.Claims.(model.SROClaims)

	roles := make(
		map[string]struct{},
		len(claims.CharacterRoles)+
			len(claims.ChatRoles)+
			len(claims.GameBackendRoles)+
			len(claims.RealmRoles),
	)

	for _, r := range claims.CharacterRoles {
		roles["sro-characters:"+r] = struct{}{}
	}

	for _, r := range claims.ChatRoles {
		roles["sro-chat:"+r] = struct{}{}
	}

	for _, r := range claims.GameBackendRoles {
		roles["sro-gamebackend:"+r] = struct{}{}
	}

	for _, r := range claims.GameBackendRoles {
		roles[r] = struct{}{}
	}

	return context.WithValue(ctx, RolesCtxKey, roles), nil
}

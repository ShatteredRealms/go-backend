package auth

import (
	"context"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/srospan"
	"github.com/golang-jwt/jwt/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

const (
	AuthorizationHeader = "authorization"
	AuthorizationScheme = "Bearer "
)

type ResourceRoles struct {
	Roles []string `json:"roles"`
}

type SROClaims struct {
	jwt.RegisteredClaims
	RealmRoles     ClaimRoles          `json:"realm_roles,omitempty"`
	ResourceAccess ClaimResourceAccess `json:"resource_access"`
	Username       string              `json:"preferred_username,omitempty"`
}

type ClaimRoles struct {
	Roles []string `json:"roles,omitempty"`
}
type ClaimResourceAccess map[string]ClaimRoles

func (s SROClaims) HasResourceRole(role *gocloak.Role, clientId string) bool {
	if resource, ok := s.ResourceAccess[clientId]; ok {
		for _, claimRole := range resource.Roles {
			if *role.Name == claimRole {
				return true
			}
		}
	}

	return false
}

func (s SROClaims) HasRole(role *gocloak.Role) bool {
	for _, claimRole := range s.RealmRoles.Roles {
		if *role.Name == claimRole {
			return true
		}
	}

	return false
}

func extractToken(ctx context.Context) (string, error) {
	val := metautils.ExtractIncoming(ctx).Get(AuthorizationHeader)
	if val == "" {
		return "", common.ErrMissingAuthorization
	}

	if !strings.HasPrefix(val, AuthorizationScheme) {
		return "", common.ErrInvalidAuthorization
	}

	return val[len(AuthorizationScheme):], nil
}

func verifyClaims(ctx context.Context, client KeycloakClient, realm string) (*jwt.Token, *SROClaims, error) {
	if client == nil {
		return nil, nil, common.ErrMissingGocloak
	}

	tokenString, err := extractToken(ctx)
	if err != nil {
		return nil, nil, err
	}

	var claims SROClaims
	token, err := client.DecodeAccessTokenCustomClaims(
		ctx,
		tokenString,
		realm,
		&claims,
	)

	if err != nil {
		log.Logger.WithContext(ctx).Infof("Error extracting claims: %v", err)
		return nil, nil, err
	}

	if !token.Valid {
		log.Logger.WithContext(ctx).Infof("Invalid token given from %s:%s", claims.Username, claims.Subject)
		return nil, nil, err
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		srospan.SourceOwnerId(claims.Subject),
		srospan.SourceOwnerUsername(claims.Username),
	)

	return token, &claims, nil
}

func AddOutgoingToken(ctx context.Context, token string) context.Context {
	return addOutgoingAuthBearer(ctx, "Bearer "+token)
}

func PassOutgoing(ctx context.Context) context.Context {
	return addOutgoingAuthBearer(
		ctx,
		metautils.ExtractIncoming(ctx).Get("authorization"),
	)
}

func addOutgoingAuthBearer(ctx context.Context, token string) context.Context {
	md := metadata.New(
		map[string]string{
			"authorization": token,
		},
	)

	return metadata.AppendToOutgoingContext(metadata.NewOutgoingContext(ctx, md))
}

package auth

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/srospan"
	"github.com/WilSimpson/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/trace"
)

func verifyClaims(ctx context.Context, client gocloak.KeycloakClient, realm string) (*jwt.Token, *SROClaims, error) {
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

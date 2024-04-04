package config

import (
	"context"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/auth"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/WilSimpson/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/trace"
)

var (
	parser jwt.Parser
)

type ServerContext struct {
	GlobalConfig   *GlobalConfig
	KeycloakClient *gocloak.GoCloak
	Tracer         trace.Tracer
	RefSROServer   *SROServer

	jwt            *gocloak.JWT
	tokenExpiresAt time.Time
}

func (srvCtx *ServerContext) GetJWT(ctx context.Context) (*gocloak.JWT, error) {
	if srvCtx.jwt != nil && time.Now().Before(srvCtx.tokenExpiresAt) {
		return srvCtx.jwt, nil
	}

	return srvCtx.loginClient(ctx)
}

func (srvCtx *ServerContext) loginClient(ctx context.Context) (*gocloak.JWT, error) {
	var err error
	srvCtx.jwt, err = srvCtx.KeycloakClient.LoginClient(
		ctx,
		srvCtx.RefSROServer.Keycloak.ClientId,
		srvCtx.RefSROServer.Keycloak.ClientSecret,
		srvCtx.GlobalConfig.Keycloak.Realm,
	)
	if err != nil {
		return nil, fmt.Errorf("login keycloak: %v", err)
	}

	claims := &jwt.RegisteredClaims{}
	_, _, err = parser.ParseUnverified(srvCtx.jwt.AccessToken, claims)
	if err != nil {
		log.Logger.Errorf("parsing access token: %v", err)
		return srvCtx.jwt, nil
	}

	// Remove 5 seconds to ensure there are no race cases with expiration
	srvCtx.tokenExpiresAt = claims.ExpiresAt.Time.Add(-5 * time.Second)
	return srvCtx.jwt, nil
}

func (srvCtx *ServerContext) OutgoingClientAuth(ctx context.Context) (context.Context, error) {
	token, err := srvCtx.GetJWT(ctx)
	if err != nil {
		return ctx, err
	}

	return auth.AddOutgoingToken(
		ctx,
		token.AccessToken,
	), nil
}

package auth

import (
	"context"
)

const (
	claimContextKey claimContextKeyType = iota
)

func RetrieveClaims(ctx context.Context) (claims *SROClaims, ok bool) {
	claims, ok = ctx.Value(claimContextKey).(*SROClaims)
	return
}

func insertClaims(ctx context.Context, claims *SROClaims) context.Context {
	return context.WithValue(ctx, claimContextKey, claims)
}

package helpers

import (
	"context"
	"github.com/ShatteredRealms/go-backend/pkg/log"
)

func Check(ctx context.Context, err error, errContext string) {
	if err != nil {
		log.Logger.WithContext(ctx).Fatalf("%s: %v", errContext, err)
	}
}

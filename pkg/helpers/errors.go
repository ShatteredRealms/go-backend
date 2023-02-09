package helpers

import (
    "context"
    log "github.com/sirupsen/logrus"
)

func Check(ctx context.Context, err error, errContext string) {
    if err != nil {
        log.WithContext(ctx).Fatalf("%s: %v", errContext, err)
    }
}

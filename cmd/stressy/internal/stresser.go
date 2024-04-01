package internal

import "context"

type Stresser interface {
	Start(ctx context.Context) error
	Shutdown() error
}

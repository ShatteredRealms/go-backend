package srv

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotAuthorized = status.Error(codes.Unauthenticated, "not authorized")
)

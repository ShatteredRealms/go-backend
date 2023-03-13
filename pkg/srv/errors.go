package srv

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrDoesNotExist = status.Error(codes.InvalidArgument, "does not exist")
)

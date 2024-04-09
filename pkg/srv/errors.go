package srv

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternalCreateCharacter = status.Error(codes.Internal, "unable to create character")
	ErrInternalEditCharacter   = status.Error(codes.Internal, "unable to edit character")
	ErrInvalidDimension        = status.Error(codes.InvalidArgument, "invalid dimension requested")
)

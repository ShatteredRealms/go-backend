package common

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrInvalidName thrown when a name has invalid characters
	ErrInvalidName = errors.New("name contains invalid character(s)")

	// ErrNameProfane thrown when a name is profane
	ErrNameProfane = errors.New("name unavailable")

	// ErrInvalidRealm thrown when a character belongs to an unknown realm
	ErrInvalidRealm = errors.New("invalid realm")

	// ErrInvalidGender thrown when a character belongs to an unknown gender
	ErrInvalidGender = errors.New("invalid gender")

	// ErrInvalidServerLocation thrown when a server location is unknown
	ErrInvalidServerLocation = errors.New("invalid server location")

	// ErrNotOwner thrown when the requester is not the owner of the target
	// and doesn't have additional priviledges
	ErrNotOwner = errors.New("not owner")

	ErrMissingContext = errors.New("missing valid context")

	ErrMissingAuthorization = errors.New("missing authorization")
	ErrInvalidAuthorization = errors.New("invalid authorization scheme")
	ErrInvalidAuth          = errors.New("invalid auth")
	ErrMissingGocloak       = errors.New("error A01")

	ErrUnauthorized  = status.New(codes.Unauthenticated, "not authorized")
	ErrDoesNotExist  = status.New(codes.InvalidArgument, "does not exist")
	ErrHandleRequest = status.New(codes.Internal, "unable to handle request")
)

package model

import "errors"

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
)

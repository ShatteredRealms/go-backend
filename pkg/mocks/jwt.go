package mocks

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

// JWT If true, returns no errors with string ok, otherwise returns error.
type JWT bool

func (t JWT) Create(time.Duration, string, jwt.MapClaims) (string, error) {
	if t {
		return "ok", nil
	} else {
		return "", fmt.Errorf("error")
	}
}

func (t JWT) Validate(token string) (interface{}, error) {
	if t {
		return "ok", nil
	} else {
		return "", fmt.Errorf("error")
	}
}

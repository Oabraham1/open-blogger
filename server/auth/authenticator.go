package auth

import (
	"time"
)

/* Authenticator is an interface for creating and verifying tokens */
type Authenticator interface {
	CreateToken(username string, duration time.Duration) (string, *AuthPayload, error)
	VerifyToken(token string) (*AuthPayload, error)
}

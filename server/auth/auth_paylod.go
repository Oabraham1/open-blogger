package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

/* Errors returned by VeirfyAuthToken */
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

/* AuthPayload is the payload of a token */
type AuthPayload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

/* NewAuthPayload creates a new token payload */
func NewAuthPayload(username string, duration time.Duration) (*AuthPayload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &AuthPayload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

/* Valid checks if the token is valid */
func (payload *AuthPayload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

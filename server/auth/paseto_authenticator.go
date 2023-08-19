package auth

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

/* PasetoAuthenticator is a pasto token */
type PasetoAuthenticator struct {
	Paseto       *paseto.V2
	SymmetricKey []byte
}

/* NewPasetoAuthenticator creates a new paseto authenticator */
func NewPasetoAuthenticator(symmetricKey string) (Authenticator, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	authenticator := &PasetoAuthenticator{
		Paseto:       paseto.NewV2(),
		SymmetricKey: []byte(symmetricKey),
	}
	return authenticator, nil
}

/* CreateToken creates a new token */
func (p *PasetoAuthenticator) CreateToken(username string, duration time.Duration) (string, *AuthPayload, error) {
	payload, err := NewAuthPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := p.Paseto.Encrypt(p.SymmetricKey, payload, nil)
	return token, payload, err
}

/* VerifyToken verifies a token */
func (p *PasetoAuthenticator) VerifyToken(token string) (*AuthPayload, error) {
	payload := &AuthPayload{}
	err := p.Paseto.Decrypt(token, p.SymmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

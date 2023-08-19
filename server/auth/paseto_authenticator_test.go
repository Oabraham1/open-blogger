package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoAuthenticatorCreateToken(t *testing.T) {
	pasetoAuthenticator, err := NewPasetoAuthenticator("01234567890123456789012345678901")
	require.NoError(t, err)

	username := "testUser"
	duration := time.Minute * 15

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := pasetoAuthenticator.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPasetoAuthenticatorVerifyToken(t *testing.T) {
	pasetoAuthenticator, err := NewPasetoAuthenticator("01234567890123456789012345678901")
	require.NoError(t, err)

	username := "testUser"
	duration := time.Minute * 15

	token, _, err := pasetoAuthenticator.CreateToken(username, duration)
	require.NoError(t, err)

	payload, err := pasetoAuthenticator.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
}

func TestPasetoAuthenticatorExpiredToken(t *testing.T) {
	pasetoAuthenticator, err := NewPasetoAuthenticator("01234567890123456789012345678901")
	require.NoError(t, err)

	username := "testUser"
	duration := -time.Minute

	token, payload, err := pasetoAuthenticator.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = pasetoAuthenticator.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

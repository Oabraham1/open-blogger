package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateUser() *User {
	return &User{
		ID:        "123",
		Username:  "testUser",
		FirstName: "test",
		LastName:  "user",
		Email:     "testuser@email.com",
		ImageURL:  "http://testuser.com/image",
	}
}

func TestNewUser(t *testing.T) {
	user := CreateUser()

	require.NotNil(t, user)
	require.Greater(t, len(user.ID), 0)
	require.Greater(t, len(user.Username), 0)
	require.Greater(t, len(user.FirstName), 0)
	require.Greater(t, len(user.LastName), 0)
	require.Greater(t, len(user.Email), 0)
	require.Greater(t, len(user.ImageURL), 0)
}

func TestValidate(t *testing.T) {
	user := CreateUser()
	err := user.Validate()
	require.Nil(t, err)
}

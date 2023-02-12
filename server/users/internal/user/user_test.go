package user

import (
	"testing"
	"time"

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
	require.Greater(t, time.Now(), user.Created)
}

func TestUpdateUser(t *testing.T) {
	user := CreateUser()
	update := CreateUser()
	update.Username = "updatedUser"
	update.ImageURL = "http://updateduser.com/image"

	err := user.UpdateUser(update)
	require.Nil(t, err)
	require.Equal(t, user.Username, update.Username)
	require.Equal(t, user.FirstName, update.FirstName)
	require.Equal(t, user.LastName, update.LastName)
	require.Equal(t, user.Email, update.Email)
	require.Equal(t, user.ImageURL, update.ImageURL)
	require.Greater(t, time.Now(), user.Updated)
}

func TestValidate(t *testing.T) {
	user := CreateUser()
	err := user.Validate()
	require.Nil(t, err)
}

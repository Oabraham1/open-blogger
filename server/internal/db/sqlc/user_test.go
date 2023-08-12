package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createDummyUser(username, email string) CreateNewUserParams {
	return CreateNewUserParams{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  "testPassword",
		FirstName: "testFirstName",
		LastName:  "testLastName",
		Interests: []string{"testInterest1", "testInterest2"},
		CreatedAt: time.Now(),
	}
}

func TestCreateNewUser(t *testing.T) {
	arg := createDummyUser("testUser", "testUser@email.com")
	user, err := testQueries.CreateNewUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Interests, user.Interests)

	// Delete the user
	err = testQueries.DeleteUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	getUserDeleted, err := testQueries.GetUserByID(context.Background(), user.ID)
	require.Error(t, err)
	require.Empty(t, getUserDeleted)
}

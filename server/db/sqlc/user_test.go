package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createDummyUser(username, email string) CreateNewUserParams {
	return CreateNewUserParams{
		Username:  username,
		Email:     email,
		Password:  "testPassword",
		FirstName: "testFirstName",
		LastName:  "testLastName",
	}
}

func TestUserCRUDOperations(t *testing.T) {
	arg := createDummyUser("testUser", "testUser@email.com")

	/* Test CreateNewUser */
	user, err := testStore.CreateNewUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)

	/* Test GetUserByID */
	getUser, err := testStore.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getUser)
	require.Equal(t, arg.Username, getUser.Username)
	require.Equal(t, arg.Email, getUser.Email)
	require.Equal(t, arg.FirstName, getUser.FirstName)
	require.Equal(t, arg.LastName, getUser.LastName)

	/* Test GetUserByUsername */
	getUserByUsername, err := testStore.GetUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, getUserByUsername)
	require.Equal(t, arg.Username, getUserByUsername.Username)
	require.Equal(t, arg.FirstName, getUserByUsername.FirstName)
	require.Equal(t, arg.LastName, getUserByUsername.LastName)
	require.Equal(t, getUser.Username, getUserByUsername.Username)
	require.Equal(t, getUser.FirstName, getUserByUsername.FirstName)
	require.Equal(t, getUser.LastName, getUserByUsername.LastName)
	require.Equal(t, getUser.Interests, getUserByUsername.Interests)

	/* Test GetPostsByUserID */
	posts, err := testStore.GetPostsByUserName(context.Background(), user.Username)
	require.NoError(t, err)
	require.Empty(t, posts)
	// Create 4 posts
	for i := 0; i < 4; i++ {
		_, err := testStore.CreateNewPost(context.Background(), CreateNewPostParams{
			Username: user.Username,
			Status:   StatusPublished,
			Category: "News",
			Title:    "testTitle",
			Body:     "testContent",
		})
		require.NoError(t, err)
	}
	posts, err = testStore.GetPostsByUserName(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, posts)
	require.Equal(t, 4, len(posts))

	/* Test GetPostsByUsername */
	getPostsByUserName, err := testStore.GetPostsByUserName(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, getPostsByUserName)
	require.Equal(t, 4, len(getPostsByUserName))
	require.Equal(t, posts[0].Body, getPostsByUserName[0].Body)

	/* Test UpdateUserInterestsByID */
	newInterests := []string{"newTestInterest1", "newTestInterest2"}
	err = testStore.UpdateUserInterestsByUsername(context.Background(), UpdateUserInterestsByUsernameParams{
		Username:  user.Username,
		Interests: newInterests,
	})
	require.NoError(t, err)
	updatedUser, err := testStore.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, newInterests, updatedUser.Interests)
	require.NotEqual(t, user.Interests, updatedUser.Interests)

	/* Test UpdatePostBodyByUserID */
	newBody := "newTestContent"
	updatedPost, err := testStore.UpdatePostBody(context.Background(), UpdatePostBodyParams{
		Username: user.Username,
		ID:       posts[0].ID,
		Body:     newBody,
	})
	require.NoError(t, err)
	require.NotEmpty(t, updatedPost)
	require.Equal(t, newBody, updatedPost.Body)
	require.NotEqual(t, posts[0].Body, updatedPost.Body)

	/* Teardown */

	// Delete all comments by user
	comments, err := testStore.GetCommentsByUserName(context.Background(), user.Username)
	require.NoError(t, err)
	for _, comment := range comments {
		err = testStore.DeleteCommentByID(context.Background(), comment.ID)
		require.NoError(t, err)
	}

	// Delete all posts by user
	for _, post := range posts {
		err = testStore.DeletePostByID(context.Background(), post.ID)
		require.NoError(t, err)
	}

	// Find all user sessions
	sessions, err := testStore.GetUserSessionsByUsername(context.Background(), user.Username)
	require.NoError(t, err)

	// Delete all user sessions
	for _, session := range sessions {
		err = testStore.DeleteSessionById(context.Background(), session.ID)
		require.NoError(t, err)
	}

	// Delete user
	err = testStore.DeleteUserAccount(context.Background(), user.Username)
	require.NoError(t, err)
	getUserDeleted, err := testStore.GetUserByID(context.Background(), user.ID)
	require.Error(t, err)
	require.Empty(t, getUserDeleted)
}

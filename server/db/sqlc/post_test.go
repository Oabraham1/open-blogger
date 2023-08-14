package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createDummyPost(t *testing.T, userId uuid.UUID, username string) CreateNewPostParams {
	return CreateNewPostParams{
		ID:        uuid.New(),
		Title:     "Test Post",
		Body:      "This is a test post",
		UserID:    userId,
		Username:  username,
		Status:    StatusDraft,
		Category:  "Test",
		CreatedAt: time.Now(),
	}
}

func TestPostCRUDOperations(t *testing.T) {
	userArg := createDummyUser("testUser123", "testUser@email.com")
	user, err := testStore.CreateNewUser(context.Background(), userArg)
	require.NoError(t, err)
	arg := createDummyPost(t, user.ID, user.Username)
	ctx := context.Background()

	/*
		Test Create Post
	*/
	post, err := testStore.CreateNewPost(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, post)
	require.Equal(t, arg.Title, post.Title)
	require.Equal(t, arg.Body, post.Body)
	require.Equal(t, arg.Username, post.Username)
	require.Equal(t, arg.Status, post.Status)
	require.Equal(t, arg.Category, post.Category)
	require.Equal(t, arg.CreatedAt.Day(), post.CreatedAt.Day())

	/*
		Test Get Post By ID
	*/
	getPost, err := testStore.GetPostById(ctx, post.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getPost)
	require.Equal(t, post.ID, getPost.ID)
	require.Equal(t, post.Title, getPost.Title)
	require.Equal(t, post.Body, getPost.Body)
	require.Equal(t, post.Username, getPost.Username)
	require.Equal(t, post.Status, getPost.Status)
	require.Equal(t, post.Category, getPost.Category)
	require.Equal(t, post.CreatedAt.Day(), getPost.CreatedAt.Day())

	/*
		Test Get Post By Category
	*/
	getPostsByCategory, err := testStore.GetPostsByCategory(ctx, post.Category)
	require.NoError(t, err)
	require.NotEmpty(t, getPostsByCategory)
	require.Greater(t, len(getPostsByCategory), 0)

	/*
		Test Get All Posts
		Create 4 posts
	*/
	for i := 0; i < 4; i++ {
		arg := createDummyPost(t, user.ID, user.Username)
		_, err := testStore.CreateNewPost(ctx, arg)
		require.NoError(t, err)
	}
	getAllPosts, err := testStore.GetAllPosts(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, getAllPosts)
	require.Equal(t, len(getAllPosts), 5)

	/*
		Test Update Post
	*/
	updatePost := UpdatePostBodyByPostIDAndUserIDParams{
		Body:   "This is an updated post",
		UserID: arg.UserID,
		ID:     post.ID,
	}
	updatedPost, err := testStore.UpdatePostBodyByPostIDAndUserID(ctx, updatePost)
	require.NoError(t, err)
	require.NotEmpty(t, updatedPost)
	require.Equal(t, updatePost.Body, updatedPost.Body)
	require.NotEqual(t, post.Body, updatedPost.Body)

	/*
		Test Delete Post
	*/
	err = testStore.DeletePostByID(ctx, post.ID)
	require.NoError(t, err)
	getPost, err = testStore.GetPostById(ctx, post.ID)
	require.Error(t, err)
	require.Empty(t, getPost)

	// Tear Down
	// Delete all posts
	posts, err := testStore.GetAllPosts(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, posts)
	for _, post := range posts {
		err = testStore.DeletePostByID(ctx, post.ID)
		require.NoError(t, err)
	}
	// Delete user
	err = testStore.DeleteUserByID(ctx, arg.UserID)
	require.NoError(t, err)
}

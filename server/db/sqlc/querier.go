// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateNewComment(ctx context.Context, arg CreateNewCommentParams) (Comment, error)
	CreateNewPost(ctx context.Context, arg CreateNewPostParams) (Post, error)
	CreateNewUser(ctx context.Context, arg CreateNewUserParams) (User, error)
	DeletePostByID(ctx context.Context, id uuid.UUID) error
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
	GetAllPosts(ctx context.Context) ([]GetAllPostsRow, error)
	GetCommentsByPostID(ctx context.Context, postID uuid.UUID) ([]Comment, error)
	GetPostById(ctx context.Context, id uuid.UUID) (Post, error)
	GetPostsByCategory(ctx context.Context, category string) ([]Post, error)
	GetPostsByUserID(ctx context.Context, userID uuid.UUID) ([]Post, error)
	GetPostsByUserName(ctx context.Context, username string) ([]Post, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	UpdatePostBodyByPostIDAndUserID(ctx context.Context, arg UpdatePostBodyByPostIDAndUserIDParams) (Post, error)
	UpdatePostStatus(ctx context.Context, arg UpdatePostStatusParams) (Post, error)
	UpdateUserInterestsByID(ctx context.Context, arg UpdateUserInterestsByIDParams) error
}

var _ Querier = (*Queries)(nil)
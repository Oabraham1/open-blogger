package db

import "context"

/* CreatePostTxParams contains the input parameters of the CreatePostTx function */
type CreatePostTxParams struct {
	CreateNewPostParams
	AfterCreate func(post Post) error
}

/* CreatePostTxResult is the result of the CreatePostTx function */
type CreatePostTxResult struct {
	Post Post
}

/* CreatePostTx creates a new post and executes the callback within a database transaction */
func (store *SQLStore) CreatePostTx(ctx context.Context, arg CreatePostTxParams) (CreatePostTxResult, error) {
	var result CreatePostTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Post, err = q.CreateNewPost(ctx, arg.CreateNewPostParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.Post)
	})

	return result, err
}

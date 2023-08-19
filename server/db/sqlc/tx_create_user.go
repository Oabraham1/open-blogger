package db

import "context"

/* CreateUserTxParams contains the input parameters of the CreateUserTx function */
type CreateUserTxParams struct {
	CreateNewUserParams
	AfterCreate func(user User) error
}

/* CreateUserTxResult is the result of the CreateUserTx function */
type CreateUserTxResult struct {
	User User
}

/* CreateUserTx creates a new user and executes the callback within a database transaction */
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateNewUser(ctx, arg.CreateNewUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}

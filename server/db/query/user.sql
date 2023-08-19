-- name: CreateNewUser :one
INSERT INTO users (username, password, email, first_name, last_name) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetPostsByUserName :many
SELECT * FROM posts WHERE username = $1;

-- name: UpdateUserInterestsByUsername :exec
UPDATE users SET interests = $1 WHERE username = $2;

-- name: UpdatePostBody :one
UPDATE posts SET body = $1, last_modified = $2 WHERE id = $3 AND username = $4 RETURNING *;

-- name: DeleteUserAccount :exec
DELETE FROM users WHERE username = $1;
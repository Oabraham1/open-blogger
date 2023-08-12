-- name: CreateNewUser :one
INSERT INTO users (username, password, email, first_name, last_name, interests, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, username, email, first_name, last_name, interests;

-- name: GetUserByUsername :one
SELECT username, first_name, last_name, interests FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT id, username, email, first_name, last_name, interests FROM users WHERE id = $1;

-- name: GetPostsByUserID :many
SELECT id, title, body, status, category, created_at, published_at, last_modified FROM posts WHERE user_id = $1;

-- name: GetPostsByUserName :many
SELECT id, title, body, status, category, created_at, published_at, last_modified FROM posts WHERE username = $1;

-- name: UpdateUserInterestsByID :exec
UPDATE users SET interests = $1 WHERE id = $2;

-- name: UpdatePostBodyByUserID :exec
UPDATE posts SET body = $1 WHERE user_id = $2;

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;
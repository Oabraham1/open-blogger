-- name: CreateNewPost :one
INSERT INTO posts (id, title, body, user_id, username, status, category, created_at, published_at, last_modified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: GetPostsByCategory :many
SELECT * FROM posts WHERE category = $1;

-- name: GetPostById :one
SELECT * FROM posts WHERE id = $1;

-- name: GetAllPosts :many
SELECT id, title, username, body, status, category, created_at, published_at, last_modified FROM posts;

-- name: GetCommentsByPostID :many
SELECT * FROM comments WHERE post_id = $1;

-- name: CreateNewComment :one
INSERT INTO comments (id, user_id, username, post_id, body, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: DeletePostByID :exec
DELETE FROM posts WHERE id = $1;
-- name: CreateNewPost :one
INSERT INTO posts (title, body, user_id, username, status, category, created_at, published_at, last_modified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, title, body, username, status, category, created_at, published_at, last_modified;

-- name: GetPostsByCategory :many
SELECT id, title, body, username, status, category, created_at, published_at, last_modified FROM posts WHERE category = $1;

-- name: GetPostById :one
SELECT id, title, body, username, status, category, created_at, published_at, last_modified FROM posts WHERE id = $1;

-- name: GetAllPosts :many
SELECT id, title, username, body, status, category, created_at, published_at, last_modified FROM posts;

-- name: GetCommentsByPostID :many
SELECT id, body, username, created_at FROM comments WHERE post_id = $1;

-- name: CreateNewComment :one
INSERT INTO comments (user_id, username, post_id, body, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, body, username, created_at;

-- name: DeletePostByID :exec
DELETE FROM posts WHERE id = $1;
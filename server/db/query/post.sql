-- name: CreateNewPost :one
INSERT INTO posts (title, body, username, status, category, published_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetPostsByCategory :many
SELECT * FROM posts WHERE category = $1;

-- name: GetPostById :one
SELECT * FROM posts WHERE id = $1;

-- name: GetAllPosts :many
SELECT id, title, username, body, status, category, created_at, published_at, last_modified FROM posts;

-- name: UpdatePostStatus :one
UPDATE posts SET status = $1, published_at = $2 WHERE id = $3 AND username = $4 RETURNING *;

-- name: CreateNewComment :one
INSERT INTO comments (username, post_id, body) VALUES ($1, $2, $3) RETURNING *;

-- name: GetCommentsByPostID :many
SELECT * FROM comments WHERE post_id = $1;

-- name: GetCommentsByUserName :many
SELECT * FROM comments WHERE username = $1;

-- name: GetCommentByID :one
SELECT * FROM comments WHERE id = $1;

-- name: DeleteCommentByID :exec
DELETE FROM comments WHERE id = $1;

-- name: DeletePostByID :exec
DELETE FROM posts WHERE id = $1;
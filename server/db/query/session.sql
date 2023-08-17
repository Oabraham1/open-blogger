-- name: CreateNewUserSession :one
INSERT INTO sessions (username, refresh_token, user_agent, client_ip, expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetSessionById :one
SELECT * FROM sessions WHERE id = $1 LIMIT 1;

-- name: GetUserSessionsByUsername :many
SELECT * FROM sessions WHERE username = $1 ORDER BY created_at DESC;

-- name: DeleteSessionById :exec
DELETE FROM sessions WHERE id = $1;
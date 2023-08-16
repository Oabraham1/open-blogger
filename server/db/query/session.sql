-- name: CreateNewUserSession :one
INSERT INTO sessions (username, refresh_token, user_agent, client_ip, expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetSessionById :one
SELECT * FROM sessions WHERE id = $1 LIMIT 1;
-- name: CreateSession :one
INSERT INTO sessions
(
  id,
  user_id,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at,
  encrypted_twitch_token
)
VALUES
($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.id = $1 AND s.user_id = $2
LIMIT 1;

-- name: GetSessionByRefreshToken :one
SELECT * FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.refresh_token = $1 AND s.is_blocked = false
LIMIT 1;

-- name: GetSessionByUserID :one
SELECT * FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE u.user_id = $1 AND is_blocked = false
ORDER BY s.created_at DESC
LIMIT 1;

-- name: ListSession :many
SELECT * FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.is_blocked = false
ORDER BY s.created_at ASC;

-- name: UpdateSession :exec
UPDATE sessions
SET
  encrypted_twitch_token = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;
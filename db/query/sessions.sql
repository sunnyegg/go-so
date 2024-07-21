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
SELECT * FROM sessions
WHERE id = $1 AND user_id = $2
LIMIT 1;

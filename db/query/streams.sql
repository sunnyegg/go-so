-- name: CreateStream :one
INSERT INTO streams
(
  user_id,
  title,
  game_name,
  started_at,
  created_by
)
VALUES
($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetStream :one
SELECT s.title, s.game_name, s.started_at, u.user_name as created_by_username
FROM streams s JOIN users u ON s.user_id = u.id
WHERE s.id = $1
  AND s.user_id = $2
LIMIT 1;

-- name: ListStreams :many
SELECT s.title, s.game_name, s.started_at, u.user_name as created_by_username
FROM streams s JOIN users u ON s.user_id = u.id
WHERE s.user_id = $3
ORDER BY s.started_at DESC
LIMIT $1
OFFSET $2;

-- name: DeleteStream :exec
DELETE FROM streams
WHERE id = $1;

-- name: GetStreamAttendanceMembers :many
SELECT s.title, s.game_name, s.started_at, am.username, am.is_shouted, am.present_at
FROM attendance_members as am JOIN streams as s ON am.stream_id = s.id
WHERE am.stream_id = $3 AND s.user_id = $4
ORDER BY am.present_at ASC
LIMIT $1
OFFSET $2;
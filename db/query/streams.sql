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
SELECT * FROM streams
WHERE id = $1 LIMIT 1;

-- name: ListStreams :many
SELECT * FROM streams
WHERE 1=1
  AND user_id = $3
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteStream :exec
DELETE FROM streams
WHERE id = $1;

-- name: GetStreamAttendanceMembers :many
SELECT s.title, s.game_name, s.started_at, am.username, am.is_shouted, am.present_at FROM attendance_members as am
JOIN streams as s ON am.stream_id = s.id
WHERE 1=1
  AND stream_id = $3
ORDER BY present_at ASC
LIMIT $1
OFFSET $2;
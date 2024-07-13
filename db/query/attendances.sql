-- name: CreateAttendance :one
INSERT INTO attendances
(
  user_id,
  stream_id
)
VALUES
($1, $2)
RETURNING *;

-- name: GetAttendance :one
SELECT * FROM attendances
WHERE id = $1 LIMIT 1;

-- name: ListAttendances :many
SELECT * FROM attendances
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteAttendance :exec
DELETE FROM attendances
WHERE id = $1;
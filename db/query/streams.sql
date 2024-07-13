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
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteStream :exec
DELETE FROM streams
WHERE id = $1;
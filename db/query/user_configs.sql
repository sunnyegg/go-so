-- name: CreateUserConfig :one
INSERT INTO user_configs
(
  user_id,
  config_type,
  "value"
)
VALUES
($1, $2, $3)
RETURNING *;

-- name: GetUserConfig :one
SELECT "value" FROM user_configs
WHERE 1=1
  AND user_id = $1
  AND config_type = $2
LIMIT 1;

-- name: UpdateUserConfig :one
UPDATE user_configs
SET "value" = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUserConfig :exec
DELETE FROM user_configs
WHERE id = $1;
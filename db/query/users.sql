-- name: CreateUser :one
INSERT INTO users
(
  user_id,
  user_login,
  user_name,
  profile_image_url,
  token
)
VALUES
($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT id, user_login, user_name, profile_image_url FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET user_login = $2, user_name = $3, profile_image_url = $4, token = $5, updated_at = now()
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserByUserID :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;
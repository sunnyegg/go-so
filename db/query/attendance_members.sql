-- name: CreateAttendanceMember :one
INSERT INTO attendance_members
(
  attendance_id,
  username,
  is_shouted,
  present_at
)
VALUES
($1, $2, $3, $4)
RETURNING *;

-- name: GetAttendanceMember :one
SELECT * FROM attendance_members
WHERE id = $1 LIMIT 1;

-- name: ListAttendanceMembers :many
SELECT * FROM attendance_members
ORDER BY id
LIMIT $1
OFFSET $2;

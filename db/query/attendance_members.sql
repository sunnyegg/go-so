-- name: CreateAttendanceMember :one
INSERT INTO attendance_members
(
  stream_id,
  username,
  is_shouted,
  present_at
)
VALUES
($1, $2, $3, $4)
RETURNING *;

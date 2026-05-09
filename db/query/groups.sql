-- name: CreateGroup :one
INSERT INTO groups (
  name,
  user_id
) VALUES (
  $1, $2
)
RETURNING id, name, user_id, created_at;

-- name: GetGroupsByUserID :many
SELECT id, name, user_id, created_at
FROM groups
WHERE user_id = $1
ORDER BY id DESC;

-- name: GetGroupByID :one
SELECT id, name, user_id, created_at
FROM groups
WHERE id = $1;

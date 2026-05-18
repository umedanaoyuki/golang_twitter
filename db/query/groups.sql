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

-- name: GetGroupsByMemberUserID :many
SELECT g.id, g.name, g.user_id, g.created_at
FROM groups g
INNER JOIN group_members gm ON g.id = gm.group_id
WHERE gm.user_id = $1
ORDER BY g.id DESC;

-- name: GetGroupByID :one
SELECT id, name, user_id, created_at
FROM groups
WHERE id = $1;

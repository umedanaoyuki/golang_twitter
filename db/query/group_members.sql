-- name: CreateGroupMember :one
INSERT INTO group_members (
  group_id,
  user_id
) VALUES (
  $1, $2
)
RETURNING group_id, user_id, created_at;

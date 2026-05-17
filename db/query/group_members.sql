-- name: CreateGroupMember :one
INSERT INTO group_members (
  group_id,
  user_id
) VALUES (
  $1, $2
)
RETURNING group_id, user_id, created_at;

-- name: ExistsGroupMember :one
SELECT EXISTS(
  SELECT 1
  FROM group_members
  WHERE group_id = $1 AND user_id = $2
) AS exists;
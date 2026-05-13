-- name: CreateGroupMember :one
INSERT INTO group_members (
  group_id,
  user_id,
  role,
  status,
  invited_by,
  accepted_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
ON CONFLICT (group_id, user_id) DO NOTHING
RETURNING id, group_id, user_id, role, status, invited_by, accepted_at, created_at;

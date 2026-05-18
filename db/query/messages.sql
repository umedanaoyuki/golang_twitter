-- name: CreateMessage :one
INSERT INTO messages (
  user_id,
  group_id,
  content
) VALUES (
  $1, $2, $3
)
RETURNING id, user_id, group_id, content, created_at;

-- name: GetMessagesByGroupID :many
SELECT id, user_id, group_id, content, created_at
FROM messages
WHERE group_id = $1
ORDER BY id DESC;

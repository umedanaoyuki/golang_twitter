-- name: CreateNotification :one
INSERT INTO notifications (
  user_id,
  actor_id,
  type,
  tweet_id,
  comment_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, user_id, actor_id, type, tweet_id, comment_id, is_read, read_at, created_at;

-- name: GetNotificationsByUserIDWithCursor :many
SELECT id, user_id, actor_id, type, tweet_id, comment_id, is_read, read_at, created_at
FROM notifications
WHERE user_id = $1
  AND (CASE WHEN $2 = 0 THEN true ELSE id < $2 END)
ORDER BY id DESC
LIMIT $3;

-- name: CountUnreadNotificationsByUserID :one
SELECT COUNT(*)::bigint AS unread_count
FROM notifications
WHERE user_id = $1 AND is_read = false;

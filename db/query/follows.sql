-- name: CreateFollow :one
INSERT INTO follows (
  user_id,
  followed_user_id
) VALUES (
  $1, $2
)
ON CONFLICT (user_id, followed_user_id) DO NOTHING
RETURNING id, user_id, followed_user_id, created_at;

-- name: DeleteFollow :exec
DELETE FROM follows
WHERE user_id = $1 AND followed_user_id = $2;

-- name: GetFollowersByUserIdWithCursor :many
SELECT id, user_id, followed_user_id, created_at
FROM follows
WHERE followed_user_id = $1
  AND (CASE WHEN $2 = 0 THEN true ELSE id < $2 END)
ORDER BY id DESC
LIMIT $3;

-- name: GetFollowingByUserIdWithCursor :many
SELECT id, user_id, followed_user_id, created_at
FROM follows
WHERE user_id = $1
  AND (CASE WHEN $2 = 0 THEN true ELSE id < $2 END)
ORDER BY id DESC
LIMIT $3;
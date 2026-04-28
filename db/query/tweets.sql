-- name: CreateTweet :one
INSERT INTO tweets (
  user_id,
  content
) VALUES (
  $1, $2
) RETURNING id, user_id, content, created_at, updated_at;

-- name: GetTweetsByUserID :many
SELECT id, user_id, content, created_at, updated_at
FROM tweets
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetAllTweets :many
SELECT id, user_id, content, created_at, updated_at
FROM tweets
ORDER BY created_at DESC
LIMIT $1;

-- name: GetTweetByID :one
SELECT id, user_id, content, created_at, updated_at
FROM tweets
WHERE id = $1;

-- name: DeleteTweet :exec
DELETE FROM tweets
WHERE id = $1 AND user_id = $2;

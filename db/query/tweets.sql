-- name: CreateTweet :one
INSERT INTO tweets (
  user_id,
  content,
  image_url
) VALUES (
  $1, $2, $3
) RETURNING id, user_id, content, image_url, created_at, updated_at;

-- name: GetTweetsByUserID :many
SELECT id, user_id, content, image_url, created_at, updated_at
FROM tweets
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetTweetsByUserIDWithCursor :many
SELECT id, user_id, content, image_url, created_at, updated_at
FROM tweets
WHERE user_id = $1
  AND (CASE WHEN $2 = 0 THEN true ELSE id < $2 END)
ORDER BY id DESC
LIMIT $3;

-- name: GetAllTweets :many
SELECT id, user_id, content, image_url, created_at, updated_at
FROM tweets
ORDER BY created_at DESC
LIMIT $1;

-- name: GetAllTweetsWithCursor :many
SELECT id, user_id, content, image_url, created_at, updated_at
FROM tweets
WHERE (CASE WHEN @cursor::int = 0 THEN true ELSE id < @cursor::int END)
ORDER BY id DESC
LIMIT @limit_count::int;

-- name: GetTweetByID :one
SELECT id, user_id, content, image_url, created_at, updated_at
FROM tweets
WHERE id = $1;

-- name: GetTweetsByIDs :many
SELECT id, user_id, content, image_url, created_at, updated_at
FROM tweets
WHERE id = ANY($1::int[]);

-- name: DeleteTweet :exec
DELETE FROM tweets
WHERE id = $1 AND user_id = $2;

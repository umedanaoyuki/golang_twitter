-- name: CreateImageTweet :one
INSERT INTO image_tweets (
  user_id,
  image_url
) VALUES (
  $1, $2
) RETURNING id, user_id, image_url, created_at, updated_at;

-- name: GetImageTweetsByUserID :many
SELECT id, user_id, image_url, created_at, updated_at
FROM image_tweets
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetImageTweetsByUserIDWithCursor :many
SELECT id, user_id, image_url, created_at, updated_at
FROM image_tweets
WHERE user_id = $1
  AND (CASE WHEN $2 = 0 THEN true ELSE id < $2 END)
ORDER BY id DESC
LIMIT $3;

-- name: GetAllImageTweets :many
SELECT id, user_id, image_url, created_at, updated_at
FROM image_tweets
ORDER BY created_at DESC
LIMIT $1;

-- name: GetImageTweetByID :one
SELECT id, user_id, image_url, created_at, updated_at
FROM image_tweets
WHERE id = $1;

-- name: GetImageTweetsByIDs :many
SELECT id, user_id, image_url, created_at, updated_at
FROM image_tweets
WHERE id = ANY($1::int[]);

-- name: DeleteImageTweet :exec
DELETE FROM image_tweets
WHERE id = $1 AND user_id = $2;

-- name: CreateLike :one
INSERT INTO likes (
  user_id,
  tweet_id
) VALUES (
  $1, $2
)
ON CONFLICT (user_id, tweet_id) DO NOTHING
RETURNING id, user_id, tweet_id, created_at;

-- name: DeleteLike :exec
DELETE FROM likes
WHERE user_id = $1 AND tweet_id = $2;

-- name: CountLikesByTweetID :one
SELECT COUNT(*)::bigint AS like_count
FROM likes
WHERE tweet_id = $1;

-- name: ExistsLike :one
SELECT EXISTS(
  SELECT 1
  FROM likes
  WHERE user_id = $1 AND tweet_id = $2
) AS exists;

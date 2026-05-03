-- name: CreateRetweet :one
INSERT INTO retweets (
  user_id,
  tweet_id
) VALUES (
  $1, $2
)
ON CONFLICT (user_id, tweet_id) DO NOTHING
RETURNING id, user_id, tweet_id, created_at;

-- name: DeleteRetweet :exec
DELETE FROM retweets
WHERE user_id = $1 AND tweet_id = $2;

-- name: CountLikesByTweetID :one
SELECT COUNT(*)::bigint AS retweet_count
FROM retweets
WHERE tweet_id = $1;

-- name: ExistsRetweet :one
SELECT EXISTS(
  SELECT 1
  FROM retweets
  WHERE user_id = $1 AND tweet_id = $2
) AS exists;
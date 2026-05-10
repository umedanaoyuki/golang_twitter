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

-- name: CountRetweetsByTweetID :one
SELECT COUNT(*)::bigint AS retweet_count
FROM retweets
WHERE tweet_id = $1;

-- name: CountRetweetsByTweetIDs :many
SELECT tweet_id, COUNT(*)::bigint AS retweet_count
FROM retweets
WHERE tweet_id = ANY($1::int[])
GROUP BY tweet_id;

-- name: ExistsRetweet :one
SELECT EXISTS(
  SELECT 1
  FROM retweets
  WHERE user_id = $1 AND tweet_id = $2
) AS exists;

-- name: GetUserRetweetsWithCursor :many
SELECT id, user_id, tweet_id, created_at
FROM retweets
WHERE user_id = $1
  AND (CASE WHEN $2 = 0 THEN true ELSE id < $2 END)
ORDER BY id DESC
LIMIT $3;
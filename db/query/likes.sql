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

-- name: CountLikesByTweetIDs :many
SELECT tweet_id, COUNT(*)::bigint AS like_count
FROM likes
WHERE tweet_id = ANY($1::int[])
GROUP BY tweet_id;

-- name: ExistsLike :one
SELECT EXISTS(
  SELECT 1
  FROM likes
  WHERE user_id = $1 AND tweet_id = $2
) AS exists;

-- name: GetUserLikesWithCursor :many
SELECT id, user_id, tweet_id, created_at
FROM likes
WHERE user_id = $1
  AND (CASE WHEN $2 = 0 THEN true ELSE id < $2 END)
ORDER BY id DESC
LIMIT $3;

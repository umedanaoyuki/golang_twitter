-- name: CreateComment :one
INSERT INTO comments (
  user_id,
  tweet_id,
  content
) VALUES (
  $1, $2, $3
)
RETURNING id, user_id, tweet_id, content, created_at;

-- name: DeleteComment :one
DELETE FROM comments
WHERE id = $1 AND user_id = $2 AND tweet_id = $3
RETURNING id;

-- name: GetCommentsByTweetIDWithCursor :many
SELECT id, user_id, tweet_id, content, created_at
FROM comments
WHERE tweet_id = $1
  AND (CASE WHEN $2 = 0 THEN true ELSE id < $2 END)
ORDER BY id DESC
LIMIT $3;

-- name: CountCommentsByTweetID :one
SELECT COUNT(*)::bigint AS comment_count
FROM comments
WHERE tweet_id = $1;

-- name: CountCommentsByTweetIDs :many
SELECT tweet_id, COUNT(*)::bigint AS comment_count
FROM comments
WHERE tweet_id = ANY($1::int[])
GROUP BY tweet_id;

-- name: ExistsComment :one
SELECT EXISTS(
  SELECT 1
  FROM comments
  WHERE user_id = $1 AND tweet_id = $2
) AS exists;

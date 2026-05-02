-- name: CreateBookmark :one
INSERT INTO bookmarks (
  user_id,
  tweet_id
) VALUES (
  $1, $2
)
ON CONFLICT (user_id, tweet_id) DO NOTHING
RETURNING id, user_id, tweet_id, created_at;

-- name: DeleteBookmark :exec
DELETE FROM bookmarks
WHERE user_id = $1 AND tweet_id = $2;

-- name: CountBookmarksByTweetID :one
SELECT COUNT(*)::bigint AS bookmark_count
FROM bookmarks
WHERE tweet_id = $1;

-- name: ExistsBookmark :one
SELECT EXISTS(
  SELECT 1
  FROM bookmarks
  WHERE user_id = $1 AND tweet_id = $2
) AS exists;
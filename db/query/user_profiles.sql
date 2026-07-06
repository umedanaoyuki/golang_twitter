-- name: CreateUserProfile :one
INSERT INTO user_profiles (
  user_id,
  name,
  bio,
  image_url,
  location
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id, user_id, name, bio, image_url, location, created_at, updated_at;

-- name: GetUserProfileByUserID :one
SELECT id, user_id, name, bio, image_url, location, created_at, updated_at
FROM user_profiles
WHERE user_id = $1
LIMIT 1;

-- name: UpdateUserProfile :one
UPDATE user_profiles
SET name = $2, bio = $3, image_url = $4, location = $5, updated_at = NOW()
WHERE user_id = $1
RETURNING id, user_id, name, bio, image_url, location, created_at, updated_at;

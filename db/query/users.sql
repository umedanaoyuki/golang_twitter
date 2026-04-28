-- name: CreateUser :one
INSERT INTO users (
  email,
  password,
  is_active
) VALUES (
  $1, $2, $3
) RETURNING id, email, password, is_active, created_at, updated_at;

-- name: UpdateUserIsActive :exec
UPDATE users
SET is_active = $2, updated_at = NOW()
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password, is_active, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserDetailByUserID :one
SELECT id, email, password, is_active, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;
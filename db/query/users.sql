-- name: CreateUser :one
INSERT INTO users (
  email,
  password,
  is_active
) VALUES (
  $1, $2, $3,
) RETURNING id, email, password, is_active, created_at, updated_at;
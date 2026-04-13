-- name: CreateUserActivation :one
INSERT INTO user_activations (
  user_id,
  token,
  expired_at
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetUserActivationByToken :one
SELECT * FROM user_activations
WHERE token = $1 AND expired_at > NOW()
LIMIT 1;

-- name: DeleteUserActivation :exec
DELETE FROM user_activations
WHERE token = $1;

-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_up, email)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(), 
  $1
)
RETURNING *;

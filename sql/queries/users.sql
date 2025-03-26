-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(), 
  $1,
  $2
)
RETURNING id, created_at, updated_at, email;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE users.email = $1;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2, email = $3, updated_at = NOW()
WHERE users.id = $1
RETURNING id, created_at, updated_at, email;

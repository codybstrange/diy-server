-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  $3,
  null
);

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE refresh_tokens.token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $3
WHERE refresh_tokens.token = $1;

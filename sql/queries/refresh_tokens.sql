-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
 VALUES ($1, NOW(), NOW(), $2,$3,$4) RETURNING *; 

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT u.id, u.email, u.created_at, u.updated_at
FROM refresh_tokens rt
JOIN users u ON rt.user_id = u.id
WHERE rt.token = $1
  AND rt.revoked_at IS NULL
  AND rt.expires_at > NOW();

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW() WHERE token=$1 RETURNING *;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING token, created_at, updated_at, expires_at;

-- name: LookupToken :one
SELECT * FROM refresh_tokens
WHERE token = $1 AND expires_at > NOW() AND revoked_at IS NULL;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

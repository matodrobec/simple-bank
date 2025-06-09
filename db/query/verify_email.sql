-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
  username,
  email,
  secret_code,
  expired_at
) VALUES (
  @username, @email, @secret_code, COALESCE(sqlc.narg(expired_at)::timestamptz, now() + interval '15 minutes')
) RETURNING *;


-- name: GetVerifyEmail :one
SELECT * FROM verify_emails
WHERE id = $1 LIMIT 1;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET
  is_used = TRUE
WHERE
  id = @id
  AND secret_code = @secret_code
  AND is_used = FALSE
  AND expired_at > now()
RETURNING *;
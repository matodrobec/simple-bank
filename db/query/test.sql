-- name: CreateTest :one
INSERT INTO test (
  note
) VALUES (
  $1
) RETURNING *;
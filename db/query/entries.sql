-- name: CreateEntry :one
INSERT INTO entries (account_id,amount) VALUES ($1,$2) RETURNING *;

-- name: GetEntryById :one
SELECT * FROM entries
WHERE id = $1;

-- name: GetEntriesByAmount :many
SELECT * FROM entries
WHERE amount = $1;

-- name: GetEntriesByAccountId :many
SELECT * FROM entries
WHERE account_id = $1;






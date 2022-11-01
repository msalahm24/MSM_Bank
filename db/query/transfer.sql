-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id,to_account_id,amount) VALUES ($1,$2,$3) RETURNING *;

-- name: GetTransferById :one
SELECT * FROM transfers
WHERE id = $1;

-- name: GetTransfersByAmount :many
SELECT * FROM transfers
WHERE amount = $1;

-- name: GetTransfersByFromAccountId :many
SELECT * FROM transfers
WHERE from_account_id = $1;

-- name: GetTransfersByToAccountId :many
SELECT * FROM transfers
WHERE to_account_id = $1;
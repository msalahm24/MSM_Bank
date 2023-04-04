-- name: CreateUser :one
INSERT INTO users 
(username,
hashed_pass,
full_name,
email
) 
VALUES ($1,$2,$3,$4) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;
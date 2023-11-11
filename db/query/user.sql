-- name: CreateUser :one
INSERT INTO users (
  username, password, full_name, email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- Ini menggunakan select biasa, sehingga tidak handling untuk antri sampai ada request yang rollback atau commit
-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;
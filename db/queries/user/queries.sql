-- name: GetUser :one
SELECT id, name, email, created_at, updated_at
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, name, email, created_at, updated_at
FROM users
ORDER BY created_at DESC;

-- name: CreateUser :exec
INSERT INTO users (id, name, email)
VALUES ($1, $2, $3);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

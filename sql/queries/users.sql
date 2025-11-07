-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = $3
WHERE id = $4
RETURNING *;

-- name: UpdateUserIsChirpyRed :exec
UPDATE users
SET is_chirpy_red = $1, updated_at = $2
WHERE id = $3;
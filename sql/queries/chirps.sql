-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2
    )
RETURNING *;

-- name: GetAllChirps :many
SELECT *
FROM chirps
WHERE user_id = coalesce(sqlc.narg('user_id'), user_id)
ORDER BY 
    CASE WHEN sqlc.arg('order')::text = 'asc' then created_at end ASC,
    CASE WHEN sqlc.arg('order') = 'desc' then created_at end DESC;

-- name: GetChirp :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;
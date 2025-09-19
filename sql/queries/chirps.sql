-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id) VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2) RETURNING *;

-- name: GetAllChirps :many
SELECT *
FROM chirps
ORDER BY
  CASE WHEN @sort_order = 'asc' THEN created_at END ASC,
  CASE WHEN @sort_order = 'desc' THEN created_at END DESC;

-- name: GetChrirp :one
SELECT * FROM chirps WHERE id=$1;

-- name: DeleteChirpByID :exec
DELETE FROM chirps WHERE id=$1;

-- name: GetChirpsByAuthor :many
SELECT *
FROM chirps
WHERE user_id = @user_id
ORDER BY
  CASE WHEN @sort_order = 'asc' THEN created_at END ASC,
  CASE WHEN @sort_order = 'desc' THEN created_at END DESC;

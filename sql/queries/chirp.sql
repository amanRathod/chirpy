-- name: CreateChirp :one
INSERT INTO chirp (id, created_at, updated_at, body, user_id)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * from chirp
ORDER BY created_at ASC;

-- name: GetChirpById :one
SELECT * from chirp
WHERE id = $1;


-- name: DeleteChirpById :exec
DELETE from chirp where id=$1;
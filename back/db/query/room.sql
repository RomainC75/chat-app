-- name: GetRoom :one
SELECT * FROM rooms WHERE id = $1;

-- name: GetRooms :one
SELECT * FROM rooms ;

-- name: CreateRoom :one
INSERT INTO rooms (
    id,
    user_id,
    name,
    created_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;


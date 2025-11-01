-- name: GetMessage :one
SELECT * FROM messages WHERE id = $1;

-- name: GetMessagesByRoom :one
SELECT * FROM messages WHERE room_id = $1 ORDER BY created_at;

-- name: CreateMessage :one
INSERT INTO messages (
    id,
    user_id,
    room_id,
    message,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;


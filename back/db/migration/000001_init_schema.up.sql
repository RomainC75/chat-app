CREATE TABLE users (
    id   uuid PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password  TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

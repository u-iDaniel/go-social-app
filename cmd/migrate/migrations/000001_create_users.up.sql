CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL, -- Using citext for case-insensitive email
    username VARCHAR(255) NOT NULL UNIQUE,
    password bytea NOT NULL,
    created_at timestamp(0) WITH TIME ZONE DEFAULT NOW()
);
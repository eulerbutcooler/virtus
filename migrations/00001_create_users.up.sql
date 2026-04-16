CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_roles AS ENUM('member','admin','institution');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    role user_roles NOT NULL DEFAULT 'member',
    password_hash TEXT NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_emails ON users(email);

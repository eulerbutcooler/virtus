CREATE TYPE contribution_status AS ENUM (
    'pending', 'completed', 'failed', 'refunded'
);

CREATE TABLE contributions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    pool_id       UUID NOT NULL REFERENCES pools(id) ON DELETE RESTRICT,
    amount        NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    currency      VARCHAR(3) NOT NULL DEFAULT 'USD',
    status        contribution_status NOT NULL DEFAULT 'pending',
    payment_ref   VARCHAR(255),
    stripe_pi_id  VARCHAR(255) UNIQUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contributions_user_id ON contributions(user_id);
CREATE INDEX idx_contributions_status  ON contributions(status);

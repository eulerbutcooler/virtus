CREATE TYPE request_status AS ENUM (
    'draft', 'submitted', 'verified', 'queued',
    'funded', 'procuring', 'delivered', 'completed', 'rejected'
);

CREATE TYPE urgency_level AS ENUM ('critical', 'high', 'standard', 'low');

CREATE TABLE requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    item_category   VARCHAR(100) NOT NULL,
    item_name       VARCHAR(255) NOT NULL,
    description     TEXT NOT NULL,
    urgency         urgency_level NOT NULL DEFAULT 'standard',
    estimated_cost  NUMERIC(15,2) NOT NULL CHECK (estimated_cost > 0),
    justification   TEXT NOT NULL,
    status          request_status NOT NULL DEFAULT 'draft',
    rejection_note  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_requests_user_id ON requests(user_id);
CREATE INDEX idx_requests_status  ON requests(status);
CREATE INDEX idx_requests_urgency ON requests(urgency);

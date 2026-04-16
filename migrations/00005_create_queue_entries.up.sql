CREATE TABLE queue_entries (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id            UUID NOT NULL UNIQUE REFERENCES requests(id) ON DELETE CASCADE,
    position              INTEGER NOT NULL DEFAULT 0,
    priority_score        NUMERIC(10,4) NOT NULL DEFAULT 0.0,
    funding_progress      NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    estimated_fulfillment DATE,
    entered_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_queue_entries_priority ON queue_entries(priority_score DESC);
CREATE INDEX idx_queue_entries_position ON queue_entries(position ASC);

CREATE TABLE impact_records (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    delivery_id          UUID NOT NULL REFERENCES deliveries(id) ON DELETE RESTRICT,
    user_id              UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    interval_label       VARCHAR(50) NOT NULL, -- 7_days, 90_days, etc
    outcome_description  TEXT,
    satisfaction_score   SMALLINT CHECK (satisfaction_score BETWEEN 1 AND 10),
    metrics              JSONB NOT NULL DEFAULT '{}',
    recorded_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_impact_records_delivery_id ON impact_records(delivery_id);
CREATE INDEX idx_impact_records_user_id     ON impact_records(user_id);

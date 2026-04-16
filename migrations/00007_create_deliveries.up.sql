CREATE TYPE delivery_status AS ENUM (
    'in_transit', 'delivered', 'failed', 'returned'
);

CREATE TABLE deliveries (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fulfillment_id   UUID NOT NULL UNIQUE REFERENCES fulfillments(id) ON DELETE RESTRICT,
    tracking_number  VARCHAR(255),
    carrier          VARCHAR(100),
    proof_photo_url  TEXT,
    status           delivery_status NOT NULL DEFAULT 'in_transit',
    delivered_at     TIMESTAMPTZ,
    verified_at      TIMESTAMPTZ,
    verified_by      UUID REFERENCES users(id),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_deliveries_fulfillment_id ON deliveries(fulfillment_id);
CREATE INDEX idx_deliveries_status         ON deliveries(status);

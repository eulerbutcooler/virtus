CREATE TYPE fulfillment_status AS ENUM (
    'pending', 'vendor_selected', 'ordered', 'shipped', 'delivered', 'cancelled'
);

CREATE TABLE fulfillments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id          UUID NOT NULL UNIQUE REFERENCES requests(id) ON DELETE RESTRICT,
    vendor_name         VARCHAR(255),
    vendor_ref          VARCHAR(255),
    actual_cost         NUMERIC(15,2),
    procurement_status  fulfillment_status NOT NULL DEFAULT 'pending',
    notes               TEXT,
    procured_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fulfillments_request_id ON fulfillments(request_id);
CREATE INDEX idx_fulfillments_status     ON fulfillments(procurement_status);

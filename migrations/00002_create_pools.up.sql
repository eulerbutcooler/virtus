CREATE TABLE pools (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance         NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    total_in        NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    total_out       NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    currency        VARCHAR(3) NOT NULL DEFAULT 'USD',
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO pools (id, currency)
VALUES ('00000000-0000-0000-0000-000000000001', 'USD');
